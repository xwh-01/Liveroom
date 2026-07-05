package dao

import (
	"context"
	"database/sql"

	"liveroom-battle/model"
)

type RecordDao struct {
	db *sql.DB
}

func NewRecordDao(db *sql.DB) *RecordDao {
	return &RecordDao{db: db}
}

func (d *RecordDao) InsertChatRecord(ctx context.Context, r model.ChatRecord) error {
	_, err := d.db.ExecContext(ctx,
		"INSERT INTO chat_records (room_id, user_id, content, created_at) VALUES (?, ?, ?, ?)",
		r.RoomID, r.UserID, r.Content, r.CreatedAt,
	)
	return err
}

func (d *RecordDao) InsertGiftRecord(ctx context.Context, r model.GiftRecord) error {
	_, err := d.db.ExecContext(ctx,
		"INSERT INTO gift_records (room_id, user_id, gift_type, gift_score, created_at) VALUES (?, ?, ?, ?, ?)",
		r.RoomID, r.UserID, r.GiftType, r.GiftScore, r.CreatedAt,
	)
	return err
}

func (d *RecordDao) ListRecentChatRecords(ctx context.Context, roomID string, limit int) ([]model.ChatRecord, error) {
	rows, err := d.db.QueryContext(ctx,
		"SELECT id, room_id, user_id, content, created_at FROM chat_records WHERE room_id = ? ORDER BY created_at DESC LIMIT ?",
		roomID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.ChatRecord
	for rows.Next() {
		var r model.ChatRecord
		if err := rows.Scan(&r.ID, &r.RoomID, &r.UserID, &r.Content, &r.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, rows.Err()
}

func (d *RecordDao) ListRecentGiftRecords(ctx context.Context, roomID string, limit int) ([]model.GiftRecord, error) {
	rows, err := d.db.QueryContext(ctx,
		"SELECT id, room_id, user_id, gift_type, gift_score, created_at FROM gift_records WHERE room_id = ? ORDER BY created_at DESC LIMIT ?",
		roomID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.GiftRecord
	for rows.Next() {
		var r model.GiftRecord
		if err := rows.Scan(&r.ID, &r.RoomID, &r.UserID, &r.GiftType, &r.GiftScore, &r.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, rows.Err()
}
