package model

type RoomState struct {
	RoomID       string `json:"room_id"`
	OnlineCount  int    `json:"online_count"`
	LimitedCount int    `json:"limited_count"`
}
