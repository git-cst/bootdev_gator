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
	if len(cmd.Args) < 2 {
		return fmt.Errorf("incorrect number of arguments passed, expected 2 (name of feed & url): %v", cmd.Args)
	}

	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]
	if !isValidURL(feedUrl) { // See helpers.go for implementation
		return fmt.Errorf("second argument should be the url, the passed argument does not fulfill url schema: %v", cmd.Args[1])
	}

	ctx := context.Background()
	// Check if feed already exists
	var feedId uuid.UUID
	feed, err := s.Db.GetFeedByUrl(ctx, feedUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Feed doesn't exist so create the feed
			params := database.CreateFeedParams{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      feedName,
				Url:       feedUrl,
				UserID:    user.ID,
			}
			createdFeed, err := s.Db.CreateFeed(ctx, params)
			if err != nil {
				return err
			}
			feedId = createdFeed.ID
			fmt.Printf("Feed using url %s did not exist.\nCreated: %v.\n", feedUrl, createdFeed)
		} else {
			// Handle other errors
			return err
		}
	} else {
		feedId = feed.ID
	}

	return followFeed(ctx, s, feedId, user)
}

// middleware auth handles user
func HandlerFollowFeed(s *config.State, cmd commands.Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("no url passed to the follow feed handler: %v", cmd.Args)
	}

	ctx := context.Background()
	feedUrl := cmd.Args[0]
	feed, err := s.Db.GetFeedByUrl(ctx, feedUrl)
	if err != nil {
		// Check if this is a no data found error
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("No feeds registered for url: %s", feedUrl)
			return err
		}
		// Handle other errors
		return err
	}

	return followFeed(ctx, s, feed.ID, user)
}

func HandlerGetFeeds(s *config.State, cmd commands.Command) error {
	ctx := context.Background()
	feedData, err := s.Db.GetFeeds(ctx)
	if err != nil {
		// Check if this is a no data found error
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("No feeds found.")
			return nil
		}
		return err
	}

	jsonData, err := json.MarshalIndent(feedData, "", "	")
	if err != nil {
		return fmt.Errorf("error whilst prettifying feed query data: %v", err)
	}

	fmt.Println(string(jsonData))

	return nil
}

// middleware auth handles user
func HandlerGetFollowing(s *config.State, cmd commands.Command, user database.User) error {
	ctx := context.Background()

	feedFollows, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		// Check if this is a no data found error
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("User isn't following any feeds yet: %s\n", user.Name)
			return nil
		}
		// Handle other errors
		return err
	}

	fmt.Printf("User %s is following:\n", user.Name)
	for _, follow := range feedFollows {
		fmt.Printf(" - %s", follow.FeedName)
	}

	return nil
}

// middleware auth handles user
func HandlerUnfollow(s *config.State, cmd commands.Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("no url passed to command: %v", cmd.Args)
	}

	ctx := context.Background()
	feedUrl := cmd.Args[0]
	feed, err := s.Db.GetFeedByUrl(ctx, feedUrl)
	if err != nil {
		// Check if no data found error
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no feed has that url, %v", feedUrl)
		}
		// Handle other errors
		return err
	}

	params := database.RemoveFollowForUserParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	s.Db.RemoveFollowForUser(ctx, params)

	return nil
}
