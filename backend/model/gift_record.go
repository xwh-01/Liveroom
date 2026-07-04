package model

import "time"

type GiftRecord struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RecordID  string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"record_id"`
	RoomID    string    `gorm:"type:varchar(64);index;not null" json:"room_id"`
	UserID    string    `gorm:"type:varchar(128);not null" json:"user_id"`
	GiftID    string    `gorm:"type:varchar(32);not null" json:"gift_id"`
	Score     int       `gorm:"not null" json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

func (GiftRecord) TableName() string {
	return "gift_records"
}
