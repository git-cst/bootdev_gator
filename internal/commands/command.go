package commands

import (
	"fmt"

	"github.com/git-cst/bootdev_gator/internal/config"
)

type Commands struct {
	HandlerFunctions map[string]func(*config.State, Command) error
}

type Command struct {
	Name string
	Args []string
}

func (c *Commands) Register(name string, f func(*config.State, Command) error) {
	if c.HandlerFunctions == nil {
		c.HandlerFunctions = make(map[string]func(*config.State, Command) error)
	}
	c.HandlerFunctions[name] = f
}

func (c *Commands) Run(s *config.State, cmd Command) error {
	function, exists := c.HandlerFunctions[cmd.Name]
	if !exists {
		return fmt.Errorf("no handler function by the name of %s registered", cmd.Name)
	}

	err := function(s, cmd)
	if err != nil {
		return fmt.Errorf("error while running function %s\nError was %v", cmd.Name, err)
	}
	return nil
}
