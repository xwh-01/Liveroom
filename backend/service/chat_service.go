package service

import (
	"context"
	"encoding/json"
	"log/slog"

	"liveroom-battle/model"
	"liveroom-battle/utils"
)

type ChatService struct {
	rateLimitSvc *RateLimitService
	hub          HubInterface
}

type HubInterface interface {
	Broadcast(roomID string, message []byte)
	SendToUser(roomID, userID string, message []byte)
}

func NewChatService(rateLimitSvc *RateLimitService, hub HubInterface) *ChatService {
	return &ChatService{
		rateLimitSvc: rateLimitSvc,
		hub:          hub,
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
	s.hub.Broadcast(msg.RoomID, payload)
}
