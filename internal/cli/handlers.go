package cli

import (
	"fmt"
)


type handler func(*state, command) error

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login requires 1 argument (username)")
	}
	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("error handling login: %w", err)
	}
	fmt.Printf("The user %s has been set\n", cmd.args[0])
	return nil
}