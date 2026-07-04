package dao

import (
	"context"

	"liveroom-battle/model"

	"gorm.io/gorm"
)

type ChatMessageDao struct {
	db *gorm.DB
}

func NewChatMessageDao(db *gorm.DB) *ChatMessageDao {
	return &ChatMessageDao{db: db}
}

func (d *ChatMessageDao) Save(ctx context.Context, msg *model.ChatMessage) error {
	return d.db.WithContext(ctx).Create(msg).Error
}

func (d *ChatMessageDao) ListRecent(ctx context.Context, roomID string, limit int) ([]model.ChatMessage, error) {
	var messages []model.ChatMessage
	err := d.db.WithContext(ctx).
		Where("room_id = ?", roomID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}
