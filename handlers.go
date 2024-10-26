package main

import (
	"context"
	"errors"
	"fmt"
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
	fetchURL := "https://www.wagslane.dev/index.xml"
	RSSFeed, err := fetchFeed(context.Background(), fetchURL)
	if err != nil {
		return err
	}

	fmt.Println(RSSFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("the login command expects two arguments, the name and URL of the feed")
	}

	name := cmd.args[0]
	url := cmd.args[1]
	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}	

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
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

	fmt.Println(feed)

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