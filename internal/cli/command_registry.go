package coco_cli

import (
	"errors"

	"github.com/urfave/cli/v3"
)

var CommandRegistry *RegistryContainer = &RegistryContainer{}

type CommandFinder interface {
	Find(string) (*cli.Command, error)
	All() []*cli.Command
}

type RegistryContainer struct {
	Commands []*cli.Command // should be a slice of cli.Commands from urfav/cli
}

func (r *RegistryContainer) Register(cmd *cli.Command) {
	r.Commands = append(r.Commands, cmd)
}

func (r *RegistryContainer) Find(key string) (*cli.Command, error) {
	var c *cli.Command

	for _, cmd := range r.Commands {
		if cmd.Name == key {
			c = cmd
		}
	}

	if c == nil {
		return nil, errors.New("command not found")
	}

	return c, nil
}

func (r *RegistryContainer) All() []*cli.Command {
	return r.Commands
}
