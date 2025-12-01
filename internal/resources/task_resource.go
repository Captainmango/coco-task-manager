package resources

import (
	"log/slog"

	"github.com/captainmango/coco-cron-parser/internal/config"
	"github.com/captainmango/coco-cron-parser/internal/crontab"
)

type TaskResource struct {
	crontabManager crontab.CrontabHandler
}

func CreateTaskResource(
	ctbeManager crontab.CrontabHandler,
) TaskResource {
	return TaskResource{
		crontabManager: ctbeManager,
	}
}

func (t TaskResource) GetAllCrontabEntries() ([]crontab.CrontabEntry, error) {
	entries, err := t.crontabManager.GetAllCrontabEntries()

	slog.Info("retrieved entries from cron file",
		slog.String("file_location", config.Config.CrontabFile),
	)

	if err != nil {
		return nil, err
	}

	return entries, nil
}
