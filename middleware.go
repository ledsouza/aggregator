package main

import (
	"context"
	"fmt"

	"github.com/ledsouza/aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("couldn't get current user: %w", err)
		}

		err = handler(s, c, user)
		if err != nil {
			return fmt.Errorf("couldn't execute handler: %w", err)
		}

		return nil
	}
}
