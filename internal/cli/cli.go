package cli

import (
	"fmt"
	"os"

	"github.com/donaldnguyen99/gator/internal/config"
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
	cli.state = &state{
		config: config,
	}

	cli.commands.register("login", handlerLogin)

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