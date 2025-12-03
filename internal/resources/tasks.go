package resources

import (
	"context"
	"log/slog"

	coco_cli "github.com/captainmango/coco-cron-parser/internal/cli"
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
			tR.ScheduleTask()
			slog.Info("do a thing")
			return nil
		},
	}
}

func init () {
	taskResource := CreateResources().TaskResource
	coco_cli.Registry.Register(createStartGameCommand(taskResource))
	coco_cli.Registry.Register(createScheduleCronCommand(taskResource))
}