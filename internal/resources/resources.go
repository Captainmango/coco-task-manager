package resources

import (
	"github.com/captainmango/coco-cron-parser/internal/config"
	"github.com/captainmango/coco-cron-parser/internal/crontab"
	"github.com/captainmango/coco-cron-parser/internal/msq"
)

type Resources struct {
	TaskResource TaskResource
}

func CreateResources() Resources {
	return Resources{
		CreateTaskResource(
			crontab.CrontabManager{},
			msq.NewRabbitMQHandler(
				msq.WithConnStr(
					config.Config.RabbitMQHost,
				),
			),
		), // Maybe need to DI this for int testing?
	}
}
