package model

import "time"

type Room struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RoomID     string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"room_id"`
	Title      string    `gorm:"type:varchar(255);not null" json:"title"`
	AnchorName string    `gorm:"type:varchar(128)" json:"anchor_name"`
	Status     int       `gorm:"type:tinyint;default:1" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Room) TableName() string {
	return "rooms"
}
