package controller

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"liveroom-battle/common"
	"liveroom-battle/model"
	"liveroom-battle/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSController struct {
	dispatcher *service.MessageDispatcher
	roomSvc    *service.RoomService
}

func NewWSController(dispatcher *service.MessageDispatcher, roomSvc *service.RoomService) *WSController {
	return &WSController{
		dispatcher: dispatcher,
		roomSvc:    roomSvc,
	}
}

func (c *WSController) HandleWS(ctx *gin.Context) {
	roomID := ctx.Query("room_id")
	userID := ctx.Query("user_id")
	if roomID == "" || userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "room_id and user_id required"})
		return
	}

	conn, err := common.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		slog.Error("websocket upgrade failed", "err", err)
		return
	}

	client := &model.Client{
		RoomID: roomID,
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	wsCtx := context.Background()
	c.roomSvc.Join(wsCtx, client)

	go c.writePump(client)
	go c.readPump(client, wsCtx)
}

func (c *WSController) readPump(client *model.Client, ctx context.Context) {
	defer func() {
		c.roomSvc.Leave(ctx, client)
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(4096)
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, raw, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				slog.Error("websocket read error", "err", err)
			}
			break
		}

		var msg model.Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			slog.Warn("invalid message format", "err", err)
			continue
		}

		msg.RoomID = client.RoomID
		msg.UserID = client.UserID

		c.dispatcher.Dispatch(ctx, client, &msg)
	}
}

func (c *WSController) writePump(client *model.Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				slog.Error("websocket write error", "err", err)
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
