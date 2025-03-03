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
		return fmt.Errorf("expected a username. Usage: gator login <username>")
	}

	userName := cmd.args[0]

	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userName,
	})
	if err != nil {
		return fmt.Errorf("user with that username already exists: %w", err)
	}

	s.cfg.SetUser(userName)

	fmt.Println("user", userName, "was created")
	fmt.Println(newUser)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ClearUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't reset the table: %w", err)
	}

	fmt.Println("Succesfully cleared the users table")

	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Println("*", user.Name, "(current)")
			continue
		}
		fmt.Println("*", user.Name)
	}

	return nil
}
