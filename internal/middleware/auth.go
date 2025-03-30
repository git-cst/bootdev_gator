package middleware

import (
	"context"

	"github.com/git-cst/bootdev_gator/internal/commands"
	"github.com/git-cst/bootdev_gator/internal/config"
	"github.com/git-cst/bootdev_gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *config.State, cmd commands.Command, user database.User) error) func(*config.State, commands.Command) error {
	return func(s *config.State, cmd commands.Command) error {
		// Check if user is already cached and username matches config user (which is the logged in user)
		if s.CurrentUser != nil && s.CurrentUser.Name == s.Config.User {
			s.LogDebug("Using cached user: %s", s.CurrentUser.Name)
			return handler(s, cmd, *s.CurrentUser)
		}

		// Otherwise retrieve from database
		user, err := s.Db.GetUser(context.Background(), s.Config.User)
		if err != nil {
			s.LogError("Tried to retrieve user %s from database, failed while doing so err: %v", s.Config.User, err)
			return err
		}

		// Cache for future use
		s.CurrentUser = &user
		s.LogDebug("Cached user: %s", user.Name)

		return handler(s, cmd, user)
	}
}
