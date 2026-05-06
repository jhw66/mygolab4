package cacherepo

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const (
	rankZSetKey       = "rank:video:hot"
	rankDirtyVideoKey = "rank:dirty_videos"
)

type RankCache interface {
	ZAddScore(ctx context.Context, vid string, score float64) error
	ZBatchAdd(ctx context.Context, entries []redis.Z) error
	GetTopVideoIDs(ctx context.Context, limit int) ([]string, error)
	MarkDirty(ctx context.Context, vid string) error
	ListDirty(ctx context.Context) ([]string, error)
	RemoveDirtyBatch(ctx context.Context, vids []interface{}) error
	RemoveVideo(ctx context.Context, vid string) error
}

type rankCache struct {
	rdb *redis.Client
}

func NewRankCache(rdb *redis.Client) RankCache {
	return &rankCache{rdb: rdb}
}

func (c *rankCache) ZAddScore(ctx context.Context, vid string, score float64) error {
	return c.rdb.ZAdd(ctx, rankZSetKey, redis.Z{Score: score, Member: vid}).Err()
}

func (c *rankCache) ZBatchAdd(ctx context.Context, entries []redis.Z) error {
	pipe := c.rdb.Pipeline()
	for _, entry := range entries {
		pipe.ZAdd(ctx, rankZSetKey, entry)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (c *rankCache) GetTopVideoIDs(ctx context.Context, limit int) ([]string, error) {
	return c.rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   rankZSetKey,
		Start: 0,
		Stop:  int64(limit - 1),
		Rev:   true,
	}).Result()
}

func (c *rankCache) MarkDirty(ctx context.Context, vid string) error {
	return c.rdb.SAdd(ctx, rankDirtyVideoKey, vid).Err()
}

func (c *rankCache) ListDirty(ctx context.Context) ([]string, error) {
	return c.rdb.SMembers(ctx, rankDirtyVideoKey).Result()
}

func (c *rankCache) RemoveDirtyBatch(ctx context.Context, vids []interface{}) error {
	if len(vids) == 0 {
		return nil
	}
	return c.rdb.SRem(ctx, rankDirtyVideoKey, vids...).Err()
}

func (c *rankCache) RemoveVideo(ctx context.Context, vid string) error {
	pipe := c.rdb.Pipeline()
	pipe.ZRem(ctx, rankZSetKey, vid)
	pipe.SRem(ctx, rankDirtyVideoKey, vid)
	_, err := pipe.Exec(ctx)
	return err
}
