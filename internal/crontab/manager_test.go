package crontab

import (
	"fmt"
	"os"
	"testing"

	"github.com/captainmango/coco-cron-parser/internal/config"
	"github.com/captainmango/coco-cron-parser/internal/data"
	"github.com/captainmango/coco-cron-parser/internal/utils"
	"github.com/stretchr/testify/assert"
)

const expectedCrontabFormat = "%s root %s | tee /tmp/log\n"

func Test_ItWritesToCrontabFile(t *testing.T) {
	config.BootstrapConfig()
	config.Config.CrontabFile = utils.BasePath("e2e/storage/crontab")

	c := exampleTestCron()

	err := WriteCronToSchedule(c, "./test-command")
	if err != nil {
		t.Fatal(err)
		return
	}

	expected := fmt.Sprintf(expectedCrontabFormat, c, "./test-command")
	
	out := readFromPath(t, config.Config.CrontabFile)

	assert.Equal(t, expected, out)

	resetFileFromPath(t, config.Config.CrontabFile)
}

func exampleTestCron() data.Cron {
	return  data.Cron{
		Data: []data.CronFragment{
			{
				Expr: "*",
				Kind: data.WILDCARD,
				FragmentType: data.MINUTE,
			},
						{
				Expr: "*",
				Kind: data.WILDCARD,
				FragmentType: data.MINUTE,
			},
						{
				Expr: "*",
				Kind: data.WILDCARD,
				FragmentType: data.MINUTE,
			},
						{
				Expr: "*",
				Kind: data.WILDCARD,
				FragmentType: data.MINUTE,
			},
						{
				Expr: "*",
				Kind: data.WILDCARD,
				FragmentType: data.MINUTE,
			},
		},
	}
}

func readFromPath(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return string(data)
}

func resetFileFromPath(t *testing.T, path string) {
	t.Helper()

	err := os.Truncate(path, 0)
    if err != nil {
        t.Fatalf("truncate failed: %v", err)
		return
    }
}