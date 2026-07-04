package model

import "time"

type ChatMessage struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	MessageID string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"message_id"`
	RoomID    string    `gorm:"type:varchar(64);index;not null" json:"room_id"`
	UserID    string    `gorm:"type:varchar(128);not null" json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}
