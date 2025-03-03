package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Armody/Gator/internal/RSS"
	"github.com/Armody/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, c command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := RSS.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, c command) error {
	if len(c.args) != 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("no user logged in")
	}
	name := c.args[0]
	url := c.args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	fmt.Println("Feed created successfully:")
	fmt.Println(feed)
	return nil
}
