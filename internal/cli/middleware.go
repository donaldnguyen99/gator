package cli

import (
	"context"

	"github.com/donaldnguyen99/gator/internal/database"
)

func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error,
) func(s *state, cmd command) error {

	newFunc := func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
	return newFunc
}