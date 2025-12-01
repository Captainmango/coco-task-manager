package coco_http

import "github.com/google/uuid"

type TaskDto struct {
	ID      uuid.UUID `json:"id"`
	Command string    `json:"command"`
	Cron    string    `json:"cron"`
}
