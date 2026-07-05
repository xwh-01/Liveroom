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
	Content   string    `json:"content"`
	GiftType  string    `json:"gift_type"`
	GiftScore int       `json:"gift_score"`
	CreatedAt time.Time `json:"created_at"`
}

type PersistState struct {
	SubmittedCount      int64 `json:"submitted_count"`
	PersistedChatCount  int64 `json:"persisted_chat_count"`
	PersistedGiftCount  int64 `json:"persisted_gift_count"`
	PersistFailedCount  int64 `json:"persist_failed_count"`
	PersistDroppedCount int64 `json:"persist_dropped_count"`
	QueueLength         int   `json:"queue_length"`
	QueueCapacity       int   `json:"queue_capacity"`
}
