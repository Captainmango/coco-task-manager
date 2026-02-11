package resources

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/captainmango/coco-cron-parser/internal/crontab"
	"github.com/captainmango/coco-cron-parser/internal/msq"
	"github.com/captainmango/coco-cron-parser/internal/parser"
	"github.com/captainmango/coco-cron-parser/internal/resources/mocks"
)

func Test_ItCanReadCronTabs(t *testing.T) {
	t.Parallel()
	id, _ := uuid.NewV7()

	mockCrontabHandler := new(mocks.MockCrontabHandler)
	mockQueueHandler := new(mocks.MockQueueHandler)

	mockCrontabHandler.On("GetAllCrontabEntries").
		Return([]crontab.CrontabEntry{
			{
				ID:   id,
				Cron: exampleTestCron(),
				Cmd:  "test-command",
			},
		}, nil)

	tR := CreateTaskResource(mockCrontabHandler, mockQueueHandler)

	d, _ := tR.GetAllCrontabEntries()

	mockCrontabHandler.AssertExpectations(t)
	assert.NotEmpty(t, d)
}

func Test_ItReturnsNoCrontabsWhenFileEmpty(t *testing.T) {
	t.Parallel()

	mockCrontabHandler := new(mocks.MockCrontabHandler)
	mockQueueHandler := new(mocks.MockQueueHandler)

	mockCrontabHandler.On("GetAllCrontabEntries").
		Return([]crontab.CrontabEntry{}, nil)

	tR := CreateTaskResource(mockCrontabHandler, mockQueueHandler)

	d, _ := tR.GetAllCrontabEntries()

	mockCrontabHandler.AssertExpectations(t)
	assert.Empty(t, d)
}

func Test_ItReturnsErrorsWhenFileCannotBeRead(t *testing.T) {
	t.Parallel()

	mockCrontabHandler := new(mocks.MockCrontabHandler)
	mockQueueHandler := new(mocks.MockQueueHandler)

	mockCrontabHandler.On("GetAllCrontabEntries").
		Return([]crontab.CrontabEntry{}, errors.New("file read error"))

	tR := CreateTaskResource(mockCrontabHandler, mockQueueHandler)

	_, err := tR.GetAllCrontabEntries()

	mockCrontabHandler.AssertExpectations(t)
	assert.Error(t, err, "file read error")
}

func Test_ItCanGetTaskById(t *testing.T) {
	t.Parallel()
	id, _ := uuid.NewV7()

	mockCrontabHandler := new(mocks.MockCrontabHandler)
	mockQueueHandler := new(mocks.MockQueueHandler)

	ctb := crontab.CrontabEntry{
		ID:   id,
		Cron: exampleTestCron(),
		Cmd:  "test-command",
	}

	mockCrontabHandler.On("GetCrontabEntryByID", id).
		Return(ctb, nil)

	tR := CreateTaskResource(mockCrontabHandler, mockQueueHandler)

	d, _ := tR.GetTaskByID(id)

	mockCrontabHandler.AssertExpectations(t)
	assert.Equal(t, ctb.ID, d.ID)
	assert.Equal(t, ctb.Cron, d.Cron)
	assert.Equal(t, ctb.Cmd, d.Cmd)
}

func Test_ItErrorsIfEntryDoesNotExistByID(t *testing.T) {
	t.Parallel()
	id, _ := uuid.NewV7()

	emptyResult := crontab.CrontabEntry{}

	mockCrontabHandler := new(mocks.MockCrontabHandler)
	mockQueueHandler := new(mocks.MockQueueHandler)
	mockCrontabHandler.On("GetCrontabEntryByID", id).
		Return(emptyResult, errors.New("entry does not exist"))

	tR := CreateTaskResource(mockCrontabHandler, mockQueueHandler)

	d, err := tR.GetTaskByID(id)

	mockCrontabHandler.AssertExpectations(t)
	assert.Error(t, err, "entry does not exist")
	assert.Equal(t, emptyResult.ID, d.ID)
	assert.Equal(t, emptyResult.Cron, d.Cron)
	assert.Equal(t, emptyResult.Cmd, d.Cmd)
}

func Test_ItCanWriteTasksToFile(t *testing.T) {
	t.Parallel()

	mockCrontabHandler := new(mocks.MockCrontabHandler)
	mockQueueHandler := new(mocks.MockQueueHandler)

	var capturedID uuid.UUID
	mockCrontabHandler.On("WriteCrontabEntries", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			ctb := args.Get(0).([]crontab.CrontabEntry)
			capturedID = ctb[0].ID
		})

	tR := CreateTaskResource(mockCrontabHandler, mockQueueHandler)

	taskId, err := tR.ScheduleTask("* * * * *", "test-command")

	mockCrontabHandler.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, capturedID, taskId)
}

func Test_ItHandlesErrorsWhenWriting(t *testing.T) {
	t.Parallel()

	mockCrontabHandler := new(mocks.MockCrontabHandler)
	mockQueueHandler := new(mocks.MockQueueHandler)
	mockCrontabHandler.On("WriteCrontabEntries", mock.Anything).
		Return(errors.New("error writing"))

	tR := CreateTaskResource(mockCrontabHandler, mockQueueHandler)

	_, err := tR.ScheduleTask("* * * * *", "test-command")

	mockCrontabHandler.AssertExpectations(t)
	assert.NotNil(t, err)
	assert.Error(t, err, "error writing")
}

func Test_ItCanRemoveCrontabFromFile(t *testing.T) {
	t.Parallel()

	id, _ := uuid.NewV7()

	mockCrontabHandler := new(mocks.MockCrontabHandler)
	mockQueueHandler := new(mocks.MockQueueHandler)
	mockCrontabHandler.On("RemoveCrontabEntryByID", id).
		Return(nil)

	tR := CreateTaskResource(mockCrontabHandler, mockQueueHandler)

	err := tR.RemoveTaskByID(id)

	mockCrontabHandler.AssertExpectations(t)
	assert.Nil(t, err)
}

func Test_ItHandlesErrorWhenRemovingFromFile(t *testing.T) {
	t.Parallel()

	id, _ := uuid.NewV7()

	mockCrontabHandler := new(mocks.MockCrontabHandler)
	mockQueueHandler := new(mocks.MockQueueHandler)
	mockCrontabHandler.On("RemoveCrontabEntryByID", id).
		Return(errors.New("error removing"))

	tR := CreateTaskResource(mockCrontabHandler, mockQueueHandler)

	err := tR.RemoveTaskByID(id)

	mockCrontabHandler.AssertExpectations(t)
	assert.NotNil(t, err)
	assert.Error(t, err, "error removing")
}

func Test_ItCanPushMessages(t *testing.T) {
	t.Parallel()

	mockCrontabHandler := new(mocks.MockCrontabHandler)
	mockQueueHandler := new(mocks.MockQueueHandler)

	routingKey := "coco_tasks.start_game"
	data := "test-payload"
	payload := msq.StartGamePayload{
		RoomId: data,
	}

	msgPayload, _ := json.Marshal(payload)

	mockQueueHandler.On("PushMessage", routingKey, string(msgPayload)).
		Return(nil)

	tR := CreateTaskResource(mockCrontabHandler, mockQueueHandler)

	err := tR.PushStartGameMessage(payload)

	mockCrontabHandler.AssertExpectations(t)
	assert.Nil(t, err)
}

func Test_WhenPushMessagesFails(t *testing.T) {
	t.Parallel()

	mockCrontabHandler := new(mocks.MockCrontabHandler)
	mockQueueHandler := new(mocks.MockQueueHandler)

	routingKey := "coco_tasks.start_game"
	data := "test-payload"
	payload := msq.StartGamePayload{
		RoomId: data,
	}

	msgPayload, _ := json.Marshal(payload)

	mockQueueHandler.On("PushMessage", routingKey, string(msgPayload)).
		Return(errors.New("push failed"))

	tR := CreateTaskResource(mockCrontabHandler, mockQueueHandler)

	err := tR.PushStartGameMessage(payload)

	mockCrontabHandler.AssertExpectations(t)
	assert.NotNil(t, err)
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
