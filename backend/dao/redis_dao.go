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
