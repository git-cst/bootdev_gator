package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/git-cst/bootdev_gator/internal/commands"
	"github.com/git-cst/bootdev_gator/internal/config"
	"github.com/git-cst/bootdev_gator/internal/database"
	"github.com/google/uuid"
)

// middleware auth handles user
func HandlerAddFeed(s *config.State, cmd commands.Command, user database.User) error {
	s.LogInfo("User %s attempting to add feed: args=%v", user.Name, cmd.Args)
	if len(cmd.Args) < 2 {
		errMsg := fmt.Sprintf("incorrect number of arguments passed, expected 2 (name of feed & url): %v", cmd.Args)
		s.LogError("Failed to add feed: %s", errMsg)
		return fmt.Errorf("incorrect number of arguments passed, expected 2 (name of feed & url): %v", cmd.Args)
	}

	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]
	if !isValidURL(feedUrl) { // See helpers.go for implementation
		errMsg := fmt.Sprintf("second argument should be the url, the passed argument does not fulfill url schema: %v", cmd.Args[1])
		s.LogError("Failed to add feed: %s", errMsg)
		return fmt.Errorf("second argument should be the url, the passed argument does not fulfill url schema: %v", cmd.Args[1])
	}

	ctx := context.Background()
	// Check if feed already exists
	var feedId uuid.UUID
	feed, err := s.Db.GetFeedByUrl(ctx, feedUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Feed doesn't exist so create the feed
			s.LogInfo("Creating new feed '%s' with URL '%s'", feedName, feedUrl)
			params := database.CreateFeedParams{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      feedName,
				Url:       feedUrl,
				UserID:    user.ID,
			}
			createdFeed, err := s.Db.CreateFeed(ctx, params)
			if err != nil {
				s.LogError("Failed to create feed in database: %v", err)
				return err
			}
			feedId = createdFeed.ID
			s.LogInfo("Successfully created feed: id=%s, name=%s", createdFeed.ID, feedName)
		} else {
			// Handle other errors
			s.LogError("Error checking for existing feed: %v", err)
			return err
		}
	} else {
		feedId = feed.ID
		s.LogInfo("Feed already exists, using existing feed: id=%s, name=%s", feed.ID, feed.Name)
	}

	err = followFeed(ctx, s, feedId, user)
	if err != nil {
		s.LogError("Failed to follow feed: %v", err)
		return err
	}

	s.LogInfo("User %s successfully added and followed feed: id=%s, name=%s", user.Name, feedId, feedName)
	return nil
}

// middleware auth handles user
func HandlerFollowFeed(s *config.State, cmd commands.Command, user database.User) error {
	s.LogInfo("User %s attempting to follow feed: args=%v", user.Name, cmd.Args)

	if len(cmd.Args) < 1 {
		s.LogError("no url passed to the follow feed handler: %v", cmd.Args)
		return fmt.Errorf("no url passed to the follow feed handler: %v", cmd.Args)
	}

	ctx := context.Background()
	feedUrl := cmd.Args[0]
	feed, err := s.Db.GetFeedByUrl(ctx, feedUrl)
	if err != nil {
		// Check if this is a no data found error
		if errors.Is(err, sql.ErrNoRows) {
			s.LogError("No feeds registered for url: %s", feedUrl)
			return err
		}
		// Handle other errors
		s.LogError("Error while following feed: %v", err)
		return err
	}

	err = followFeed(ctx, s, feed.ID, user)
	if err != nil {
		s.LogError("Failed to follow feed: %v", err)
		return err
	}

	s.LogInfo("User %s successfully followed feed: id=%s, name=%s", user.Name, feed.ID, feed.Name)
	return nil
}

func HandlerGetFeeds(s *config.State, cmd commands.Command) error {
	s.LogInfo("Getting all feeds")

	ctx := context.Background()
	feedData, err := s.Db.GetFeeds(ctx)
	if err != nil {
		// Check if this is a no data found error
		if errors.Is(err, sql.ErrNoRows) {
			s.LogInfo("No feeds found in database")
			return nil
		}
		s.LogError("Failed to query feeds: %v", err)
		return err
	}

	jsonData, err := json.MarshalIndent(feedData, "", "	")
	if err != nil {
		s.LogError("Error marshaling feed data to JSON: %v", err)
		return fmt.Errorf("error whilst prettifying feed query data: %v", err)
	}

	s.LogInfo("Successfully retrieved %d feeds", len(feedData))
	fmt.Println(string(jsonData))

	return nil
}

// middleware auth handles user
func HandlerGetFollowing(s *config.State, cmd commands.Command, user database.User) error {
	s.LogInfo("Getting feeds that %s (%v) is following", user.Name, user.ID)
	ctx := context.Background()

	feedFollows, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		// Check if this is a no data found error
		if errors.Is(err, sql.ErrNoRows) {
			s.LogInfo("%s (%v) is not following any feeds yet", user.Name, user.ID)
			fmt.Printf("User isn't following any feeds yet: %s\n", user.Name)
			return nil
		}
		// Handle other errors
		return err
	}

	for _, follow := range feedFollows {
		s.LogInfo("%s (%v) is following: %s", user.Name, user.ID, follow.FeedName)
	}

	s.LogInfo("Successfully retrieved follows for %s (%v) user", user.Name, user.ID)
	return nil
}

// middleware auth handles user
func HandlerUnfollow(s *config.State, cmd commands.Command, user database.User) error {
	if len(cmd.Args) < 1 {
		s.LogError("no url passed to the unfollow feed handler: %v", cmd.Args)
		return fmt.Errorf("no url passed to command: %v", cmd.Args)
	}

	s.LogInfo("Unfollowing feed with url: %s", cmd.Args[0])
	ctx := context.Background()
	feedUrl := cmd.Args[0]
	feed, err := s.Db.GetFeedByUrl(ctx, feedUrl)
	if err != nil {
		// Check if no data found error
		if errors.Is(err, sql.ErrNoRows) {
			s.LogInfo("No feed registered with url: %s", feedUrl)
			return fmt.Errorf("no feed has that url, %v", feedUrl)
		}
		// Handle other errors
		s.LogError("Error while retrieving feed registered with url: %v", err)
		return err
	}

	params := database.RemoveFollowForUserParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.Db.RemoveFollowForUser(ctx, params)
	if err != nil {
		s.LogError("Error while executing query to unfollow feed registered with url: %v", err)
		return err
	}

	s.LogInfo("Successfully unfollowed feed with url: %s", feedUrl)
	return nil
}
