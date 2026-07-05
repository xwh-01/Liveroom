package model

import "time"

type ChatRecord struct {
	ID        int64     `json:"id"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type GiftRecord struct {
	ID        int64     `json:"id"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	GiftType  string    `json:"gift_type"`
	GiftScore int       `json:"gift_score"`
	CreatedAt time.Time `json:"created_at"`
}

type PersistEvent struct {
	Type      string    `json:"type"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content,omitempty"`
	GiftType  string    `json:"gift_type,omitempty"`
	GiftScore int       `json:"gift_score,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
