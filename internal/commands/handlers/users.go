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
	s.LogInfo("Starting login process")
	if len(cmd.Args) < 1 {
		s.LogError("No user passed to the login handler: %v", cmd.Args)
		return fmt.Errorf("no user passed to the login handler: %v", cmd.Args)
	}

	loginUserRequest := cmd.Args[0]
	s.LogInfo("User attempting login: %s", loginUserRequest)
	databaseUser, err := s.Db.GetUser(context.Background(), loginUserRequest)
	if err != nil {
		// Check if this is a no data found error
		if err == sql.ErrNoRows {
			s.LogInfo("User %s does not exist", loginUserRequest)
			os.Exit(1)
		}
		// Handle other errors
		s.LogError("Error retrieving user %s from database, error was %v", loginUserRequest, err)
	}

	setConfigUser(s, databaseUser.Name)
	s.LogInfo("User %s (%v) successfully logged in", databaseUser.Name, databaseUser.ID)

	return nil
}

func HandlerRegister(s *config.State, cmd commands.Command) error {
	s.LogInfo("Starting user registration process")
	if len(cmd.Args) < 1 {
		s.LogError("No user passed to the register handler: %v", cmd.Args)
		return fmt.Errorf("no user passed to the register handler: %v", cmd.Args)
	}

	registerUser := cmd.Args[0]
	params := database.CreateUserParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      registerUser,
	}

	s.LogInfo("Attempting to create user: %s", registerUser)
	createdUser, err := s.Db.CreateUser(context.Background(), params)
	if err != nil {
		// Check if this is a unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation code
				s.LogInfo("User %s already exists", registerUser)
				os.Exit(1)
			}
		}
		// Handle other errors
		s.LogError("Error creating user %s in database, error was %v", registerUser, err)
		os.Exit(1)
	}

	setConfigUser(s, createdUser.Name)
	s.LogInfo("Successfully registered user %s (%v)", createdUser.Name, createdUser.ID)

	return nil
}

func HandlerUsers(s *config.State, cmd commands.Command) error {
	s.LogInfo("Retrieving users")
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		// Check if this is a no data found error
		if err == sql.ErrNoRows {
			s.LogError("No users in users table")
			os.Exit(1)
		}
		// Handle other errors
		s.LogError("Error retrieving users from database, error was %v", err)
	}

	for _, user := range users {
		message := fmt.Sprintf("* %s", user.Name)
		if user.Name == s.Config.User {
			message += " (current)"
		}

		s.LogInfo("%s\n", message)
	}

	s.LogInfo("Successfully retrieved users from database")
	return nil
}

func HandlerReset(s *config.State, cmd commands.Command) error {
	s.LogInfo("Starting reset users process")
	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		s.LogError("Error resetting user table: %v\n", err)
		os.Exit(1)
	}

	s.LogInfo("Successfully reset users table in database")

	return nil
}
