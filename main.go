package main

import (
	"fmt"
	"os"

	"github.com/git-cst/bootdev_gator/internal/commands"
	"github.com/git-cst/bootdev_gator/internal/config"
)

func main() {
	configFile, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	state := config.State{
		Config: &configFile,
	}
	cmds := commands.Commands{
		HandlerFunctions: make(map[string]func(*config.State, commands.Command) error),
	}

	cmds.Register("login", commands.HandlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Usage: program_name command [args...]")
		os.Exit(1)
	}

	cmd := commands.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = cmds.Run(&state, cmd)
	if err != nil {
		fmt.Println("Error running command:", err)
		os.Exit(1)
	}
}
