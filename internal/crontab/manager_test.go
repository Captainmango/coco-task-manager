package crontab

import (
	"fmt"
	"os"
	"testing"

	"github.com/captainmango/coco-cron-parser/internal/config"
	"github.com/captainmango/coco-cron-parser/internal/data"
	"github.com/captainmango/coco-cron-parser/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const expectedCrontabFormat = "%s root %s | tee /tmp/log\n"

type CronTabManagerTestSuite struct {
	suite.Suite
	cron data.Cron
}

func Test_RunManagerTestSuite(t *testing.T) {
	suite.Run(t, new(CronTabManagerTestSuite))
}

func (s *CronTabManagerTestSuite) SetupTest() {
	config.BootstrapConfig()
	config.Config.CrontabFile = utils.BasePath("e2e/storage/crontab")
	s.cron = exampleTestCron()
}

func (s *CronTabManagerTestSuite) TearDownTest() {
	resetFileFromPath(s.T(), config.Config.CrontabFile)
}


func (s *CronTabManagerTestSuite) Test_ItWritesToCrontabFile() {
	err := WriteCronToSchedule(s.cron, "./test-command")
	if err != nil {
		s.T().Fatal(err)
		return
	}

	expected := fmt.Sprintf(expectedCrontabFormat, s.cron, "./test-command")
	
	out := readFromPath(s.T(), config.Config.CrontabFile)

	assert.Equal(s.T(), expected, out)
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