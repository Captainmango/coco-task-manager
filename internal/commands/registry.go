package commands

var Registry *registry = &registry{}

type registry struct {
	Commands []any // should be a slice of cli.Commands from urfav/cli
}

func (r *registry) Register(cmd any) {
	r.Commands = append(r.Commands, cmd)
}

