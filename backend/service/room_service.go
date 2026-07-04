package service

import (
	"context"
	"log/slog"

	"liveroom-battle/dao"
	"liveroom-battle/hub"
	"liveroom-battle/model"
)

type RoomService struct {
	redisDao     *dao.RedisDao
	hub          *hub.RoomHub
	metricsStore *model.MetricsStore
}

func NewRoomService(redisDao *dao.RedisDao, hub *hub.RoomHub, metricsStore *model.MetricsStore) *RoomService {
	return &RoomService{redisDao: redisDao, hub: hub, metricsStore: metricsStore}
}

func (s *RoomService) Join(ctx context.Context, client *model.Client) {
	s.hub.Join(client.RoomID, client)
	if err := s.redisDao.AddOnlineUser(ctx, client.RoomID, client.UserID); err != nil {
		slog.Error("add online user failed", "room_id", client.RoomID, "user_id", client.UserID, "err", err)
	}
	s.broadcastOnline(ctx, client.RoomID)
}

func (s *RoomService) Leave(ctx context.Context, client *model.Client) {
	s.hub.Leave(client.RoomID, client)
	if err := s.redisDao.RemoveOnlineUser(ctx, client.RoomID, client.UserID); err != nil {
		slog.Error("remove online user failed", "room_id", client.RoomID, "user_id", client.UserID, "err", err)
	}
	s.broadcastOnline(ctx, client.RoomID)
}

func (s *RoomService) broadcastOnline(ctx context.Context, roomID string) {
	count := s.hub.OnlineCount(roomID)
	msg, err := model.NewResponse("online", roomID, "", model.OnlineData{Count: count})
	if err != nil {
		slog.Error("marshal online message failed", "err", err)
		return
	}
	s.hub.Broadcast(roomID, msg)
}

func (s *RoomService) GetRoomState(ctx context.Context, roomID string) *model.RoomState {
	onlineCount := s.hub.OnlineCount(roomID)

	limitedCount, err := s.redisDao.GetLimitedCount(ctx, roomID)
	if err != nil {
		slog.Error("get limited count failed", "room_id", roomID, "err", err)
	}

	chatCount := s.metricsStore.GetChatCount(roomID)
	giftCount := s.metricsStore.GetGiftCount(roomID)

	rankings, err := s.redisDao.GetTopRank(ctx, roomID, 10)
	if err != nil {
		slog.Error("get rank failed", "room_id", roomID, "err", err)
	}
	if rankings == nil {
		rankings = []model.RankItem{}
	}

	return &model.RoomState{
		RoomID:          roomID,
		OnlineCount:     onlineCount,
		ConnectionCount: onlineCount,
		LimitedCount:    limitedCount,
		ChatCount:       chatCount,
		GiftCount:       giftCount,
		Rankings:        rankings,
	}
}
