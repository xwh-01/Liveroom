package hub

import (
	"log/slog"
	"sync"
	"time"

	"liveroom-battle/model"
)

type RoomHub struct {
	rooms map[string]map[*model.Client]bool
	mu    sync.RWMutex
}

func NewRoomHub() *RoomHub {
	return &RoomHub{
		rooms: make(map[string]map[*model.Client]bool),
	}
}

func (h *RoomHub) Join(roomID string, client *model.Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*model.Client]bool)
	}
	h.rooms[roomID][client] = true
	slog.Info("client joined", "room_id", roomID, "user_id", client.UserID)
}

func (h *RoomHub) Leave(roomID string, client *model.Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.rooms[roomID]; ok {
		delete(clients, client)
		close(client.Send)
		if len(clients) == 0 {
			delete(h.rooms, roomID)
		}
	}
	slog.Info("client left", "room_id", roomID, "user_id", client.UserID)
}

func (h *RoomHub) Broadcast(roomID string, messageType string, message []byte) {
	start := time.Now()
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.rooms[roomID]
	if !ok {
		return
	}
	target := len(clients)
	dropped := 0
	for client := range clients {
		select {
		case client.Send <- message:
		default:
			dropped++
			slog.Warn("client send buffer full, dropping message", "user_id", client.UserID)
		}
	}
	elapsed := time.Since(start)
	if target > 0 {
		slog.Info("broadcast finished",
			"room_id", roomID,
			"type", messageType,
			"target_clients", target,
			"dropped_clients", dropped,
			"latency_us", elapsed.Microseconds(),
		)
	}
}

func (h *RoomHub) SendToUser(roomID, userID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.rooms[roomID]
	if !ok {
		return
	}
	for client := range clients {
		if client.UserID == userID {
			select {
			case client.Send <- message:
			default:
			}
			return
		}
	}
}

func (h *RoomHub) OnlineCount(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.rooms[roomID]; ok {
		return len(clients)
	}
	return 0
}

func (h *RoomHub) RoomIDs() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	ids := make([]string, 0, len(h.rooms))
	for id := range h.rooms {
		ids = append(ids, id)
	}
	return ids
}
