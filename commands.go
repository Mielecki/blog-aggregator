package main

import (
	"errors"
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	cmdHandler map[string]func(*state, command) error
}

// This function handles the login command, setting username in the given state.config struct
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the login command expects a single argument, the username")
	}

	username := cmd.args[0]

	if err := s.config.SetUser(username); err != nil {
		return err
	}

	fmt.Println("The user " + username + " has been set")

	return nil
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