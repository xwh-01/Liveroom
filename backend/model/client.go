package model

import "github.com/gorilla/websocket"

type Client struct {
	RoomID string
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}
