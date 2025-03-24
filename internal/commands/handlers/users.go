package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/git-cst/bootdev_gator/internal/commands"
	"github.com/git-cst/bootdev_gator/internal/config"
	"github.com/git-cst/bootdev_gator/internal/database"
	"github.com/lib/pq"
)

func HandlerLogin(s *config.State, cmd commands.Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("no user passed to the login handler: %v", cmd.Args)
	}

	loginUserRequest := cmd.Args[0]
	databaseUser, err := s.Db.GetUser(context.Background(), loginUserRequest)
	if err != nil {
		// Check if this is a no data found error
		if err == sql.ErrNoRows {
			fmt.Println("User does not exist")
			os.Exit(1)
		}
		// Handle other errors
		fmt.Printf("Error retrieving user: %s\nError was: %v\n", loginUserRequest, err)
	}

	setConfigUser(s, databaseUser.Name)

	return nil
}

func HandlerRegister(s *config.State, cmd commands.Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("no user passed to the register handler: %v", cmd.Args)
	}

	registerUser := cmd.Args[0]
	params := database.CreateUserParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      registerUser,
	}

	createdUser, err := s.Db.CreateUser(context.Background(), params)
	if err != nil {
		// Check if this is a unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation code
				fmt.Println("User already exists")
				os.Exit(1)
			}
		}
		// Handle other errors
		fmt.Printf("Error creating user: %v\n", err)
		os.Exit(1)
	}

	setConfigUser(s, createdUser.Name)
	fmt.Printf("Created user: %v\n", createdUser)

	return nil
}

func HandlerUsers(s *config.State, cmd commands.Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		// Check if this is a no data found error
		if err == sql.ErrNoRows {
			fmt.Println("User does not exist")
			os.Exit(1)
		}
		// Handle other errors
		fmt.Printf("Error retrieving users\nError was: %v\n", err)
	}

	for _, user := range users {
		message := fmt.Sprintf("* %s", user.Name)
		if user.Name == s.Config.User {
			message += " (current)"
		}
		fmt.Printf("%s\n", message)
	}

	return nil
}

func HandlerReset(s *config.State, cmd commands.Command) error {
	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		fmt.Printf("Error resetting user table: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully reset user table")

	return nil
}
