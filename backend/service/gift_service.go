package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"liveroom-battle/dao"
	"liveroom-battle/model"
	"liveroom-battle/utils"
)

type GiftService struct {
	redisDao      *dao.RedisDao
	hub           HubInterface
	giftRecordDao *dao.GiftRecordDao
}

func NewGiftService(redisDao *dao.RedisDao, hub HubInterface, giftRecordDao *dao.GiftRecordDao) *GiftService {
	return &GiftService{redisDao: redisDao, hub: hub, giftRecordDao: giftRecordDao}
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
	s.hub.Broadcast(msg.RoomID, payload)

	rankings, err := s.redisDao.GetTopRank(ctx, msg.RoomID, 10)
	if err != nil {
		slog.Error("get top rank failed", "err", err)
		return
	}

	rankPayload, _ := model.NewResponse("rank", msg.RoomID, "", model.RankData{Rankings: rankings})
	s.hub.Broadcast(msg.RoomID, rankPayload)

	s.saveGiftRecord(ctx, msg.RoomID, msg.UserID, giftData.GiftType, score)
}

func (s *GiftService) saveGiftRecord(ctx context.Context, roomID, userID, giftType string, score int) {
	record := &model.GiftRecord{
		RecordID:  utils.NewUUID(),
		RoomID:    roomID,
		UserID:    userID,
		GiftID:    giftType,
		Score:     score,
		CreatedAt: time.Now(),
	}
	if err := s.giftRecordDao.Save(ctx, record); err != nil {
		slog.Error("save gift record failed", "err", err, "room_id", roomID, "user_id", userID)
	}
}
