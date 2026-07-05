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

type HubInterface interface {
	Broadcast(roomID string, messageType string, message []byte)
	SendToUser(roomID string, userID string, message []byte)
}

type ChatService struct {
	rateLimitSvc *RateLimitService
	hub          HubInterface
	redisDao     *dao.RedisDao
	persistSvc   *PersistService
}

func NewChatService(rateLimitSvc *RateLimitService, hub HubInterface, redisDao *dao.RedisDao, persistSvc *PersistService) *ChatService {
	return &ChatService{
		rateLimitSvc: rateLimitSvc,
		hub:          hub,
		redisDao:     redisDao,
		persistSvc:   persistSvc,
	}
}

func (s *ChatService) HandleChat(ctx context.Context, client *model.Client, msg *model.Message) {
	var chatData model.ChatData
	if err := json.Unmarshal(msg.Data, &chatData); err != nil {
		slog.Error("unmarshal chat data failed", "err", err)
		return
	}

	limited, err := s.rateLimitSvc.IsLimited(ctx, msg.RoomID, msg.UserID)
	if err != nil {
		return
	}

	if limited {
		s.rateLimitSvc.IncrLimitedCount(ctx, msg.RoomID)
		reply, _ := model.NewResponse("system", msg.RoomID, "", model.SystemData{
			Content: "你发送太快了，已被限流",
		})
		s.hub.SendToUser(msg.RoomID, msg.UserID, reply)
		return
	}

	chatData.Timestamp = utils.NowStr()
	payload, _ := model.NewResponse("chat", msg.RoomID, msg.UserID, chatData)
	s.hub.Broadcast(msg.RoomID, "chat", payload)

	if _, err := s.redisDao.IncrChatCount(ctx, msg.RoomID); err != nil {
		slog.Error("incr chat count failed", "err", err)
	}

	if s.persistSvc != nil {
		s.persistSvc.Submit(model.PersistEvent{
			Type:      "chat",
			RoomID:    msg.RoomID,
			UserID:    msg.UserID,
			Content:   chatData.Content,
			CreatedAt: time.Now(),
		})
	}
}
