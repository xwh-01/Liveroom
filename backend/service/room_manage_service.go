package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
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

func (s *RoomManageService) CreateRoom(ctx context.Context, title, ownerName string) (*model.RoomMeta, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if ownerName == "" {
		ownerName = "anonymous"
	}

	roomID := generateRoomID()
	now := time.Now()
	createdAt := strconv.FormatInt(now.Unix(), 10)

	meta := model.RoomMeta{
		RoomID:    roomID,
		Title:     title,
		OwnerName: ownerName,
		Status:    "live",
		CreatedAt: createdAt,
	}

	if err := s.redisDao.CreateRoom(ctx, meta); err != nil {
		return nil, err
	}

	meta.OnlineCount = s.hub.OnlineCount(roomID)

	chatCount, _ := s.redisDao.GetChatCount(ctx, roomID)
	meta.ChatCount = chatCount
	giftCount, _ := s.redisDao.GetGiftCount(ctx, roomID)
	meta.GiftCount = giftCount

	return &meta, nil
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

func (s *RoomManageService) EnsureDefaultRoom(ctx context.Context) {
	meta, err := s.redisDao.GetRoomMeta(ctx, "1001")
	if err != nil || meta == nil {
		defaultMeta := model.RoomMeta{
			RoomID:    "1001",
			Title:     "默认直播间",
			OwnerName: "system",
			Status:    "live",
			CreatedAt: strconv.FormatInt(time.Now().Unix(), 10),
		}
		if err := s.redisDao.CreateRoom(ctx, defaultMeta); err != nil {
			return
		}
	}
}

func generateRoomID() string {
	now := time.Now().UnixMilli()
	r := rand.Intn(10000)
	return fmt.Sprintf("%d%04d", now, r)
}
