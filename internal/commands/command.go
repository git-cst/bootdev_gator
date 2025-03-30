package commands

import (
	"fmt"

	"github.com/git-cst/bootdev_gator/internal/config"
)

type Commands struct {
	HandlerFunctions map[string]func(*config.State, Command) error
	CommandList      map[string]Command
}

type Command struct {
	Name        string
	Description string
	Args        []string
}

func (c *Command) GetName() string {
	return c.Name
}

func (c *Command) GetDescription() string {
	return c.Description
}

func (c *Command) GetArgs() []string {
	return c.Args
}

func (c *Commands) Register(name string, description string, f func(*config.State, Command) error) {
	if c.HandlerFunctions == nil {
		c.HandlerFunctions = make(map[string]func(*config.State, Command) error)
	}
	c.HandlerFunctions[name] = f

	// Store the command with its description
	command := Command{
		Name:        name,
		Description: description,
		Args:        []string{},
	}

	// We need to keep track of commands and their descriptions
	if c.CommandList == nil {
		c.CommandList = make(map[string]Command)
	}
	c.CommandList[name] = command
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

func (c *Commands) ListCommands(s *config.State, cmd Command) error {
	fmt.Println("Available commands:")
	fmt.Println("------------------")

	// You might want to sort the commands for consistent output
	for name, command := range c.CommandList {
		s.LogInfo("%-15s - %s", name, command.Description)
	}

	return nil
}
