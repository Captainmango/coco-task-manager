package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli/v3"
)

// Mock of CrontabHandler interface. Used only in tests
type MockCommandRegistry struct {
	mock.Mock
}

func (mcr *MockCommandRegistry) Find(key string) (*cli.Command, error) {
	args := mcr.Called()
	return args.Get(0).(*cli.Command), args.Error(1)
}

func (mcr *MockCommandRegistry) All() []*cli.Command {
	args := mcr.Called()
	return args.Get(0).([]*cli.Command)
}
