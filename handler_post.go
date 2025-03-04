package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Armody/Gator/internal/database"
)

func handlerBrowse(s *state, c command, user database.User) error {
	limit := 2
	if len(c.args) > 0 {
		l, err := strconv.Atoi(c.args[0])
		if err != nil {
			return fmt.Errorf("usage: browse <limit>")
		}
		limit = l
	}

	posts, err := s.db.GetPostsByUser(context.Background(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error getting posts: %w", err)
	}

	for _, post := range posts {
		fmt.Println(post.Title, post.Url)
		fmt.Println("posted at", post.PublishedAt.Time.Format("02 Jan 2006"))
		fmt.Println("----------------------------")
		fmt.Println(post.Description.String)
		fmt.Println("============================")
	}

	return nil
}
