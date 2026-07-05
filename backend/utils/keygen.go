package utils

import "fmt"

const (
	RoomOnlinePrefix    = "room:online"
	GiftRankPrefix      = "gift:rank"
	ChatRateLimitPrefix = "rate_limit:chat"
	LimitedCountPrefix  = "room:limited"
	ChatCountPrefix     = "room:chat_count"
	GiftCountPrefix     = "room:gift_count"
	RoomMetaPrefix      = "room:meta"
	LiveRoomsKey        = "room:live"
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

func ChatCountKey(roomID string) string {
	return fmt.Sprintf("%s:%s", ChatCountPrefix, roomID)
}

func GiftCountKey(roomID string) string {
	return fmt.Sprintf("%s:%s", GiftCountPrefix, roomID)
}

func RoomMetaKey(roomID string) string {
	return fmt.Sprintf("%s:%s", RoomMetaPrefix, roomID)
}

func LiveRoomsSetKey() string {
	return LiveRoomsKey
}
