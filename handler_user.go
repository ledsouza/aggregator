package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ledsouza/aggregator/internal/database"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete users")
	}

	fmt.Println("The database was succesfully reseted!")
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is required")
	}

	username := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("couldn't get the user: %w", err)
	}

	err = s.config.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set user: %w", err)
	}

	fmt.Printf("the user %s has been set \n", username)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("name is required")
	}

	username := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		return fmt.Errorf("couldn't create an user: %w", err)
	}

	err = s.config.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set user: %w", err)
	}

	fmt.Printf("the user %s has been set \n", username)
	printUser(user)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get users: %w", err)
	}

	for _, user := range users {
		if s.config.CurrentUserName == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
