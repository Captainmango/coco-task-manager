package resources

import (
	"github.com/captainmango/coco-cron-parser/internal/crontab"
)

type Resources struct {
	TaskResource TaskResource
}

func CreateResources() Resources {
	return Resources{
		CreateTaskResource(
			crontab.CrontabManager{},
		), // Maybe need to DI this for int testing?
	}
}
