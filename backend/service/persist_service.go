package service

import (
	"context"
	"log/slog"
	"sync/atomic"

	"liveroom-battle/dao"
	"liveroom-battle/model"
)

type PersistService struct {
	recordDao    *dao.RecordDao
	queue        chan model.PersistEvent
	droppedCount atomic.Int64
}

func NewPersistService(recordDao *dao.RecordDao, queueSize int) *PersistService {
	return &PersistService{
		recordDao: recordDao,
		queue:     make(chan model.PersistEvent, queueSize),
	}
}

func (s *PersistService) Start(ctx context.Context, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go s.worker(ctx)
	}
}

func (s *PersistService) Submit(event model.PersistEvent) bool {
	select {
	case s.queue <- event:
		return true
	default:
		s.droppedCount.Add(1)
		slog.Warn("persist queue full, event dropped",
			"type", event.Type,
			"room_id", event.RoomID,
			"user_id", event.UserID,
			"dropped_total", s.droppedCount.Load(),
		)
		return false
	}
}

func (s *PersistService) DroppedCount() int64 {
	return s.droppedCount.Load()
}

func (s *PersistService) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-s.queue:
			s.handle(ctx, event)
		}
	}
}

func (s *PersistService) handle(ctx context.Context, event model.PersistEvent) {
	switch event.Type {
	case "chat":
		record := model.ChatRecord{
			RoomID:    event.RoomID,
			UserID:    event.UserID,
			Content:   event.Content,
			CreatedAt: event.CreatedAt,
		}
		if err := s.recordDao.InsertChatRecord(ctx, record); err != nil {
			slog.Error("persist chat record failed", "err", err,
				"room_id", event.RoomID, "user_id", event.UserID)
		}
	case "gift":
		record := model.GiftRecord{
			RoomID:    event.RoomID,
			UserID:    event.UserID,
			GiftType:  event.GiftType,
			GiftScore: event.GiftScore,
			CreatedAt: event.CreatedAt,
		}
		if err := s.recordDao.InsertGiftRecord(ctx, record); err != nil {
			slog.Error("persist gift record failed", "err", err,
				"room_id", event.RoomID, "user_id", event.UserID)
		}
	default:
		slog.Warn("unknown persist event type", "type", event.Type)
	}
}
