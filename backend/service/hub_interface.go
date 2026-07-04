package service

type HubInterface interface {
	Broadcast(roomID string, message []byte)
	SendToUser(roomID, userID string, message []byte)
}
