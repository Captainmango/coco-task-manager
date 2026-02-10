package coco_cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/urfave/cli/v3"

	"github.com/captainmango/coco-cron-parser/internal/msq"
	"github.com/captainmango/coco-cron-parser/internal/resources"
)

func createStartGameCommand(tR resources.TaskResource) *cli.Command {
	return &cli.Command{
		Name:        "start-game",
		Description: "Sends the start game message to the Dealer API",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "room_id",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			roomIdString := c.StringArg("room_id")
			if roomIdString == "" {
				return cli.Exit("room_id argument is required", 1)
			}

			startGamePayload := msq.StartGamePayload{
				RoomId: roomIdString,
			}

			if err := tR.PushStartGameMessage(startGamePayload); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return nil
		},
	}
}

func createPullMessagesCommand(tR resources.TaskResource) *cli.Command {
	return &cli.Command{
		Name:        "pull-messages",
		Description: "Pulls messages from the given topic. Used primarily for debugging",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "topic",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			topicString := c.StringArg("topic")
			if topicString == "" {
				return cli.Exit("topic argument is required", 1)
			}

			forever := make(chan struct{})
			// nolint:errcheck // This is for debugging purposes.
			go func() error {
				err := tR.ProcessMessages(ctx, "coco_tasks.start_game", func(msg amqp.Delivery) error {
					var t struct {
						RoomId string `json:"room_id"`
					}

					json.NewDecoder(bytes.NewReader(msg.Body)).Decode(&t)

					slog.Info("Got Json", slog.String("payload", t.RoomId))
					msg.Ack(false)

					return nil
				})
				if err != nil {
					return cli.Exit(err.Error(), 1)
				}

				return nil
			}()

			fmt.Println("Waiting for messages. CTRL + C to terminate")
			<-forever

			return nil
		},
	}
}

func createScheduleCronCommand(tR resources.TaskResource) *cli.Command {
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
	taskResource := resources.CreateResources().TaskResource
	CommandRegistry.Register(createStartGameCommand(taskResource))
	CommandRegistry.Register(createScheduleCronCommand(taskResource))
	CommandRegistry.Register(createPullMessagesCommand(taskResource))
}
