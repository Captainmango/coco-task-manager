package crontab

import (
	"errors"
	"fmt"
	"os"

	"github.com/captainmango/coco-cron-parser/internal/config"
	"github.com/captainmango/coco-cron-parser/internal/data"
)

const cronFormat = "%s root %s | tee /tmp/log\n"

var (
	errCrontabFileNotSet = errors.New("crontab file not set")
)

// Sets the cron printing mode to RAW_EXPRESSION and writes to the configured crontab file
func WriteCronToSchedule(cron data.Cron, cmd string) error {
	cron.PrintingMode = data.RAW_EXPRESSION
	file := config.Config.CrontabFile

	if file == "" {
		return errCrontabFileNotSet
	}

	crontab, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(crontab, cronFormat, cron, cmd)

	if err != nil {
		return err
	}

	return nil
}