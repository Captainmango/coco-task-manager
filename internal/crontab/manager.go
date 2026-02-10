package crontab

import (
	"errors"

	"github.com/google/uuid"
)

const cronFormat = "%s root /app/%s 2>&1 | tee -a /tmp/log # %s\n"

var (
	errCrontabFileNotSet = errors.New("crontab file not set")
)

type CrontabHandler interface {
	WriteCrontabEntries([]CrontabEntry) error
	GetAllCrontabEntries() ([]CrontabEntry, error)
	GetCrontabEntryByID(uuid.UUID) (CrontabEntry, error)
	RemoveCrontabEntryByID(uuid.UUID) error
}
