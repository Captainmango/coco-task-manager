package coco_http

import "net/http"

func (a *app) handleGetScheduledTasks(w http.ResponseWriter, r *http.Request) {
	entries, err := a.resources.TaskResource.GetAllCrontabEntries()

	// Check what the error is at some point
	if err != nil {
		errRes := NewResponse(
			WithError(
				err,
				[]ScheduledTaskDto{},
				tMeta{"hello": "world"},
			),
		)

		a.writeJSON(w, http.StatusBadRequest, errRes, nil)
		return
	}

	var out []ScheduledTaskDto = []ScheduledTaskDto{}

	for _, item := range entries {
		i := ScheduledTaskDto{
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
