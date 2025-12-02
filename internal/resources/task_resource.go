package resources

import (
	"log/slog"

	"github.com/captainmango/coco-cron-parser/internal/commands"
	"github.com/captainmango/coco-cron-parser/internal/config"
	"github.com/captainmango/coco-cron-parser/internal/crontab"
	"github.com/google/uuid"
)

type TaskResource struct {
	crontabManager crontab.CrontabHandler
	commandRegistry commands.CommandFinder
}

func CreateTaskResource(
	ctbeManager crontab.CrontabHandler,
	cmdRegistry commands.CommandFinder,
) TaskResource {
	return TaskResource{
		crontabManager: ctbeManager,
		commandRegistry: cmdRegistry,
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

func (t TaskResource) GetAllAvailableCommands() []any {
	return t.commandRegistry.All()
}

func (t TaskResource) ScheduleTask() error {
	return nil
}

func (t TaskResource) GetTaskByID(id uuid.UUID) (any, error) {
	ctbE, err := t.crontabManager.GetCrontabEntryByID(id)

	if err != nil {
		return crontab.CrontabEntry{}, nil
	}

	return ctbE, nil
}

func (t TaskResource) RemoveTaskByID(id uuid.UUID) error {
	err := t.crontabManager.RemoveCrontabEntryByID(id)
	return err
}
