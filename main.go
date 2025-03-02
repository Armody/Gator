package main

import (
	"log"
	"os"

	"github.com/Armody/Gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	programState := &state{&cfg}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

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
