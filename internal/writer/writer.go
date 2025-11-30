package writer

import "github.com/captainmango/coco-cron-parser/internal/data"

const cronFormat = "%s root %s | tee /tmp/log"

func WriteCronToSchedule(cron *data.Cron, cmd string) {
	
}