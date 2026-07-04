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
