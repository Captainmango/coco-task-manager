package crontab

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/captainmango/coco-cron-parser/internal/config"
	"github.com/captainmango/coco-cron-parser/internal/parser"
	"github.com/captainmango/coco-cron-parser/internal/utils"
)

const expectedCrontabFormat = cronFormat

type CronTabManagerTestSuite struct {
	suite.Suite
	cron parser.Cron
	cM   CrontabHandler
}

func Test_RunManagerTestSuite(t *testing.T) {
	suite.Run(t, new(CronTabManagerTestSuite))
}

func (s *CronTabManagerTestSuite) SetupTest() {
	config.BootstrapConfig()
	config.Config.CrontabFile = utils.BasePath("e2e/storage/crontab")
	s.cron = exampleTestCron()
	s.cM = &CrontabManager{}
}

func (s *CronTabManagerTestSuite) TearDownTest() {
	resetFileFromPath(s.T(), config.Config.CrontabFile)
}

func (s *CronTabManagerTestSuite) Test_ItWritesToCrontabFile() {
	fakeUuID, _ := uuid.NewUUID()
	err := s.cM.WriteCrontabEntries([]CrontabEntry{
		{
			ID:   fakeUuID,
			Cron: s.cron,
			Cmd:  "./test-command",
		},
	})

	assert.NoError(s.T(), err)

	expected := fmt.Sprintf(expectedCrontabFormat, s.cron, "./test-command", fakeUuID.String())

	out := readFromPath(s.T(), config.Config.CrontabFile)

	assert.Equal(s.T(), expected, out)
}

func (s *CronTabManagerTestSuite) Test_ItGetsCrontabsInFile() {
	fakeUuIDOne, _ := uuid.NewUUID()
	fakeUuIDTwo, _ := uuid.NewUUID()

	err := s.cM.WriteCrontabEntries(fixtureCrontabs(fakeUuIDOne, fakeUuIDTwo))

	assert.NoError(s.T(), err)

	entries, err := s.cM.GetAllCrontabEntries()

	assert.NoError(s.T(), err)
	assert.Len(s.T(), entries, 2)
}

func (s *CronTabManagerTestSuite) Test_ItGetsCrontabByID() {
	fakeUuIDOne, _ := uuid.NewUUID()
	fakeUuIDTwo, _ := uuid.NewUUID()

	err := s.cM.WriteCrontabEntries(fixtureCrontabs(fakeUuIDOne, fakeUuIDTwo))

	assert.NoError(s.T(), err)

	ctbE, err := s.cM.GetCrontabEntryByID(fakeUuIDTwo)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), ctbE)
	assert.Equal(s.T(), fakeUuIDTwo, ctbE.ID)
}

func (s *CronTabManagerTestSuite) Test_ItDeletesCrontabByID() {
	fakeUuIDOne, _ := uuid.NewUUID()
	fakeUuIDTwo, _ := uuid.NewUUID()

	err := s.cM.WriteCrontabEntries(fixtureCrontabs(fakeUuIDOne, fakeUuIDTwo))

	assert.NoError(s.T(), err)

	err = s.cM.RemoveCrontabEntryByID(fakeUuIDTwo)

	assert.NoError(s.T(), err)
	entries, _ := s.cM.GetAllCrontabEntries()
	assert.Len(s.T(), entries, 1)
}

func exampleTestCron() parser.Cron {
	return parser.Cron{
		Data: []parser.CronFragment{
			{
				Expr:         "*",
				Kind:         parser.WILDCARD,
				FragmentType: parser.MINUTE,
			},
			{
				Expr:         "*",
				Kind:         parser.WILDCARD,
				FragmentType: parser.MINUTE,
			},
			{
				Expr:         "*",
				Kind:         parser.WILDCARD,
				FragmentType: parser.MINUTE,
			},
			{
				Expr:         "*",
				Kind:         parser.WILDCARD,
				FragmentType: parser.MINUTE,
			},
			{
				Expr:         "*",
				Kind:         parser.WILDCARD,
				FragmentType: parser.MINUTE,
			},
		},
	}
}

func fixtureCrontabs(uuids ...uuid.UUID) []CrontabEntry {
	var out []CrontabEntry

	for _, item := range uuids {
		out = append(out, CrontabEntry{
			ID:   item,
			Cron: exampleTestCron(),
			Cmd:  "./test-command",
		})
	}

	return out
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
