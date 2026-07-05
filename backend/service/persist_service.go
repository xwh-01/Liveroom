package service

import (
	"context"
	"log/slog"
	"sync/atomic"

	"liveroom-battle/dao"
	"liveroom-battle/model"
)

type PersistService struct {
	recordDao *dao.RecordDao
	queue     chan model.PersistEvent

	submittedCount      atomic.Int64
	persistedChatCount  atomic.Int64
	persistedGiftCount  atomic.Int64
	failedCount         atomic.Int64
	droppedCount        atomic.Int64
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
		s.submittedCount.Add(1)
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

func (s *PersistService) State() model.PersistState {
	return model.PersistState{
		SubmittedCount:      s.submittedCount.Load(),
		PersistedChatCount:  s.persistedChatCount.Load(),
		PersistedGiftCount:  s.persistedGiftCount.Load(),
		PersistFailedCount:  s.failedCount.Load(),
		PersistDroppedCount: s.droppedCount.Load(),
		QueueLength:         len(s.queue),
		QueueCapacity:       cap(s.queue),
	}
}

func (s *PersistService) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-s.queue:
			s.handleEvent(ctx, event)
		}
	}
}

func (s *PersistService) handleEvent(ctx context.Context, event model.PersistEvent) {
	switch event.Type {
	case "chat":
		err := s.recordDao.InsertChatRecord(ctx, model.ChatRecord{
			RoomID:    event.RoomID,
			UserID:    event.UserID,
			Content:   event.Content,
			CreatedAt: event.CreatedAt,
		})
		if err != nil {
			s.failedCount.Add(1)
			slog.Error("persist chat record failed", "err", err, "room_id", event.RoomID, "user_id", event.UserID)
			return
		}
		s.persistedChatCount.Add(1)
	case "gift":
		err := s.recordDao.InsertGiftRecord(ctx, model.GiftRecord{
			RoomID:    event.RoomID,
			UserID:    event.UserID,
			GiftType:  event.GiftType,
			GiftScore: event.GiftScore,
			CreatedAt: event.CreatedAt,
		})
		if err != nil {
			s.failedCount.Add(1)
			slog.Error("persist gift record failed", "err", err, "room_id", event.RoomID, "user_id", event.UserID)
			return
		}
		s.persistedGiftCount.Add(1)
	}
}
