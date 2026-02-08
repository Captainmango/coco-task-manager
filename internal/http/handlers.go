package coco_http

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a *app) handleLivez(w http.ResponseWriter, r *http.Request) {
	a.writeJSON(w, http.StatusOK, map[string]any{
		"status":   "OK",
		"host":     r.Host,
		"protocol": r.Proto,
	}, nil)
}

func (a *app) handleGetScheduledTasks(w http.ResponseWriter, r *http.Request) {
	entries, err := a.resources.TaskResource.GetAllCrontabEntries()

	if err != nil {
		errRes := NewResponse(
			WithError(
				err,
				[]ScheduledTaskResponse{},
			),
		)

		a.writeJSON(w, http.StatusBadRequest, errRes, nil)
		return
	}

	var out []ScheduledTaskResponse = []ScheduledTaskResponse{}

	for _, item := range entries {
		i := ScheduledTaskResponse{
			ID:      item.ID,
			Command: item.Cmd,
			Cron:    item.Cron.String(),
		}

		out = append(out, i)
	}

	res := NewResponse(
		WithData(
			SCHEDULED_TASK,
			out,
		),
	)

	a.writeJSON(w, http.StatusOK, res, nil)
}

func (a *app) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	cmds := a.commandsRegistry.All()

	var out []TaskResponse
	for _, cmd := range cmds {
		tOut := TaskResponse{
			Slug: cmd.Name,
			Args: slices.Concat(cmd.Args().Slice()),
		}

		out = append(out, tOut)
	}

	res := NewResponse(WithData(TASK, out))
	a.writeJSON(w, http.StatusOK, res, nil)
}

func (a *app) handleScheduleTask(w http.ResponseWriter, r *http.Request) {
	var input ScheduleTaskRequest

	err := a.readJSON(w, r, &input)
	if err != nil {
		res := NewResponse(WithError(err, ScheduledTaskResponse{}))
		a.writeJSON(w, http.StatusInternalServerError, res, nil)
		return
	}

	cmd, err := a.commandsRegistry.Find(input.TaskId)
	if err != nil {
		res := NewResponse(WithError(err, ScheduledTaskResponse{}))
		a.writeJSON(w, http.StatusUnprocessableEntity, res, nil)
		return
	}

	var cmdStringBuilder strings.Builder
	cmdStringBuilder.WriteString(fmt.Sprintf("cli %s ", cmd.Name))
	cmdStringBuilder.WriteString(fmt.Sprintf("%s", input.Args.RoomId))

	id, err := a.resources.TaskResource.ScheduleTask(input.ScheduledTime, cmdStringBuilder.String())
	if err != nil {
		res := NewResponse(WithError(err, ScheduledTaskResponse{}))
		a.writeJSON(w, http.StatusUnprocessableEntity, res, nil)
		return
	}

	scheduledTaskDto := ScheduledTaskResponse{
		ID:      id,
		Command: cmdStringBuilder.String(),
		Cron:    input.ScheduledTime,
	}

	res := NewResponse(WithData(SCHEDULED_TASK, scheduledTaskDto))

	a.writeJSON(w, http.StatusAccepted, res, nil)
}

func (a *app) handleRemoveTask(w http.ResponseWriter, r *http.Request) {
	taskUUID := chi.URLParam(r, "uuid")

	taskId, err := uuid.Parse(taskUUID)
	if err != nil {
		res := NewResponse(WithError(err, ""))
		a.writeJSON(w, http.StatusInternalServerError, res, nil)
		return
	}

	a.resources.TaskResource.RemoveTaskByID(taskId)
	a.writeJSON(w, http.StatusNoContent, "", nil)
}
