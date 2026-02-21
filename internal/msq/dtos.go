package msq

type StartGamePayload struct {
	RoomId string `json:"room_id"`
}

func (sgp StartGamePayload) GetRoutingKey() string {
	return "coco_tasks.start_game"
}