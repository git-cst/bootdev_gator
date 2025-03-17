package commands

import (
	"fmt"

	"github.com/git-cst/bootdev_gator/internal/config"
)

func HandlerLogin(s *config.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("no user passed to the login: %v", cmd.Args)
	}

	userLogin := cmd.Args[0]

	s.Config.User = userLogin
	err := config.WriteConfig(s.Config)
	if err != nil {
		return fmt.Errorf("error while setting user")
	}

	fmt.Printf("Successfully set user %s\n", userLogin)

	return nil
}
