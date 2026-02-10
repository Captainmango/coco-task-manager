package resources

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/captainmango/coco-cron-parser/internal/crontab"
	"github.com/captainmango/coco-cron-parser/internal/parser"
)

func Test_ItCanReadCronTabs(t *testing.T) {
	id, _ := uuid.NewV7()

	mockCrontabHandler := new(mockCrontabHandler)
	mockQueueHandler := new(mockQueueHandler)

	mockCrontabHandler.On("GetAllCrontabEntries").
		Return([]crontab.CrontabEntry{
			{
				ID:   id,
				Cron: exampleTestCron(),
				Cmd:  "test-command",
			},
		}, nil)

	tR := CreateTaskResource(mockCrontabHandler, mockQueueHandler)

	d, _ := tR.crontabManager.GetAllCrontabEntries()

	mockCrontabHandler.AssertExpectations(t)
	assert.NotEmpty(t, d)
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
