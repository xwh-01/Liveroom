package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"liveroom-battle/dao"
	"liveroom-battle/hub"
	"liveroom-battle/model"
)

type RoomManageService struct {
	redisDao *dao.RedisDao
	hub      *hub.RoomHub
}

func NewRoomManageService(redisDao *dao.RedisDao, hub *hub.RoomHub) *RoomManageService {
	return &RoomManageService{redisDao: redisDao, hub: hub}
}

var defaultRooms = []model.RoomMeta{
	{
		RoomID:     "1001",
		Title:      "游戏开黑直播间",
		AnchorName: "小黑",
		Status:     "live",
	},
	{
		RoomID:     "1002",
		Title:      "音乐闲聊直播间",
		AnchorName: "Echo",
		Status:     "live",
	},
	{
		RoomID:     "1003",
		Title:      "学习摸鱼直播间",
		AnchorName: "阿码",
		Status:     "live",
	},
}

func (s *RoomManageService) EnsureDefaultRooms(ctx context.Context) error {
	now := time.Now().Unix()
	baseTS := strconv.FormatInt(now, 10)

	for i := range defaultRooms {
		defaultRooms[i].Cover = ""
		ts, _ := strconv.ParseInt(baseTS, 10, 64)
		defaultRooms[i].CreatedAt = strconv.FormatInt(ts+int64(i), 10)
		if err := s.redisDao.UpsertRoomMeta(ctx, defaultRooms[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s *RoomManageService) ListLiveRooms(ctx context.Context, limit int) ([]model.RoomMeta, error) {
	rooms, err := s.redisDao.ListLiveRooms(ctx, limit)
	if err != nil {
		return nil, err
	}
	for i := range rooms {
		rooms[i].OnlineCount = s.hub.OnlineCount(rooms[i].RoomID)
		chatCount, _ := s.redisDao.GetChatCount(ctx, rooms[i].RoomID)
		rooms[i].ChatCount = chatCount
		giftCount, _ := s.redisDao.GetGiftCount(ctx, rooms[i].RoomID)
		rooms[i].GiftCount = giftCount
	}
	return rooms, nil
}

func (s *RoomManageService) GetRoom(ctx context.Context, roomID string) (*model.RoomMeta, error) {
	meta, err := s.redisDao.GetRoomMeta(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if meta == nil {
		return nil, nil
	}
	meta.OnlineCount = s.hub.OnlineCount(roomID)
	chatCount, _ := s.redisDao.GetChatCount(ctx, roomID)
	meta.ChatCount = chatCount
	giftCount, _ := s.redisDao.GetGiftCount(ctx, roomID)
	meta.GiftCount = giftCount
	return meta, nil
}

func (s *RoomManageService) CloseRoom(ctx context.Context, roomID string) error {
	meta, err := s.redisDao.GetRoomMeta(ctx, roomID)
	if err != nil {
		return err
	}
	if meta == nil {
		return errors.New("room not found")
	}
	if meta.Status == "closed" {
		return errors.New("room already closed")
	}
	return s.redisDao.CloseRoom(ctx, roomID)
}
