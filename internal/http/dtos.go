package coco_http

import "github.com/google/uuid"

const (
	SCHEDULED_TASK = "scheduled_task" // refers to the type the client will receive
	TASK           = "task"
)

type ScheduledTaskResponse struct {
	ID      uuid.UUID `json:"id"`
	Command string    `json:"command"`
	Cron    string    `json:"cron"`
}

type TaskResponse struct {
	Slug string   `json:"task_id"`
	Args []string `json:"args"`
}

type ScheduleTaskRequest struct {
	TaskId        string `json:"task_id"`
	ScheduledTime string `json:"scheduled_time"`
	Args          struct {
		RoomId string `json:"room_id"`
	} `json:"args"`
}
