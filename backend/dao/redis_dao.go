package dao

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"liveroom-battle/model"
	"liveroom-battle/utils"

	"github.com/redis/go-redis/v9"
)

type RedisDao struct {
	rdb *redis.Client
}

func NewRedisDao(rdb *redis.Client) *RedisDao {
	return &RedisDao{rdb: rdb}
}

func (d *RedisDao) AddGiftScore(ctx context.Context, roomID, userID string, score int) error {
	key := utils.GiftRankKey(roomID)
	return d.rdb.ZIncrBy(ctx, key, float64(score), userID).Err()
}

func (d *RedisDao) GetTopRank(ctx context.Context, roomID string, topN int64) ([]model.RankItem, error) {
	key := utils.GiftRankKey(roomID)
	results, err := d.rdb.ZRevRangeWithScores(ctx, key, 0, topN-1).Result()
	if err != nil {
		return nil, err
	}
	items := make([]model.RankItem, 0, len(results))
	for _, z := range results {
		items = append(items, model.RankItem{
			UserID: z.Member.(string),
			Score:  int(z.Score),
		})
	}
	return items, nil
}

func (d *RedisDao) CheckChatRateLimit(ctx context.Context, roomID, userID string, limit int64) (bool, error) {
	key := utils.ChatRateLimitKey(roomID, userID)
	val, err := d.rdb.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if val == 1 {
		d.rdb.Expire(ctx, key, 1*time.Second)
	}
	return val > limit, nil
}

func (d *RedisDao) IncrLimitedCount(ctx context.Context, roomID string) (int64, error) {
	key := utils.LimitedCountKey(roomID)
	return d.rdb.Incr(ctx, key).Result()
}

func (d *RedisDao) GetLimitedCount(ctx context.Context, roomID string) (int, error) {
	key := utils.LimitedCountKey(roomID)
	val, err := d.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

func (d *RedisDao) AddOnlineUser(ctx context.Context, roomID, userID string) error {
	key := utils.RoomOnlineKey(roomID)
	return d.rdb.SAdd(ctx, key, userID).Err()
}

func (d *RedisDao) RemoveOnlineUser(ctx context.Context, roomID, userID string) error {
	key := utils.RoomOnlineKey(roomID)
	return d.rdb.SRem(ctx, key, userID).Err()
}

func (d *RedisDao) GetOnlineCount(ctx context.Context, roomID string) (int, error) {
	key := utils.RoomOnlineKey(roomID)
	val, err := d.rdb.SCard(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

func (d *RedisDao) RemoveRoomOnline(ctx context.Context, roomID string) {
	key := utils.RoomOnlineKey(roomID)
	if err := d.rdb.Del(ctx, key).Err(); err != nil {
		slog.Error("failed to remove room online set", "room_id", roomID, "err", err)
	}
}

func (d *RedisDao) IncrChatCount(ctx context.Context, roomID string) (int64, error) {
	key := utils.ChatCountKey(roomID)
	return d.rdb.Incr(ctx, key).Result()
}

func (d *RedisDao) GetChatCount(ctx context.Context, roomID string) (int64, error) {
	key := utils.ChatCountKey(roomID)
	val, err := d.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

func (d *RedisDao) IncrGiftCount(ctx context.Context, roomID string) (int64, error) {
	key := utils.GiftCountKey(roomID)
	return d.rdb.Incr(ctx, key).Result()
}

func (d *RedisDao) GetGiftCount(ctx context.Context, roomID string) (int64, error) {
	key := utils.GiftCountKey(roomID)
	val, err := d.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

func (d *RedisDao) UpsertRoomMeta(ctx context.Context, meta model.RoomMeta) error {
	key := utils.RoomMetaKey(meta.RoomID)
	exists, err := d.rdb.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return nil
	}
	pipe := d.rdb.Pipeline()
	pipe.HSet(ctx, key,
		"room_id", meta.RoomID,
		"title", meta.Title,
		"anchor_name", meta.AnchorName,
		"cover", meta.Cover,
		"status", meta.Status,
		"created_at", meta.CreatedAt,
	)
	timestamp, err := strconv.ParseInt(meta.CreatedAt, 10, 64)
	if err == nil {
		pipe.ZAdd(ctx, utils.LiveRoomsSetKey(), redis.Z{
			Score:  float64(timestamp),
			Member: meta.RoomID,
		})
	}
	_, err = pipe.Exec(ctx)
	return err
}

func (d *RedisDao) GetRoomMeta(ctx context.Context, roomID string) (*model.RoomMeta, error) {
	key := utils.RoomMetaKey(roomID)
	val, err := d.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(val) == 0 {
		return nil, nil
	}
	return &model.RoomMeta{
		RoomID:     val["room_id"],
		Title:      val["title"],
		AnchorName: val["anchor_name"],
		Cover:      val["cover"],
		Status:     val["status"],
		CreatedAt:  val["created_at"],
	}, nil
}

func (d *RedisDao) ListLiveRooms(ctx context.Context, limit int) ([]model.RoomMeta, error) {
	roomIDs, err := d.rdb.ZRevRange(ctx, utils.LiveRoomsSetKey(), 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}
	rooms := make([]model.RoomMeta, 0, len(roomIDs))
	for _, id := range roomIDs {
		meta, err := d.GetRoomMeta(ctx, id)
		if err != nil || meta == nil {
			continue
		}
		rooms = append(rooms, *meta)
	}
	return rooms, nil
}

func (d *RedisDao) CloseRoom(ctx context.Context, roomID string) error {
	key := utils.RoomMetaKey(roomID)
	pipe := d.rdb.Pipeline()
	pipe.HSet(ctx, key, "status", "closed")
	pipe.ZRem(ctx, utils.LiveRoomsSetKey(), roomID)
	_, err := pipe.Exec(ctx)
	return err
}

func (d *RedisDao) SetUserTeam(ctx context.Context, roomID, userID, team string) error {
	return d.rdb.Set(ctx, utils.PKUserTeamKey(roomID, userID), team, 0).Err()
}

func (d *RedisDao) GetUserTeam(ctx context.Context, roomID, userID string) (string, error) {
	val, err := d.rdb.Get(ctx, utils.PKUserTeamKey(roomID, userID)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (d *RedisDao) AddTeamUser(ctx context.Context, roomID, team string, userID string) error {
	var key string
	if team == "red" {
		key = utils.PKRedUsersKey(roomID)
	} else {
		key = utils.PKBlueUsersKey(roomID)
	}
	return d.rdb.SAdd(ctx, key, userID).Err()
}

func (d *RedisDao) GetTeamUserCount(ctx context.Context, roomID, team string) (int64, error) {
	var key string
	if team == "red" {
		key = utils.PKRedUsersKey(roomID)
	} else {
		key = utils.PKBlueUsersKey(roomID)
	}
	return d.rdb.SCard(ctx, key).Result()
}

func (d *RedisDao) AddPKScore(ctx context.Context, roomID, userID, team string, score int) error {
	pipe := d.rdb.Pipeline()
	var scoreKey string
	if team == "red" {
		scoreKey = utils.PKRedScoreKey(roomID)
	} else {
		scoreKey = utils.PKBlueScoreKey(roomID)
	}
	pipe.IncrBy(ctx, scoreKey, int64(score))
	pipe.ZIncrBy(ctx, utils.PKTeamRankKey(roomID, team), float64(score), userID)
	_, err := pipe.Exec(ctx)
	return err
}

func (d *RedisDao) GetPKScore(ctx context.Context, roomID string) (int64, int64, error) {
	redScore, err := d.rdb.Get(ctx, utils.PKRedScoreKey(roomID)).Result()
	if err == redis.Nil {
		redScore = "0"
	} else if err != nil {
		return 0, 0, err
	}
	blueScore, err := d.rdb.Get(ctx, utils.PKBlueScoreKey(roomID)).Result()
	if err == redis.Nil {
		blueScore = "0"
	} else if err != nil {
		return 0, 0, err
	}
	r, _ := strconv.ParseInt(redScore, 10, 64)
	b, _ := strconv.ParseInt(blueScore, 10, 64)
	return r, b, nil
}

func (d *RedisDao) SetPKState(ctx context.Context, state model.PKState) error {
	key := utils.PKStateKey(state.RoomID)
	return d.rdb.HSet(ctx, key,
		"room_id", state.RoomID,
		"status", state.Status,
		"red_score", state.RedScore,
		"blue_score", state.BlueScore,
		"red_users", state.RedUsers,
		"blue_users", state.BlueUsers,
		"winner", state.Winner,
		"start_time", state.StartTime,
		"end_time", state.EndTime,
	).Err()
}

func (d *RedisDao) GetPKState(ctx context.Context, roomID string) (*model.PKState, error) {
	key := utils.PKStateKey(roomID)
	val, err := d.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(val) == 0 {
		return nil, nil
	}
	redScore, _ := strconv.ParseInt(val["red_score"], 10, 64)
	blueScore, _ := strconv.ParseInt(val["blue_score"], 10, 64)
	redUsers, _ := strconv.ParseInt(val["red_users"], 10, 64)
	blueUsers, _ := strconv.ParseInt(val["blue_users"], 10, 64)
	return &model.PKState{
		RoomID:    val["room_id"],
		Status:    val["status"],
		RedScore:  redScore,
		BlueScore: blueScore,
		RedUsers:  redUsers,
		BlueUsers: blueUsers,
		Winner:    val["winner"],
		StartTime: val["start_time"],
		EndTime:   val["end_time"],
	}, nil
}

func (d *RedisDao) GetTeamRank(ctx context.Context, roomID, team string, topN int) ([]model.RankItem, error) {
	key := utils.PKTeamRankKey(roomID, team)
	results, err := d.rdb.ZRevRangeWithScores(ctx, key, 0, int64(topN-1)).Result()
	if err != nil {
		return nil, err
	}
	items := make([]model.RankItem, 0, len(results))
	for _, z := range results {
		items = append(items, model.RankItem{
			UserID: z.Member.(string),
			Score:  int(z.Score),
		})
	}
	return items, nil
}
