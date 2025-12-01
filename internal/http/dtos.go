package coco_http

import "github.com/google/uuid"

const (
	SCHEDULED_TASK = "scheduled_task" // refers to the type the client will receive
)

type ScheduledTaskDto struct {
	ID      uuid.UUID `json:"id"`
	Command string    `json:"command"`
	Cron    string    `json:"cron"`
}
