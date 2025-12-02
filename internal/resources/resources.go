package resources

import (
	"github.com/captainmango/coco-cron-parser/internal/commands"
	"github.com/captainmango/coco-cron-parser/internal/crontab"
)

type Resources struct {
	TaskResource TaskResource
}

func CreateResources() Resources {
	return Resources{
		CreateTaskResource(
			crontab.CrontabManager{},
			commands.Registry,
		), // Maybe need to DI this for int testing?
	}
}
