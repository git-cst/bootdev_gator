package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/git-cst/bootdev_gator/internal/commands"
	"github.com/git-cst/bootdev_gator/internal/config"
	"github.com/git-cst/bootdev_gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	configFile, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", configFile.DbURL)
	if err != nil {
		fmt.Println("Error opening connection to DB:", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	state := config.State{
		Config: &configFile,
		Db:     dbQueries,
	}

	cmds := commands.Commands{
		HandlerFunctions: make(map[string]func(*config.State, commands.Command) error),
	}

	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)

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
