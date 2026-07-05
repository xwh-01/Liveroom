package model

type RoomMeta struct {
	RoomID      string `json:"room_id"`
	Title       string `json:"title"`
	OwnerName   string `json:"owner_name"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	OnlineCount int    `json:"online_count"`
	ChatCount   int64  `json:"chat_count"`
	GiftCount   int64  `json:"gift_count"`
}
