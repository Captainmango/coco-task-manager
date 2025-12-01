package commands

var Registry *RegistryContainer = &RegistryContainer{}

type RegistryContainer struct {
	Commands []any // should be a slice of cli.Commands from urfav/cli
}

func (r *RegistryContainer) Register(cmd any) {
	r.Commands = append(r.Commands, cmd)
}

