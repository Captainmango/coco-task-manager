package coco_cli

import "github.com/urfave/cli/v3"

var Registry *RegistryContainer = &RegistryContainer{}

type CommandFinder interface {
	Find(string) (any, error)
	All() []*cli.Command
}

type RegistryContainer struct {
	Commands []*cli.Command // should be a slice of cli.Commands from urfav/cli
}

func (r *RegistryContainer) Register(cmd *cli.Command) {
	r.Commands = append(r.Commands, cmd)
}

func (r *RegistryContainer) Find(key string) (any, error) {
	var c any

	for _, cmd := range r.Commands {
		if cmd.Name == key {
			c = cmd
		}
	}

	return c, nil
}

func (r *RegistryContainer) All() []*cli.Command {
	return r.Commands
}