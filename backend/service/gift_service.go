package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"liveroom-battle/dao"
	"liveroom-battle/model"
)

type GiftService struct {
	redisDao   *dao.RedisDao
	hub        HubInterface
	persistSvc *PersistService
	pkSvc      *PKService
}

func NewGiftService(redisDao *dao.RedisDao, hub HubInterface, persistSvc *PersistService, pkSvc *PKService) *GiftService {
	return &GiftService{redisDao: redisDao, hub: hub, persistSvc: persistSvc, pkSvc: pkSvc}
}

func (s *GiftService) HandleGift(ctx context.Context, client *model.Client, msg *model.Message) {
	var giftData model.GiftData
	if err := json.Unmarshal(msg.Data, &giftData); err != nil {
		slog.Error("unmarshal gift data failed", "err", err)
		return
	}

	score := model.GetGiftScore(giftData.GiftType)
	if score == 0 {
		slog.Warn("invalid gift type", "gift_type", giftData.GiftType)
		return
	}

	if err := s.redisDao.AddGiftScore(ctx, msg.RoomID, msg.UserID, score); err != nil {
		slog.Error("add gift score failed", "err", err, "room_id", msg.RoomID, "user_id", msg.UserID)
		return
	}

	giftData.GiftScore = score
	giftData.Sender = msg.UserID
	payload, _ := model.NewResponse("gift", msg.RoomID, msg.UserID, giftData)
	s.hub.Broadcast(msg.RoomID, "gift", payload)

	rankings, err := s.redisDao.GetTopRank(ctx, msg.RoomID, 10)
	if err != nil {
		slog.Error("get top rank failed", "err", err)
		return
	}

	rankPayload, _ := model.NewResponse("rank", msg.RoomID, "", model.RankData{Rankings: rankings})
	s.hub.Broadcast(msg.RoomID, "rank", rankPayload)

	if _, err := s.redisDao.IncrGiftCount(ctx, msg.RoomID); err != nil {
		slog.Error("incr gift count failed", "err", err)
	}

	if s.persistSvc != nil {
		s.persistSvc.Submit(model.PersistEvent{
			Type:      "gift",
			RoomID:    msg.RoomID,
			UserID:    msg.UserID,
			GiftType:  giftData.GiftType,
			GiftScore: score,
			CreatedAt: time.Now(),
		})
	}

	if s.pkSvc != nil {
		pkState, err := s.pkSvc.AddScore(ctx, msg.RoomID, msg.UserID, score)
		if err == nil && pkState != nil {
			s.broadcastPKState(ctx, msg.RoomID)
		}
	}
}

func (s *GiftService) broadcastPKState(ctx context.Context, roomID string) {
	state, err := s.pkSvc.GetPKState(ctx, roomID)
	if err != nil || state == nil {
		return
	}
	payload, _ := model.NewResponse("pk_state", roomID, "", state)
	s.hub.Broadcast(roomID, "pk_state", payload)
}
