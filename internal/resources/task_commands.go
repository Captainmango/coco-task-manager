package resources

import (
	"context"

	"github.com/urfave/cli/v3"
)

func createStartGameCommand(tR TaskResource) *cli.Command {
	return &cli.Command{
		Name: "start-game",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "room_id",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			_ = tR
			// The thing that will post to rabbitmq
			return nil
		},
	}
}

func createScheduleCronCommand(tR TaskResource) *cli.Command {
	return &cli.Command{
		Name: "schedule-task",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "cron",
			},
			&cli.StringArg{
				Name: "task",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			// cmds := []string{}
			// slog.Info("do a thing", slog.String("cmd_name", cmds[0].Name))
			return nil
		},
	}
}

func init() {
	taskResource := CreateResources().TaskResource
	CommandRegistry.Register(createStartGameCommand(taskResource))
	CommandRegistry.Register(createScheduleCronCommand(taskResource))
}
