package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Armody/Gator/internal/RSS"
	"github.com/Armody/Gator/internal/database"
)

func handlerAgg(s *state, c command, user database.User) error {
	if len(c.args) != 1 {
		return fmt.Errorf("usage: agg <interval-between-requests>")
	}

	timeBetweenRequests, err := time.ParseDuration(c.args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s, user)
	}
}

func scrapeFeeds(s *state, user database.User) {
	feed, err := s.db.GetNextFeedToFetch(context.Background(), user.ID)
	if err != nil {
		return
	}
	if err := s.db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		return
	}

	rssFeed, err := RSS.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return
	}

	fmt.Println(rssFeed.Channel.Title)
	for i, post := range rssFeed.Channel.Item {
		if i < 10 {
			fmt.Println("	-", post.Title)
		}
	}
}
