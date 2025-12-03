package resources

import (
	coco_cli "github.com/captainmango/coco-cron-parser/internal/cli"
	"github.com/captainmango/coco-cron-parser/internal/crontab"
)

type Resources struct {
	TaskResource TaskResource
}

func CreateResources() Resources {
	return Resources{
		CreateTaskResource(
			crontab.CrontabManager{},
			coco_cli.Registry,
		),
	}
}
