package main

import (
	"fmt"
)

func handlerLogin(state *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Too few arguements for the login function.")
	}

	userName := cmd.Args[0]
	state.cfg.SetUser(userName)

	fmt.Printf("User has been set to %s.\n", userName)
	return nil
}
