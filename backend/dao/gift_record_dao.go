package dao

import (
	"context"

	"liveroom-battle/model"

	"gorm.io/gorm"
)

type GiftRecordDao struct {
	db *gorm.DB
}

func NewGiftRecordDao(db *gorm.DB) *GiftRecordDao {
	return &GiftRecordDao{db: db}
}

func (d *GiftRecordDao) Save(ctx context.Context, record *model.GiftRecord) error {
	return d.db.WithContext(ctx).Create(record).Error
}

func (d *GiftRecordDao) ListRecent(ctx context.Context, roomID string, limit int) ([]model.GiftRecord, error) {
	var records []model.GiftRecord
	err := d.db.WithContext(ctx).
		Where("room_id = ?", roomID).
		Order("created_at DESC").
		Limit(limit).
		Find(&records).Error
	return records, err
}
