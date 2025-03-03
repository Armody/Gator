package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Armody/Gator/internal/config"
	"github.com/Armody/Gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

const dbURL = "postgres://postgres:1234@localhost:5432/gator"

func main() {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	programState := &state{dbQueries, &cfg}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	args := os.Args
	if len(args) < 3 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	cmdName, cmdArgs := args[1], args[2:]
	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}

	err = cmds.run(programState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
