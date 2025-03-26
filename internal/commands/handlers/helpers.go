package handlers

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/git-cst/bootdev_gator/internal/config"
	"github.com/git-cst/bootdev_gator/internal/database"
	"github.com/google/uuid"
)

// Used in users.go
func setConfigUser(s *config.State, u string) error {
	s.Config.User = u
	err := config.WriteConfig(s.Config)
	if err != nil {
		return fmt.Errorf("error while setting user: %s", u)
	}

	fmt.Printf("Successfully set user %s\n", u)

	return nil
}

// Used in feeds.go
func isValidURL(str string) bool {
	// Parse the string into a URL structure
	u, err := url.Parse(str)

	// Check if there was an error during parsing
	if err != nil {
		return false
	}

	// A valid URL should have a scheme (like http, https) and a host
	return u.Scheme != "" && u.Host != ""
}

// Used in feeds.go
func followFeed(ctx context.Context, s *config.State, feedId uuid.UUID, user database.User) error {
	params := database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedId,
	}

	follow, err := s.Db.CreateFeedFollow(ctx, params)
	if err != nil {
		return fmt.Errorf("error whilst creating feed follow: %v", err)
	}

	fmt.Printf("Feed: %s has been followed by user: %s\n", follow.FeedName, follow.Username)
	return nil
}

func parseTimeString(timeString string) time.Time {
	return time.Time{}
}
