package model

type PKState struct {
	RoomID           string `json:"room_id"`
	Status           string `json:"status"`
	RedScore         int64  `json:"red_score"`
	BlueScore        int64  `json:"blue_score"`
	RedUsers         int64  `json:"red_users"`
	BlueUsers        int64  `json:"blue_users"`
	Winner           string `json:"winner"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	RemainingSeconds int    `json:"remaining_seconds"`
}

type TeamData struct {
	Team string `json:"team"`
}

type PKStartReq struct {
	RoomID          string `json:"room_id"`
	DurationSeconds int    `json:"duration_seconds"`
}
