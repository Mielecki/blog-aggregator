package main

import (
	"context"
	"fmt"
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
		fmt.Println(item.Title)
	}
	
	return nil
}