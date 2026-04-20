package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/PlinyTheYounger0/aggregatr/internal/config"
	"github.com/PlinyTheYounger0/aggregatr/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("Error Opening DB: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmds := commands{make(map[string]func(state *state, cmd command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

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
