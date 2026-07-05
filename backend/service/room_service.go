package service

import (
	"context"
	"log/slog"

	"liveroom-battle/dao"
	"liveroom-battle/hub"
	"liveroom-battle/model"
)

type RoomService struct {
	redisDao *dao.RedisDao
	hub      *hub.RoomHub
}

func NewRoomService(redisDao *dao.RedisDao, hub *hub.RoomHub) *RoomService {
	return &RoomService{redisDao: redisDao, hub: hub}
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
	s.hub.Broadcast(roomID, "online", msg)
}

func (s *RoomService) GetRoomState(ctx context.Context, roomID string) *model.RoomState {
	onlineCount := s.hub.OnlineCount(roomID)

	limitedCount, err := s.redisDao.GetLimitedCount(ctx, roomID)
	if err != nil {
		slog.Error("get limited count failed", "room_id", roomID, "err", err)
	}

	chatCount, err := s.redisDao.GetChatCount(ctx, roomID)
	if err != nil {
		slog.Error("get chat count failed", "room_id", roomID, "err", err)
	}

	giftCount, err := s.redisDao.GetGiftCount(ctx, roomID)
	if err != nil {
		slog.Error("get gift count failed", "room_id", roomID, "err", err)
	}

	return &model.RoomState{
		RoomID:      roomID,
		OnlineCount: onlineCount,
		LimitedCount: limitedCount,
		ChatCount:   chatCount,
		GiftCount:   giftCount,
	}
}
