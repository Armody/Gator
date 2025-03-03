package RSS

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't form request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't get response: %w", err)
	}
	defer res.Body.Close()

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read data: %w", err)
	}

	var feed RSSFeed
	if err := xml.Unmarshal(dat, &feed); err != nil {
		return nil, fmt.Errorf("couldn't unwrap response: %w", err)
	}

	feed = cleanFeed(&feed)

	return &feed, nil
}

func cleanFeed(feed *RSSFeed) RSSFeed {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return *feed
}
