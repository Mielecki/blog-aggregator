package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Mielecki/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

// This function handles the login command, setting username in the given state.config struct if it exisits in db
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the login command expects a single argument, the username")
	}

	username := cmd.args[0]

	if _, err := s.db.GetUser(context.Background(), username); err != nil {
		return err
	}

	if err := s.config.SetUser(username); err != nil {
		return err
	}

	fmt.Println("The user " + username + " has been set")

	return nil
}

// This function handles the register command, inserting new user into users table and setting username in the given state.config struct
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the register command expects a single argument, the username")
	}

	username := cmd.args[0]

	_, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: username,
	})
	if err != nil {
		return err
	}

	if err := s.config.SetUser(username); err != nil {
		return err
	}

	fmt.Println("The user " + username + " has been created")

	return nil
}

// Thus fucntion handles the users command, listing all existing users
func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.config.CurrentUserName {
			fmt.Println("* " + user.Name + " (current)")
		} else {
			fmt.Println("* " + user.Name)
		}
	}

	return nil
}

// This function handles the agg command, fetching RSS from a given website
func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the agg command expects a single argument, the time_between_reqs")
	}

	time_unparsed := cmd.args[0]
	time_between_reqs, err := time.ParseDuration(time_unparsed)
	if err != nil {
		return err
	}

	fmt.Println("Collecting feeds every " + time_between_reqs.String())
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Println(err.Error())
		}
	}			
}

// This function handles the addfeed command, adding a new feed to feeds table
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("the addfeed command expects two arguments, the name and URL of the feed")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	_, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: user.ID,
	})
	if err != nil {
		return err
	}

	err = handlerFollow(s, command{
		name: "follow",
		args: []string{url},
	}, user)

	if err != nil {
		return nil
	}

	return nil
}

// This funciton handles feeds commands, listing all existing feeds
func handlerFeeds(s *state, cmd command) error {
	feedsUser, err := s.db.GetFeedsWithUsername(context.Background())
	if err != nil {
		return err
	}

	for _, item := range feedsUser {
		fmt.Printf(`* Feed: "%s" URL: "%s" User: "%s"\n`, item.Name, item.Url, item.Name_2)
		fmt.Println()
	}

	return nil
}

// This function handles the follow command, creating a new record in feed_follow table joining feeds with users table
func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("the follow command expects one argument, the URL of the feed to follow")
	}

	url := cmd.args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	feed_follows, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf(`User: "%s" Followed feed: "%s"`, feed_follows.Name_2, feed_follows.Name)
	fmt.Println()
	return nil
}

// This function handles the following command, listing all feeds followed by the current user
func handlerFollowing(s *state, cmd command, user database.User) error {
	feed_follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, item := range feed_follows {
		fmt.Printf(`* Feed: "%s" Created by: "%s"`, item.FeedName, item.CreatedBy)
		fmt.Println()
	}

	return nil
}

// This function handles the unfollow command, removing record from feed_follows table
func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("the unfollow command expects one argument, the URL of the feed to unfollow")
	}

	url := cmd.args[0]
	err := s.db.DeleteFeedFollowRecord(context.Background(), database.DeleteFeedFollowRecordParams{
		UserID: user.ID,
		Url: url,
	})
	if err != nil {
		return err
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) > 0 {
		atoi, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
		limit = atoi
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return err
	}

	for _, item := range posts {
		fmt.Println("***********************************")
		fmt.Println("Published at: \n" + item.PublishedAt.Time.Format(time.RFC822))
		fmt.Println("From: " + item.Name)
		fmt.Println("Title: \n" + item.Title)
		fmt.Println("Descrption: \n" + item.Description.String)
		fmt.Println("Link: " + item.Url)
		fmt.Println("***********************************")
	}
	return nil 
}

// This function handles the reset command, deleting all records from users table (ONLY FOR DEVELOPMENT PURPOSE!!!)
func handlerReset(s *state, cmd command) error {
	if err := s.db.Reset(context.Background()); err != nil {
		return err
	}

	fmt.Println("Resetting was successful")

	return nil
}