package resources

import (
	"log/slog"

	"github.com/captainmango/coco-cron-parser/internal/config"
	"github.com/captainmango/coco-cron-parser/internal/crontab"
	"github.com/captainmango/coco-cron-parser/internal/parser"
	"github.com/google/uuid"
	"github.com/urfave/cli/v3"
)

type TaskResource struct {
	crontabManager crontab.CrontabHandler
	commandRegistry CommandFinder
}

func CreateTaskResource(
	ctbeManager crontab.CrontabHandler,
	cmdRegistry CommandFinder,
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

func (t TaskResource) GetAllAvailableCommands() []*cli.Command {
	return t.commandRegistry.All()
}

func (t TaskResource) ScheduleTask(cron, task string) error {
	p, err := parser.NewParser(parser.WithInput(cron, true))

	if err != nil {
		return err
	}

	parsedExpr, err := p.Parse()
	if err != nil {
		return err
	}

	id, err := uuid.NewV7()

	if err != nil {
		return err
	}

	ctbEntry := crontab.CrontabEntry{
		ID: id,
		Cron: parsedExpr,
		Cmd: task,
	}

	if err = t.crontabManager.WriteCrontabEntries([]crontab.CrontabEntry{ctbEntry}); err != nil {
		return err
	}

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
