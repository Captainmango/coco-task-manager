package crontab

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/captainmango/coco-cron-parser/internal/config"
	"github.com/captainmango/coco-cron-parser/internal/data"
	"github.com/google/uuid"
)

const cronFormat = "%s root %s | tee /tmp/log # %s\n"

var (
	errCrontabFileNotSet = errors.New("crontab file not set")
)

// Sets the cron printing mode to RAW_EXPRESSION and writes to the configured crontab file
func WriteCronToSchedule(cron data.Cron, cmd string, id string) error {
	cron.PrintingMode = data.RAW_EXPRESSION

	err := withCrontab(func(f *os.File) error {
		_, err := fmt.Fprintf(f, cronFormat, cron, cmd, id)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func GetAllCrontabEntries() ([]CrontabEntry, error) {
	var out []CrontabEntry

	err := withCrontab(func(f *os.File) error {
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

func GetCrontabEntryByID(id uuid.UUID) (CrontabEntry, error) {
	var ctbE CrontabEntry
	allEntries, err := GetAllCrontabEntries()
	if err != nil {
		return ctbE, nil
	}

	for _, item := range allEntries {
		if item.ID == id {
			ctbE = item
		}
	}

	if ctbE.ID == uuid.Nil {
		return ctbE, fmt.Errorf("did not fine ctbE with ID of %s", id)
	}

	return ctbE, nil
}

func RemoveCrontabEntryByID(id uuid.UUID) error {
	allEntries, err := GetAllCrontabEntries()
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

	if err = emptyCrontab(); err != nil {
		return err
	}

	err = withCrontab(func(f *os.File) error {
		for _, item := range entriesToKeep {
			_, err := fmt.Fprintf(f, cronFormat, item.Cron, item.Cmd, item.ID)
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

func withCrontab(fn func(f *os.File)error) error {
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

func emptyCrontab() error {
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