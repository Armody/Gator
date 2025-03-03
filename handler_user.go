package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Armody/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("expected a username. Usage: gator login <username>")
	}

	userName := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return errors.New("user doesn't exist")
	}

	err = s.cfg.SetUser(userName)
	if err != nil {
		return err
	}
	fmt.Println("User has been set to", userName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("expected a username. Usage: gator login <username>")
	}

	userName := cmd.args[0]

	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userName,
	})
	if err != nil {
		return errors.New("user with that username already exists")
	}

	s.cfg.SetUser(userName)

	fmt.Println("user", userName, "was created")
	fmt.Println(newUser)

	return nil
}
