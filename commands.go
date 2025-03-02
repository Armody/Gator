package main

import (
	"errors"
)

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, exists := c.registeredCommands[cmd.name]
	if !exists {
		return errors.New("command not found")
	}
	return f(s, cmd)
}
