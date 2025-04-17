package main

import (
	"context"
	"encoding/xml"
	"fmt"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't create the request: %w", err)
	}
	request.Header.Set("User-Agent", "gator")

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("couldn't create the client: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to fetch feed: received status code %d", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't do the request: %w", err)
	}

	rssFeed := &RSSFeed{}
	err = xml.Unmarshal(data, rssFeed)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal the response data: %w", err)
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(rssFeed.Channel.Item[i].Description)
	}

	return rssFeed, nil
}
