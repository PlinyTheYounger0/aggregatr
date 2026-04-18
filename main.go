package main

import (
	"log"
	"os"

	"github.com/PlinyTheYounger0/aggregatr/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	programState := &state{
		cfg: &cfg,
	}

	cmds := commands{make(map[string]func(state *state, cmd command) error)}
	cmds.register("login", handlerLogin)

	input := os.Args
	if len(input) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmd := command{
		Name: input[1],
		Args: input[2:],
	}

	err = cmds.run(programState, command{Name: cmd.Name, Args: cmd.Args})
	if err != nil {
		log.Fatalf("Error running %s: %v", cmd.Name, err)
	}
}
