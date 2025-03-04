package main

import (
	"context"
	"fmt"

	"github.com/Armody/Gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, c command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("no logged in user: %w", err)
		}

		return handler(s, c, user)
	}
}
