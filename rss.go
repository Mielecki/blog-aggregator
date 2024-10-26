package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// This function fetches the RSS feed from a given website URL
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Add("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()
	
	rss_feed := RSSFeed{}
	if err := xml.Unmarshal(data, &rss_feed); err != nil {
		return &RSSFeed{}, err
	}

	rss_feed.unescapeFeed()

	return &rss_feed, nil
}

// This function decodes escaped HTML entites in the entire RSSFeed
func (rss_feed *RSSFeed) unescapeFeed() {
	rss_feed.Channel.Title = html.UnescapeString(rss_feed.Channel.Title)
	rss_feed.Channel.Description = html.UnescapeString(rss_feed.Channel.Description)

	for i := range rss_feed.Channel.Item {
		rss_feed.Channel.Item[i].unescapeItem()
	}
}

// This function decodes escaped HTML entites in the RSSItem
func (rss_item *RSSItem) unescapeItem() {
	rss_item.Title = html.UnescapeString(rss_item.Title)
	rss_item.Description = html.UnescapeString(rss_item.Description)
}