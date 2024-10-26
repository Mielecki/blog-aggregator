package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Mielecki/blog-aggregator/internal/config"
	"github.com/Mielecki/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db *database.Queries
	config *config.Config
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("not enough arguments provided")
	}
	config_file, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", config_file.DbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	s := state{
		db: dbQueries,
		config: &config_file,
	}
	cmds := commands{make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)
	cmd := command{
		name: args[1],
		args: args[2:],
	}
	if err := cmds.run(&s, cmd); err != nil {
		log.Fatal(err)
	}
}
