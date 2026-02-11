package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/captainmango/coco-cron-parser/internal/crontab"
	"github.com/captainmango/coco-cron-parser/internal/msq"
)

// Mock of CrontabHandler interface. Used only in tests
type MockCrontabHandler struct {
	mock.Mock
}

func (mch *MockCrontabHandler) GetAllCrontabEntries() ([]crontab.CrontabEntry, error) {
	args := mch.Called()
	return args.Get(0).([]crontab.CrontabEntry), args.Error(1)
}

func (mch *MockCrontabHandler) GetCrontabEntryByID(id uuid.UUID) (crontab.CrontabEntry, error) {
	args := mch.Called(id)
	return args.Get(0).(crontab.CrontabEntry), args.Error(1)
}

func (mch *MockCrontabHandler) WriteCrontabEntries(entries []crontab.CrontabEntry) error {
	args := mch.Called(entries)
	return args.Error(0)
}

func (mch *MockCrontabHandler) RemoveCrontabEntryByID(id uuid.UUID) error {
	args := mch.Called(id)
	return args.Error(0)
}

// Mock of AdvancedMessageQueueHandler interface. Used only in tests
type MockQueueHandler struct {
	mock.Mock
}

func (mqh *MockQueueHandler) PushMessage(routingKey, payload string) error {
	args := mqh.Called(routingKey, payload)
	return args.Error(0)
}

func (mqh *MockQueueHandler) ConsumeMessages(
	ctx context.Context,
	routingKey string,
	fn msq.ConsumeMessageFn,
) error {

	args := mqh.Called(ctx, routingKey, fn)
	return args.Error(0)
}
