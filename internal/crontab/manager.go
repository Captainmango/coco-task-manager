package crontab

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/captainmango/coco-cron-parser/internal/config"
	"github.com/captainmango/coco-cron-parser/internal/data"
	"github.com/captainmango/coco-cron-parser/internal/parser"
	"github.com/google/uuid"
)

const cronFormat = "%s root %s | tee /tmp/log # %s\n"

var (
	errCrontabFileNotSet = errors.New("crontab file not set")
)

// Sets the cron printing mode to RAW_EXPRESSION and writes to the configured crontab file
func WriteCronToSchedule(cron data.Cron, cmd string, id string) error {
	cron.PrintingMode = data.RAW_EXPRESSION
	file := config.Config.CrontabFile

	if file == "" {
		return errCrontabFileNotSet
	}

	crontab, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer crontab.Close()

	_, err = fmt.Fprintf(crontab, cronFormat, cron, cmd, id)

	if err != nil {
		return err
	}

	return nil
}

func GetAllCrontabEntries() ([]parser.CrontabEntry, error) {
	var out []parser.CrontabEntry
	file := config.Config.CrontabFile

	if file == "" {
		return nil, errCrontabFileNotSet
	}

	crontab, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer crontab.Close()

	scanner := bufio.NewScanner(crontab)

	for scanner.Scan() {
		line := scanner.Text()

		ctbE, err := parser.NewCrontabEntryFromString(line)
		if err != nil {
			return nil, err
		}

		out = append(out, ctbE)
	}

	return out, nil
}

func GetCrontabEntryByID(id uuid.UUID) (parser.CrontabEntry, error) {
	var ctbE parser.CrontabEntry
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

	var entriesToKeep []parser.CrontabEntry

	for _, item := range allEntries {
		if item.ID == id {
			continue
		}

		entriesToKeep = append(entriesToKeep, item)
	}

	err = emptyCrontab()
	if err != nil {
		return err
	}

	for _, item := range entriesToKeep {
		err = WriteCronToSchedule(item.Cron, item.Cmd, item.ID.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func emptyCrontab() error {
	file := config.Config.CrontabFile
	crontab, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	err = os.Truncate(file, 0)
	if err != nil {
		return err
	}
	crontab.Close()
	return nil
}