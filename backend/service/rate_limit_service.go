package service

import (
	"context"
	"log/slog"

	"liveroom-battle/dao"
)

type RateLimitService struct {
	dao   *dao.RedisDao
	limit int64
}

func NewRateLimitService(dao *dao.RedisDao, limit int64) *RateLimitService {
	return &RateLimitService{dao: dao, limit: limit}
}

func (s *RateLimitService) IsLimited(ctx context.Context, roomID, userID string) (bool, error) {
	limited, err := s.dao.CheckChatRateLimit(ctx, roomID, userID, s.limit)
	if err != nil {
		slog.Error("rate limit check failed", "room_id", roomID, "user_id", userID, "err", err)
		return false, err
	}
	return limited, nil
}

func (s *RateLimitService) IncrLimitedCount(ctx context.Context, roomID string) (int64, error) {
	return s.dao.IncrLimitedCount(ctx, roomID)
}
