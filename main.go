package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Mielecki/blog-aggregator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("not enough arguments provided")
	}
	config_file, err := config.Read()
	if err != nil {
		fmt.Println(err.Error())
	}
	s := state{&config_file}
	cmds := commands{make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmd := command{
		name: "login",
		args: args[2:],
	}
	if err := cmds.run(&s, cmd); err != nil {
		log.Fatalf(err.Error())
	}
}
