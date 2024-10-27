package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/Mielecki/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func scrapeFeeds(s *state,) error {
	to_fetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	if err := s.db.MarkFeedFetched(context.Background(), to_fetch.ID); err != nil {
		return err
	}

	fetchURL := to_fetch.Url
	RSSFeed, err := fetchFeed(context.Background(), fetchURL)
	if err != nil {
		return err
	}

	for _, item := range RSSFeed.Channel.Item {
		published_at := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			published_at = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: item.Title,
			Url: item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: published_at,
			FeedID: to_fetch.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
				continue
			}
			log.Println("Couldnot create post: " + err.Error())
		}
	}

	log.Printf(`Feed: "%s" collected`, to_fetch.Name)
	return nil
}