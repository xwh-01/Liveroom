package service

import (
	"context"
	"log/slog"

	"liveroom-battle/model"
)

type MessageHandler func(ctx context.Context, client *model.Client, msg *model.Message)

type MessageDispatcher struct {
	handlers map[string]MessageHandler
}

func NewMessageDispatcher() *MessageDispatcher {
	return &MessageDispatcher{
		handlers: make(map[string]MessageHandler),
	}
}

func (d *MessageDispatcher) Register(msgType string, handler MessageHandler) {
	d.handlers[msgType] = handler
}

func (d *MessageDispatcher) Dispatch(ctx context.Context, client *model.Client, msg *model.Message) {
	handler, ok := d.handlers[msg.Type]
	if !ok {
		slog.Warn("unknown message type", "type", msg.Type)
		return
	}
	handler(ctx, client, msg)
}
