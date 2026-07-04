package model

type RoomState struct {
	RoomID       string     `json:"room_id"`
	OnlineCount  int        `json:"online_count"`
	LimitedCount int        `json:"limited_count"`
	ChatCount    int64      `json:"chat_count"`
	GiftCount    int64      `json:"gift_count"`
	Rankings     []RankItem `json:"rankings,omitempty"`
}
