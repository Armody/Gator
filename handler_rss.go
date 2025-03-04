package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Armody/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, c command, user database.User) error {
	if len(c.args) != 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
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

	if _, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}); err != nil {
		return fmt.Errorf("error following feed: %w", err)
	}

	fmt.Println("Feed created successfully:")
	fmt.Println(feed)
	return nil
}

func handlerListFeeds(s *state, c command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return nil
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error getting user: %w", err)
		}
		fmt.Println("Feed:")
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println("Feed created by:", user.Name)
		fmt.Println("===========================")
	}

	return nil
}

func handlerFollowFeed(s *state, c command, user database.User) error {
	if len(c.args) != 1 {
		return fmt.Errorf("usage: follow <url>")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), c.args[0])
	if err != nil {
		return fmt.Errorf("error getting feed: %w", err)
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error following: %w", err)
	}

	fmt.Println(follow.UserName, "followed", follow.FeedName)

	return nil
}

func handlerUnfollowFeed(s *state, c command, user database.User) error {
	if len(c.args) != 1 {
		return fmt.Errorf("usage: unfollow <url>")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), c.args[0])
	if err != nil {
		return fmt.Errorf("error getting feed: %w", err)
	}

	if err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return fmt.Errorf("error deleting feed: %w", err)
	}

	fmt.Println(user.Name, "unfollowed", feed.Name, "feed")

	return nil
}

func handlerFollowList(s *state, c command) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting follow list: %w", err)
	}

	if len(follows) == 0 {
		fmt.Println("No feed follows found for user", s.cfg.CurrentUserName)
		return nil
	}

	fmt.Println("User", follows[0].UserName, "follows:")
	for _, follow := range follows {
		fmt.Println(follow.FeedName)
	}

	return nil
}
