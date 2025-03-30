package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/git-cst/bootdev_gator/internal/commands"
	"github.com/git-cst/bootdev_gator/internal/config"
	"github.com/git-cst/bootdev_gator/internal/database"
)

const (
	ColorReset = "\033[0m"
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
)

// middleware auth handles user
func HandlerBrowse(s *config.State, cmd commands.Command, user database.User) error {
	s.LogDebug("User %s requesting to see posts: args=%v", user.Name, cmd.Args)
	var numPosts int32
	if len(cmd.Args) < 1 {
		numPosts = 2
	} else {
		num, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			s.LogError("Could not convert %v to an integer", cmd.Args[0])
		}
		numPosts = int32(num)
	}

	ctx := context.Background()
	postParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  numPosts,
	}

	s.LogInfo("Browsing %d number of posts:", numPosts)
	posts, err := s.Db.GetPostsForUser(ctx, postParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No posts retrieved
			s.LogInfo("No posts retrieved. Check if you are following any feeds.")
			return nil
		}
		// Handle other errors
		s.LogError(fmt.Sprintf("%v", err))
	}

	for i := int32(0); i < numPosts; i++ {
		post := posts[i]
		s.LogInfo(ColorGreen+"Title:"+ColorReset+" %v | "+ColorGreen+"Link:"+ColorReset+" %v ("+ColorGreen+"Published:"+ColorReset+" %v)", post.Title, post.Url, post.PublishedAt)
	}

	return nil
}
