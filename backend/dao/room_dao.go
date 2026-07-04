package dao

import (
	"context"

	"liveroom-battle/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoomDao struct {
	db *gorm.DB
}

func NewRoomDao(db *gorm.DB) *RoomDao {
	return &RoomDao{db: db}
}

func (d *RoomDao) CreateIfNotExists(ctx context.Context, roomID string) (*model.Room, error) {
	room := &model.Room{
		RoomID:     roomID,
		Title:      "Live Room " + roomID,
		AnchorName: "主播 " + roomID,
		Status:     1,
	}
	result := d.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "room_id"}},
		DoNothing: true,
	}).Create(room)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		var existing model.Room
		if err := d.db.WithContext(ctx).Where("room_id = ?", roomID).First(&existing).Error; err != nil {
			return nil, err
		}
		return &existing, nil
	}
	return room, nil
}

func (d *RoomDao) GetByRoomID(ctx context.Context, roomID string) (*model.Room, error) {
	var room model.Room
	err := d.db.WithContext(ctx).Where("room_id = ?", roomID).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}
