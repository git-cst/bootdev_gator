package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/git-cst/bootdev_gator/internal/commands"
	"github.com/git-cst/bootdev_gator/internal/commands/handlers"
	"github.com/git-cst/bootdev_gator/internal/config"
	"github.com/git-cst/bootdev_gator/internal/database"
	"github.com/git-cst/bootdev_gator/internal/middleware"
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

	// On Unix-like systems (Linux, macOS)
	username := os.Getenv("USER")

	// On Windows
	if username == "" {
		username = os.Getenv("USERNAME")
	}

	clientSetup := config.ClientOptions{
		Timeout:     time.Duration(time.Second * 60),
		UserAgent:   fmt.Sprintf("%s_gator", username),
		ContentType: "application/rss+xml, application/atom+xml, application/xml, text/xml",
		Headers:     make(map[string]string),
	}

	logger := config.CreateLogger()
	defer logger.Close()

	state := config.State{
		Config: &configFile,
		Db:     dbQueries,
		Logger: logger,
	}

	state.Config.Client = config.NewClient(clientSetup)

	cmds := commands.Commands{
		HandlerFunctions: make(map[string]func(*config.State, commands.Command) error),
	}

	// user related commands
	cmds.Register("login", "Log the specified user in. 1st argument is the user.", handlers.HandlerLogin)
	cmds.Register("register", "Register the user in the database.", handlers.HandlerRegister)
	cmds.Register("reset", "Reset the users table.", handlers.HandlerReset)
	cmds.Register("users", "Retrieve the available users in the database.", handlers.HandlerUsers)

	// rss feed related commands
	cmds.Register("feeds", "Retrieve the available feeds in the database.", handlers.HandlerGetFeeds)
	cmds.Register("addfeed", "Add a new feed to be fetched.", middleware.MiddlewareLoggedIn(handlers.HandlerAddFeed))
	cmds.Register("follow", "Follow a registered feed.", middleware.MiddlewareLoggedIn(handlers.HandlerFollowFeed))
	cmds.Register("following", "Retrieve what another specified user is following.", middleware.MiddlewareLoggedIn(handlers.HandlerGetFollowing))
	cmds.Register("unfollow", "Unfollow a feed.", middleware.MiddlewareLoggedIn(handlers.HandlerUnfollow))

	// service related commands
	cmds.Register("agg", "Start the aggregator service.", handlers.HandlerAgg)
	cmds.Register("browse", "Browse X feeds where X is the argument passed to the command.", middleware.MiddlewareLoggedIn(handlers.HandlerBrowse))

	// application commands
	cmds.Register("help", "Display the commands available to you.", cmds.ListCommands)

	state.LogDebug("Starting RSS aggregator application")

	if len(os.Args) < 2 {
		state.LogError("No command provided: Usage: program_name command [args...]")
		os.Exit(1)
	}

	cmd := commands.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if state.Config.User != "" {
		state.LogDebug("User '%s' executing command: '%s' with args: %v", state.Config.User, cmd.Name, cmd.Args)
	} else {
		state.LogDebug("Unauthenticated user executing command: '%s' with args: %v", cmd.Name, cmd.Args)
	}

	err = cmds.Run(&state, cmd)
	if err != nil {
		state.LogError("Error executing command '%s': %v", cmd.Name, err)
		os.Exit(1)
	}

	state.LogDebug("Command '%s' completed successfully", cmd.Name)
}
