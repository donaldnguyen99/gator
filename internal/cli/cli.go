package cli

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/donaldnguyen99/gator/internal/config"
	"github.com/donaldnguyen99/gator/internal/database"
	_ "github.com/lib/pq"
)

type CLI struct {
	name     string
	commands *commands
	state    *state
}

func NewCLI(name string) *CLI {
	commands := &commands{
		handlers: make(map[string]handler),
	}
	return &CLI{
		name:     name,
		commands: commands,
	}
}

// Need to pass in parsed args
func (cli *CLI) Run() error {
	config, err := config.Read()
	if err != nil {
		return err
	}
	db, err := sql.Open("postgres", config.DbUrl)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	cli.state = &state{
		db: database.New(db),
		config: config,
	}

	cli.commands.register("login", handlerLogin)
	cli.commands.register("register", handlerRegisterUser)
	cli.commands.register("users", handlerGetUsers)
	cli.commands.register("reset", handlerReset)
	cli.commands.register("agg", handlerAggregateFeeds)
	cli.commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cli.commands.register("feeds", handlerGetFeeds)
	cli.commands.register("follow", middlewareLoggedIn(handlerFollowFeed))
	cli.commands.register("following", middlewareLoggedIn(handlerListFollows))
	cli.commands.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))
	cli.commands.register("browse", middlewareLoggedIn(handlerBrowsePosts))

	args := os.Args[1:]
	if len(args) == 0 {
		return fmt.Errorf("no %s subcommand provided", cli.name)
	}
	command := command{
		name: args[0],
		args: args[1:],
	}

	if err := cli.commands.run(cli.state, command); err != nil {
		return err
	}
	return nil
}