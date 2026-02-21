package coco_http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli/v3"

	coco_cli_mock "github.com/captainmango/coco-cron-parser/internal/cli/mocks"
	"github.com/captainmango/coco-cron-parser/internal/crontab"
	"github.com/captainmango/coco-cron-parser/internal/parser"
	"github.com/captainmango/coco-cron-parser/internal/resources"
	"github.com/captainmango/coco-cron-parser/internal/resources/mocks"
)

type mockAppWithResources struct {
	*app
	mockCrontab         *mocks.MockCrontabHandler
	mockQueue           *mocks.MockQueueHandler
	mockCommandRegistry *coco_cli_mock.MockCommandRegistry
}

func Test_handleLivez(t *testing.T) {
	t.Run("returns OK status with host and protocol", func(t *testing.T) {
		mockApp := getMockApp(t)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/livez", nil)
		req.Host = "localhost:3000"
		w := httptest.NewRecorder()

		mockApp.handleLivez(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		var out map[string]any
		err := json.NewDecoder(res.Body).Decode(&out)
		assert.NoError(t, err)

		assert.Equal(t, "OK", out["status"])
		assert.Equal(t, "localhost:3000", out["host"])
		assert.Equal(t, "HTTP/1.1", out["protocol"])
	})
}

func Test_handleGetScheduledTasks(t *testing.T) {
	t.Run("returns scheduled tasks successfully", func(t *testing.T) {
		mockApp := getMockApp(t)

		cronExpr, _ := parser.NewParser(parser.WithInput("* * * * *", true))
		parsedCron, _ := cronExpr.Parse()

		expectedEntries := []crontab.CrontabEntry{
			{
				ID:   uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				Cron: parsedCron,
				Cmd:  "cli start-game room1",
			},
			{
				ID:   uuid.MustParse("660e8400-e29b-41d4-a716-446655440001"),
				Cron: parsedCron,
				Cmd:  "cli start-game room2",
			},
		}

		mockApp.mockCrontab.On("GetAllCrontabEntries").Return(expectedEntries, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/scheduled", nil)
		w := httptest.NewRecorder()

		mockApp.handleGetScheduledTasks(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		var out Response[[]ScheduledTaskResponse]
		err := json.NewDecoder(res.Body).Decode(&out)
		assert.NoError(t, err)

		assert.Equal(t, SCHEDULED_TASK, out.Type)
		assert.Len(t, out.Data, 2)
		assert.Equal(t, "cli start-game room1", out.Data[0].Command)
		assert.Equal(t, "cli start-game room2", out.Data[1].Command)
		assert.Empty(t, out.Error)
		mockApp.mockCrontab.AssertExpectations(t)
	})

	t.Run("returns error when GetAllCrontabEntries fails", func(t *testing.T) {
		mockApp := getMockApp(t)

		mockApp.mockCrontab.On("GetAllCrontabEntries").Return([]crontab.CrontabEntry{}, errors.New("database connection failed"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/scheduled", nil)
		w := httptest.NewRecorder()

		mockApp.handleGetScheduledTasks(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		var out Response[[]ScheduledTaskResponse]
		err := json.NewDecoder(res.Body).Decode(&out)
		assert.NoError(t, err)

		assert.Equal(t, "error", out.Type)
		assert.Contains(t, out.Error, "database connection failed")
		mockApp.mockCrontab.AssertExpectations(t)
	})

	t.Run("returns empty array when no scheduled tasks", func(t *testing.T) {
		mockApp := getMockApp(t)

		mockApp.mockCrontab.On("GetAllCrontabEntries").Return([]crontab.CrontabEntry{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/scheduled", nil)
		w := httptest.NewRecorder()

		mockApp.handleGetScheduledTasks(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		var out Response[[]ScheduledTaskResponse]
		err := json.NewDecoder(res.Body).Decode(&out)
		assert.NoError(t, err)

		assert.Equal(t, SCHEDULED_TASK, out.Type)
		assert.Len(t, out.Data, 0)
		mockApp.mockCrontab.AssertExpectations(t)
	})
}

func Test_handleGetTasks(t *testing.T) {
	t.Run("returns all registered tasks", func(t *testing.T) {
		mockApp := getMockApp(t)

		mockApp.mockCommandRegistry.On("All").Return([]*cli.Command{
			{
				Name: "start-game",
			},
			{
				Name: "pull-messages",
			},
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
		w := httptest.NewRecorder()

		mockApp.handleGetTasks(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		var out Response[[]TaskResponse]
		err := json.NewDecoder(res.Body).Decode(&out)
		assert.NoError(t, err)

		assert.Equal(t, TASK, out.Type)
		assert.Len(t, out.Data, 2)
		assert.Equal(t, "start-game", out.Data[0].Slug)
		assert.Equal(t, "pull-messages", out.Data[1].Slug)
		assert.Empty(t, out.Error)

		mockApp.mockCommandRegistry.AssertExpectations(t)
	})

	t.Run("returns empty array when no tasks registered", func(t *testing.T) {
		mockApp := getMockApp(t)

		mockApp.mockCommandRegistry.On("All").Return([]*cli.Command{})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
		w := httptest.NewRecorder()

		mockApp.handleGetTasks(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		var out Response[[]TaskResponse]
		err := json.NewDecoder(res.Body).Decode(&out)
		assert.NoError(t, err)

		assert.Equal(t, TASK, out.Type)
		assert.Len(t, out.Data, 0)
		assert.Empty(t, out.Error)

		mockApp.mockCommandRegistry.AssertExpectations(t)
	})
}

func Test_handleScheduleTask(t *testing.T) {
	t.Run("schedules task successfully", func(t *testing.T) {
		mockApp := getMockApp(t)

		mockApp.mockCommandRegistry.On("Find").Return(&cli.Command{
			Name: "start-game",
		}, nil)
		mockApp.mockCrontab.On("WriteCrontabEntries", mock.Anything).Return(nil)

		input := ScheduleTaskRequest{
			TaskId:        "start-game",
			ScheduledTime: "*/5 * * * *",
			Args: struct {
				RoomId string `json:"room_id"`
			}{
				RoomId: "room123",
			},
		}

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockApp.handleScheduleTask(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusAccepted, res.StatusCode)

		var out Response[ScheduledTaskResponse]
		err := json.NewDecoder(res.Body).Decode(&out)
		assert.NoError(t, err)

		assert.Equal(t, SCHEDULED_TASK, out.Type)
		assert.NotEqual(t, uuid.Nil, out.Data.ID)
		assert.Equal(t, "*/5 * * * *", out.Data.Cron)
		assert.Equal(t, "cli start-game room123", out.Data.Command)

		mockApp.mockCommandRegistry.AssertExpectations(t)
		mockApp.mockCrontab.AssertExpectations(t)
		mockApp.mockQueue.AssertExpectations(t)
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		mockApp := getMockApp(t)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockApp.handleScheduleTask(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		var out Response[ScheduledTaskResponse]
		err := json.NewDecoder(res.Body).Decode(&out)
		assert.NoError(t, err)

		assert.Equal(t, "error", out.Type)
		assert.Contains(t, out.Error, "bad JSON")
	})

	t.Run("returns error when task_id not found", func(t *testing.T) {
		mockApp := getMockApp(t)

		mockApp.mockCommandRegistry.On("Find").Return((*cli.Command)(nil), errors.New("command not found"))

		input := ScheduleTaskRequest{
			TaskId:        "non-existent-task",
			ScheduledTime: "*/5 * * * *",
			Args: struct {
				RoomId string `json:"room_id"`
			}{
				RoomId: "room123",
			},
		}

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockApp.handleScheduleTask(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

		var out Response[ScheduledTaskResponse]
		err := json.NewDecoder(res.Body).Decode(&out)
		assert.NoError(t, err)

		assert.Equal(t, "error", out.Type)
		assert.Contains(t, out.Error, "command not found")

		mockApp.mockCommandRegistry.AssertExpectations(t)
	})

	t.Run("returns error for invalid cron expression", func(t *testing.T) {
		mockApp := getMockApp(t)

		mockApp.mockCommandRegistry.On("Find").Return(&cli.Command{
			Name: "start-game",
		}, nil)

		input := ScheduleTaskRequest{
			TaskId:        "start-game",
			ScheduledTime: "invalid-cron",
			Args: struct {
				RoomId string `json:"room_id"`
			}{
				RoomId: "room123",
			},
		}

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockApp.handleScheduleTask(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

		var out Response[ScheduledTaskResponse]
		err := json.NewDecoder(res.Body).Decode(&out)
		assert.NoError(t, err)

		assert.Equal(t, "error", out.Type)
		assert.NotEmpty(t, out.Error)

		mockApp.mockCommandRegistry.AssertExpectations(t)
	})

	t.Run("returns error when WriteCrontabEntries fails", func(t *testing.T) {
		mockApp := getMockApp(t)

		mockApp.mockCommandRegistry.On("Find").Return(&cli.Command{
			Name: "start-game",
		}, nil)
		mockApp.mockCrontab.On("WriteCrontabEntries", mock.Anything).Return(errors.New("permission denied"))

		input := ScheduleTaskRequest{
			TaskId:        "start-game",
			ScheduledTime: "*/5 * * * *",
			Args: struct {
				RoomId string `json:"room_id"`
			}{
				RoomId: "room123",
			},
		}

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockApp.handleScheduleTask(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

		var out Response[ScheduledTaskResponse]
		err := json.NewDecoder(res.Body).Decode(&out)
		assert.NoError(t, err)

		assert.Equal(t, "error", out.Type)
		assert.Contains(t, out.Error, "permission denied")

		mockApp.mockCrontab.AssertExpectations(t)
	})

	t.Run("returns error for unknown fields in JSON", func(t *testing.T) {
		mockApp := getMockApp(t)

		jsonBody := `{"task_id": "start-game", "scheduled_time": "*/5 * * * *", "args": {"room_id": "room123"}, "unknown_field": "value"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", strings.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockApp.handleScheduleTask(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		var out Response[ScheduledTaskResponse]
		err := json.NewDecoder(res.Body).Decode(&out)
		assert.NoError(t, err)

		assert.Equal(t, "error", out.Type)
		assert.Contains(t, out.Error, "unknown key")
	})

	t.Run("returns error for empty request body", func(t *testing.T) {
		mockApp := getMockApp(t)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", strings.NewReader(""))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockApp.handleScheduleTask(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		var out Response[ScheduledTaskResponse]
		err := json.NewDecoder(res.Body).Decode(&out)
		assert.NoError(t, err)

		assert.Equal(t, "error", out.Type)
		assert.Contains(t, out.Error, "body must not be empty")
	})
}

func Test_handleRemoveTask(t *testing.T) {
	t.Run("removes task successfully", func(t *testing.T) {
		mockApp := getMockApp(t)

		taskUUID := "550e8400-e29b-41d4-a716-446655440000"
		parsedUUID := uuid.MustParse(taskUUID)

		mockApp.mockCrontab.On("RemoveCrontabEntryByID", parsedUUID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+taskUUID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", taskUUID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		mockApp.handleRemoveTask(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		mockApp.mockCrontab.AssertExpectations(t)
	})

	t.Run("returns error for invalid UUID format", func(t *testing.T) {
		mockApp := getMockApp(t)

		invalidUUID := "not-a-valid-uuid"
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+invalidUUID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", invalidUUID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		mockApp.handleRemoveTask(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		var out Response[string]
		err := json.NewDecoder(res.Body).Decode(&out)
		assert.NoError(t, err)

		assert.Equal(t, "error", out.Type)
		assert.Contains(t, out.Error, "invalid UUID")
	})
}

func getMockApp(t *testing.T) *mockAppWithResources {
	mockCrontab := &mocks.MockCrontabHandler{}
	mockQueue := &mocks.MockQueueHandler{}
	mockCommandRegistry := &coco_cli_mock.MockCommandRegistry{}

	appInstance := &app{
		logger: newTestLogger(t),
		resources: resources.Resources{
			TaskResource: resources.CreateTaskResource(
				mockCrontab,
				mockQueue,
			),
		},
		commandsRegistry: mockCommandRegistry,
	}

	return &mockAppWithResources{
		app:                 appInstance,
		mockCrontab:         mockCrontab,
		mockQueue:           mockQueue,
		mockCommandRegistry: mockCommandRegistry,
	}
}

func newTestLogger(t *testing.T) *slog.Logger {
	t.Helper()
	handler := slog.NewTextHandler(&testWriter{t}, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(handler)
}

type testWriter struct {
	t *testing.T
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	w.t.Log(string(p))
	return len(p), nil
}
