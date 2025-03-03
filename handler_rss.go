package main

import (
	"context"
	"fmt"

	"github.com/Armody/Gator/internal/RSS"
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
