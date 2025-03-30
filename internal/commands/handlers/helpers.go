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

func parseTimeString(timeString string) (time.Time, error) {
	// Try several common RSS date formats
	formats := []string{
		"Mon, 02 Jan 2006 15:04:05 -0700", // RFC1123 with timezone
		"Mon, 02 Jan 2006 15:04:05 MST",   // RFC1123 with timezone abbreviation
		"2006-01-02T15:04:05-07:00",       // ISO8601/RFC3339
		"2006-01-02T15:04:05Z",            // ISO8601/RFC3339 UTC
		"2006-01-02 15:04:05 -0700",       // Another common format
		"02 Jan 2006 15:04:05 -0700",      // Another variation
		// Add more formats as you encounter them
	}

	var firstErr error
	for _, format := range formats {
		publishedAt, err := time.Parse(format, timeString)
		if err == nil {
			return publishedAt, nil
		}
		if firstErr == nil {
			firstErr = err
		}
	}

	// If we got here, none of the formats worked
	return time.Time{}, fmt.Errorf("could not parse time '%s': %v", timeString, firstErr)
}
