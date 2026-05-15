package cacherepo

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type FavoriteCache interface {
	ExistsCount(ctx context.Context, vid string) (bool, error)
	SetCount(ctx context.Context, vid string, count int64, ttl time.Duration) error
	GetCountUint64(ctx context.Context, vid string) (uint64, error)
	IncrCount(ctx context.Context, vid string) error
	DecrCount(ctx context.Context, vid string) error
	DelCount(ctx context.Context, vid string) error
	TryWarmupLock(ctx context.Context, vid string, ttl time.Duration) bool
}

type favoriteCache struct {
	rdb *redis.Client
}

func NewFavoriteCache(rdb *redis.Client) FavoriteCache {
	return &favoriteCache{rdb: rdb}
}

func buildFavoriteCountKey(vid string) string {
	return "favorite_count:video:" + vid
}

func buildFavoriteWarmupLockKey(vid string) string {
	return fmt.Sprintf("warmup_lock:favorite:%s", vid)
}

func (c *favoriteCache) ExistsCount(ctx context.Context, vid string) (bool, error) {
	exists, err := c.rdb.Exists(ctx, buildFavoriteCountKey(vid)).Result()
	return exists > 0, err
}

func (c *favoriteCache) SetCount(ctx context.Context, vid string, count int64, ttl time.Duration) error {
	return c.rdb.Set(ctx, buildFavoriteCountKey(vid), count, ttl).Err()
}

func (c *favoriteCache) GetCountUint64(ctx context.Context, vid string) (uint64, error) {
	return c.rdb.Get(ctx, buildFavoriteCountKey(vid)).Uint64()
}

func (c *favoriteCache) IncrCount(ctx context.Context, vid string) error {
	return c.rdb.Incr(ctx, buildFavoriteCountKey(vid)).Err()
}

func (c *favoriteCache) DecrCount(ctx context.Context, vid string) error {
	return c.rdb.Decr(ctx, buildFavoriteCountKey(vid)).Err()
}

func (c *favoriteCache) DelCount(ctx context.Context, vid string) error {
	return c.rdb.Del(ctx, buildFavoriteCountKey(vid)).Err()
}

func (c *favoriteCache) TryWarmupLock(ctx context.Context, vid string, ttl time.Duration) bool {
	ok, err := c.rdb.SetNX(ctx, buildFavoriteWarmupLockKey(vid), 1, ttl).Result()
	if err != nil {
		return false
	}
	return ok
}
