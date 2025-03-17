package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/donaldnguyen99/gator/internal/database"
	"github.com/donaldnguyen99/gator/internal/rss"
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

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("get users requires 0 arguments")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error handling get users: %w", err)
	}

	for _, u := range users {
		if u.Name == s.config.CurrentUserName {
			fmt.Printf("* %v (current)\n", u.Name)
		} else {
			fmt.Printf("* %v\n", u.Name)
		}
	}
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

func handlerAggregateFeeds(s *state, cmd command) error {
	// TODO: may need to change behavior
	if len(cmd.args) == 0 {
		return fmt.Errorf("feed requires at least 1 argument")
	}

	feeds, err := rss.AggregateFeeds(context.Background(), cmd.args)
	if err != nil {
		return fmt.Errorf("error handling fetch feed: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("Feed %v has been fetched\n", feed)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("addfeed requires 2 arguments (name, url), got %v arguments", len(cmd.args))
	}

	name, url := cmd.args[0], cmd.args[1]

	if s.config.CurrentUserName == "" {
		return fmt.Errorf("no user logged in")
	}
	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("no user found: %w", err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error handling add feed: %w", err)
	}
	fmt.Printf("Feed %v, %v has been added\n", feed.Name, feed.Url)

	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeedsWithUserName(context.Background())
	if err != nil {
		return fmt.Errorf("error handling list feeds: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("Feed %v, %v, %v\n", feed.Name, feed.Url, feed.UserName)
	}
	return nil
}
