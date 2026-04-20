package main

import (
	"context"
	"fmt"
	"time"

	"github.com/PlinyTheYounger0/aggregatr/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %v <name>", cmd.Name)
	}

	name := cmd.Args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("Error Creating User: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("%s created successfully.\n", user.Name)
	printUser(user)

	return nil
}

func printUser(user database.User) {
	fmt.Printf("ID: %v\n", user.ID)
	fmt.Printf("Name: %v\n", user.Name)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: login <name>")
	}

	userName := cmd.Args[0]
	_, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("Unable to find %s: %w", userName, err)
	}

	err = s.cfg.SetUser(userName)
	if err != nil {
		return fmt.Errorf("Unable to set user %s: %w", userName, err)
	}

	fmt.Printf("User has been set to %s.\n", userName)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error Fetching Users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("Error Reseting DB: %w", err)
	}

	return nil
}
