package config

import (
	"log"

	"github.com/git-cst/bootdev_gator/internal/database"
)

type State struct {
	Config *Config
	Db     *database.Queries
	Logger *log.Logger
}
