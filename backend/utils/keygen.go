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
	PKStatePrefix       = "room:pk:state"
	PKScorePrefix       = "room:pk:score"
	PKUsersPrefix       = "room:pk:users"
	PKUserTeamPrefix    = "room:pk:user:team"
	PKTeamRankPrefix    = "room:pk:rank"
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

func PKStateKey(roomID string) string {
	return fmt.Sprintf("%s:%s", PKStatePrefix, roomID)
}

func PKRedScoreKey(roomID string) string {
	return fmt.Sprintf("%s:%s:red", PKScorePrefix, roomID)
}

func PKBlueScoreKey(roomID string) string {
	return fmt.Sprintf("%s:%s:blue", PKScorePrefix, roomID)
}

func PKRedUsersKey(roomID string) string {
	return fmt.Sprintf("%s:%s:red", PKUsersPrefix, roomID)
}

func PKBlueUsersKey(roomID string) string {
	return fmt.Sprintf("%s:%s:blue", PKUsersPrefix, roomID)
}

func PKUserTeamKey(roomID, userID string) string {
	return fmt.Sprintf("%s:%s:%s:team", PKUserTeamPrefix, roomID, userID)
}

func PKTeamRankKey(roomID, team string) string {
	return fmt.Sprintf("%s:%s:%s", PKTeamRankPrefix, roomID, team)
}
