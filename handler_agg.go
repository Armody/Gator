package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Armody/Gator/internal/RSS"
	"github.com/Armody/Gator/internal/database"
	"github.com/google/uuid"
)

const timeLayout = "Mon, 02 Jan 2006 15:04:05 -0700"

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
		log.Println("error getting feed: ", err)
		return
	}
	if err := s.db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		log.Println("error marking feed: ", err)
		return
	}

	rssFeed, err := RSS.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Println("error fetching rss feed: ", err)
		return
	}

	for _, post := range rssFeed.Channel.Item {
		pubDate, err := time.Parse(timeLayout, post.PubDate)
		if err != nil {
			log.Println("error formatting time:", err)
			return
		}

		if _, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     post.Title,
			Url:       post.Link,
			Description: sql.NullString{
				String: post.Description,
				Valid:  true,
			},
			PublishedAt: sql.NullTime{
				Time:  pubDate,
				Valid: true,
			},
			FeedID: feed.ID,
		}); err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Println("error creating post:", err)
			return
		}
	}
	fmt.Println("added", len(rssFeed.Channel.Item), "posts to database")
}
