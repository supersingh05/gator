package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/supersingh05/gator/internal/config"
	"github.com/supersingh05/gator/internal/database"
)

type state struct {
	config *config.Config
	db     *database.Queries
	dbObj  *sql.DB
}

func main() {
	c, err := config.Read("/home/asavirkalla/.gatorconfig.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", c.DbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	s := state{config: c, db: dbQueries, dbObj: db}

	cmds := commands{make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handleUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handler_addfeed))
	cmds.register("feeds", handler_feeds)
	cmds.register("follow", middlewareLoggedIn(handler_followfeeds))
	cmds.register("following", middlewareLoggedIn(handler_following))
	cmds.register("unfollow", middlewareLoggedIn(handler_unfollow))
	cmds.register("browse", middlewareLoggedIn(handler_browse))
	args := os.Args[1:]
	if len(args) < 1 {
		os.Exit(1)
	}
	cmd := command{name: args[0], args: args[1:]}

	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		var name sql.NullString
		name.Scan(s.config.CurrentUserName)
		user, err := s.db.GetUserByName(context.Background(), name)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
