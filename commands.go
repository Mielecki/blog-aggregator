package main

import (
	"errors"
)

type command struct {
	name string
	args []string
}

type commands struct {
	cmdHandler map[string]func(*state, command) error
}

// This method registers a new handler function for a command name
func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdHandler[name] = f
}

// This method runs a given command with the provided state if it exists
func (c *commands) run(s *state, cmd command) error {
	if s == nil {
		return errors.New("a state does not exists")
	}
	
	handler, ok := c.cmdHandler[cmd.name]
	if !ok {
		return errors.New("a given command does not exists")
	}

	if err := handler(s, cmd); err != nil {
		return err
	}

	return nil
}