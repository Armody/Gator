package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("expected a username. Usage: gator login <username>")
	}

	userName := cmd.args[0]

	err := s.cfg.SetUser(userName)
	if err != nil {
		return err
	}
	fmt.Println("User has been set to", userName)
	return nil
}
