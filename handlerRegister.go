package main

import (
	"context"
	"fmt"
	"time"

	"github.com/PlinyTheYounger0/aggregatr/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(state *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %v <name>", cmd.Name)
	}

	name := cmd.Args[0]
	user, err := state.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("Error Creating User: %w", err)
	}

	err = state.cfg.SetUser(user.Name)
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
