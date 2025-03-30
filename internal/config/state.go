package config

import (
	"github.com/git-cst/bootdev_gator/internal/database"
)

type State struct {
	Config      *Config
	Db          *database.Queries
	Logger      *LogInstance
	CurrentUser *database.User
}
