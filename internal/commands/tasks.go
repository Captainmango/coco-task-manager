package commands

import "github.com/urfave/cli/v3"

func createDeleteCronCommand() *cli.Command {
	return &cli.Command{}
}

func createScheduleCronCommand() *cli.Command {
	return &cli.Command{}
}

func init () {
	Registry.Register(createDeleteCronCommand())
	Registry.Register(createScheduleCronCommand())
}