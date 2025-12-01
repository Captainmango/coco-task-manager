package coco_http

import "github.com/google/uuid"

type ValidResponseData[T any] interface {
	[]T
}

type BaseResponse[T any] struct {
	Type   string
	Data   []T
}

type TaskDto struct {
	ID      uuid.UUID
	Command string
	Cron    string
}

type GetTasksResponse struct{ BaseResponse[TaskDto] }
func NewGetTasksResponse(data []TaskDto) GetTasksResponse {
	return GetTasksResponse{
		BaseResponse: BaseResponse[TaskDto]{
			Type: "task",
			Data: data,
		},
	}
}
