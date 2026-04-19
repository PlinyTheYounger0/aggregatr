package main

import (
	"context"
	"fmt"
)

func handlerLogin(state *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Too few arguements for the login function.")
	}

	userName := cmd.Args[0]
	_, err := state.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("Unable to find %s: %w", userName, err)
	}

	err = state.cfg.SetUser(userName)
	if err != nil {
		return fmt.Errorf("Unable to set user %s: %w", userName, err)
	}

	fmt.Printf("User has been set to %s.\n", userName)
	return nil
}
