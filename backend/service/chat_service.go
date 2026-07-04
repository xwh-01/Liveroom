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

type ChatService struct {
	rateLimitSvc   *RateLimitService
	hub            HubInterface
	chatMessageDao *dao.ChatMessageDao
}

func NewChatService(rateLimitSvc *RateLimitService, hub HubInterface, chatMessageDao *dao.ChatMessageDao) *ChatService {
	return &ChatService{
		rateLimitSvc:   rateLimitSvc,
		hub:            hub,
		chatMessageDao: chatMessageDao,
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

	s.saveChatMessage(ctx, msg.RoomID, msg.UserID, chatData.Content)
}

func (s *ChatService) saveChatMessage(ctx context.Context, roomID, userID, content string) {
	record := &model.ChatMessage{
		MessageID: utils.NewUUID(),
		RoomID:    roomID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}
	if err := s.chatMessageDao.Save(ctx, record); err != nil {
		slog.Error("save chat message failed", "err", err, "room_id", roomID, "user_id", userID)
	}
}
