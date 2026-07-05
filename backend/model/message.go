package model

import "encoding/json"

type Message struct {
	Type   string          `json:"type"`
	RoomID string          `json:"room_id"`
	UserID string          `json:"user_id"`
	Data   json.RawMessage `json:"data"`
}

type ChatData struct {
	Content   string `json:"content"`
	Timestamp string `json:"timestamp,omitempty"`
	Team      string `json:"team,omitempty"`
}

type GiftData struct {
	GiftType  string `json:"gift_type"`
	GiftScore int    `json:"gift_score,omitempty"`
	Sender    string `json:"sender,omitempty"`
}

type OnlineData struct {
	Count int `json:"count"`
}

type RankData struct {
	Rankings []RankItem `json:"rankings"`
}

type RankItem struct {
	UserID string `json:"user_id"`
	Score  int    `json:"score"`
}

type SystemData struct {
	Content string `json:"content"`
}

func NewResponse(typ string, roomID string, userID string, data interface{}) ([]byte, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	msg := Message{
		Type:   typ,
		RoomID: roomID,
		UserID: userID,
		Data:   raw,
	}
	return json.Marshal(msg)
}
