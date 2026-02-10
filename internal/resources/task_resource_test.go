package resources

import (
	"testing"

	"github.com/captainmango/coco-cron-parser/internal/config"
	"github.com/captainmango/coco-cron-parser/internal/parser"
	"github.com/captainmango/coco-cron-parser/internal/utils"
	"github.com/stretchr/testify/suite"
)

type TaskResourceTestSuite struct {
	suite.Suite
	tR          TaskResource
	cron        parser.Cron
	mockResetFn func()
}

func Test_RunTaskResourceTestSuite(t *testing.T) {
	suite.Run(t, new(TaskResourceTestSuite))
}

func (s *TaskResourceTestSuite) SetupTest() {
	config.BootstrapConfig()
	config.Config.CrontabFile = utils.BasePath("e2e/storage/crontab")
}

