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
	roomDao  *dao.RoomDao
	hub      *hub.RoomHub
}

func NewRoomService(redisDao *dao.RedisDao, hub *hub.RoomHub, roomDao *dao.RoomDao) *RoomService {
	return &RoomService{redisDao: redisDao, hub: hub, roomDao: roomDao}
}

func (s *RoomService) Join(ctx context.Context, client *model.Client) {
	s.ensureRoomExists(ctx, client.RoomID)

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

func (s *RoomService) ensureRoomExists(ctx context.Context, roomID string) {
	if _, err := s.roomDao.CreateIfNotExists(ctx, roomID); err != nil {
		slog.Error("ensure room exists failed", "room_id", roomID, "err", err)
	}
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
	return &model.RoomState{
		RoomID:       roomID,
		OnlineCount:  onlineCount,
		LimitedCount: limitedCount,
	}
}
