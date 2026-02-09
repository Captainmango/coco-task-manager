package resources

import (
	"log/slog"

	"github.com/captainmango/coco-cron-parser/internal/config"
	"github.com/captainmango/coco-cron-parser/internal/crontab"
	"github.com/captainmango/coco-cron-parser/internal/msq"
)

type Resources struct {
	TaskResource TaskResource
}

func CreateResources() Resources {
	queueHandler, err := msq.NewRabbitMQHandler(msq.WithConnStr(config.Config.RabbitMQHost))

	if err != nil {
		slog.Error(err.Error())
	}

	return Resources{
		CreateTaskResource(
			crontab.CrontabManager{},
			queueHandler,
		), // Maybe need to DI this for int testing?
	}
}
