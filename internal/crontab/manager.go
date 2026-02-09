package crontab

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"

	"github.com/captainmango/coco-cron-parser/internal/config"
	"github.com/captainmango/coco-cron-parser/internal/parser"
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

type CrontabManager struct{
	mu sync.Mutex
}

func NewCrontabManager() *CrontabManager {
	return &CrontabManager{}
}

// Sets the cron printing mode to RAW_EXPRESSION and writes to the configured crontab file
func (cM *CrontabManager) WriteCrontabEntries(crontabs []CrontabEntry) error {
	err := cM.withCrontab(func(f *os.File) error {
		for _, ctbE := range crontabs {
			ctbE.Cron.PrintingMode = parser.RAW_EXPRESSION

			_, err := fmt.Fprintf(f, cronFormat, ctbE.Cron, ctbE.Cmd, ctbE.ID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (cM *CrontabManager) GetAllCrontabEntries() ([]CrontabEntry, error) {
	var out []CrontabEntry

	err := cM.withCrontab(func(f *os.File) error {
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			line := scanner.Text()

			ctbE, err := NewCrontabEntryFromString(line)
			if err != nil {
				return err
			}

			out = append(out, ctbE)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (cM *CrontabManager) GetCrontabEntryByID(id uuid.UUID) (CrontabEntry, error) {
	var ctbE CrontabEntry
	allEntries, err := cM.GetAllCrontabEntries()
	if err != nil {
		return ctbE, nil
	}

	for _, item := range allEntries {
		if item.ID == id {
			ctbE = item
		}
	}

	if ctbE.ID == uuid.Nil {
		return ctbE, fmt.Errorf("did not find ctbE with ID of %s", id)
	}

	return ctbE, nil
}

func (cM *CrontabManager) RemoveCrontabEntryByID(id uuid.UUID) error {
	allEntries, err := cM.GetAllCrontabEntries()
	if err != nil {
		return err
	}

	var entriesToKeep []CrontabEntry

	for _, item := range allEntries {
		if item.ID == id {
			continue
		}

		entriesToKeep = append(entriesToKeep, item)
	}

	if err = cM.emptyCrontab(); err != nil {
		return err
	}

	err = cM.withCrontab(func(f *os.File) error {
		for _, item := range entriesToKeep {
			_, err = fmt.Fprintf(f, cronFormat, item.Cron, item.Cmd, item.ID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (cM *CrontabManager) withCrontab(fn func(f *os.File) error) error {
	cM.mu.Lock()
	defer cM.mu.Unlock()
	file := config.Config.CrontabFile
	if file == "" {
		return errCrontabFileNotSet
	}

	crontab, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer crontab.Close()

	err = fn(crontab)

	if err != nil {
		return err
	}

	return nil
}

func (cM *CrontabManager) emptyCrontab() error {
	cM.mu.Lock()
	defer cM.mu.Unlock()
	file := config.Config.CrontabFile
	if file == "" {
		return errCrontabFileNotSet
	}

	err := os.Truncate(file, 0)
	if err != nil {
		return err
	}

	return nil
}
