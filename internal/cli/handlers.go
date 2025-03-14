package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/donaldnguyen99/gator/internal/database"
	"github.com/google/uuid"
)


type handler func(*state, command) error

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login requires 1 argument (username)")
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("error handling login: %w", err)
	}
	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("error handling login: %w", err)
	}
	fmt.Printf("The user %s has been set\n", cmd.args[0])
	return nil
}

func handlerRegisterUser(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login requires 1 argument (username)")
	}
	u, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("error handling register user: %w", err)
	}

	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("error handling register user: %w", err)
	}
	fmt.Printf("The user %v has been created and set\n", u)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("reset requires 0 arguments")
	}
	if err := s.config.SetUser(""); err != nil {
		return fmt.Errorf("error handling reset: %w", err)
	}
	if err := s.db.DeleteAllUsers(context.Background()); err != nil {
		return fmt.Errorf("error handling reset: %w", err)
	}
	fmt.Println("The database has been reset")
	return nil
}