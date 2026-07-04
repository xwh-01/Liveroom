package utils

import "fmt"

const (
	RoomOnlinePrefix    = "room:online"
	GiftRankPrefix      = "gift:rank"
	ChatRateLimitPrefix = "rate_limit:chat"
	LimitedCountPrefix  = "room:limited"
)

func GiftRankKey(roomID string) string {
	return fmt.Sprintf("%s:%s", GiftRankPrefix, roomID)
}

func ChatRateLimitKey(roomID, userID string) string {
	return fmt.Sprintf("%s:%s:%s", ChatRateLimitPrefix, roomID, userID)
}

func LimitedCountKey(roomID string) string {
	return fmt.Sprintf("%s:%s", LimitedCountPrefix, roomID)
}

func RoomOnlineKey(roomID string) string {
	return fmt.Sprintf("%s:%s", RoomOnlinePrefix, roomID)
}
