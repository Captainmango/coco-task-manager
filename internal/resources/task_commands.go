package resources

import (
	"context"
	"log/slog"

	"github.com/urfave/cli/v3"
)

func createStartGameCommand(tR TaskResource) *cli.Command {
	return &cli.Command{
		Name: "start-game",
		Description: "Sends the start game message to the Dealer API",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "room_id",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			_ = tR
			slog.Info("testing this out", slog.String("room_id", c.StringArg("room_id")))
			return nil
		},
	}
}

func createScheduleCronCommand(tR TaskResource) *cli.Command {
	return &cli.Command{
		Name:        "schedule-task",
		Description: "Schedules a task to be run via the scheduler.",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "cron",
			},
			&cli.StringArg{
				Name: "task",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			cronString := c.StringArg("cron")
			taskString := c.StringArg("task")

			if cronString == "" || taskString == "" {
				return cli.Exit("cron and task arguments are required", 1)
			}

			_, err := tR.ScheduleTask(cronString, taskString)
			if err != nil {
				return err
			}

			slog.Info("Scheduled task",
				slog.String("cron", cronString),
				slog.String("task", taskString),
			)

			return nil
		},
	}
}

func init() {
	taskResource := CreateResources().TaskResource
	CommandRegistry.Register(createStartGameCommand(taskResource))
	CommandRegistry.Register(createScheduleCronCommand(taskResource))
}
