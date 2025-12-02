package commands

var Registry *RegistryContainer = &RegistryContainer{}

type CommandFinder interface {
	Find(string) (any, error)
	All() []any
}

type RegistryContainer struct {
	Commands []any // should be a slice of cli.Commands from urfav/cli
}

func (r *RegistryContainer) Register(cmd any) {
	r.Commands = append(r.Commands, cmd)
}

func (r *RegistryContainer) Find(key string) (any, error) {
	var c any

	for _, cmd := range r.Commands {
		if cmd == key {
			c = cmd
		}
	}

	return c, nil
}

func (r *RegistryContainer) All() []any {
	return r.Commands
}