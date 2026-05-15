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
	ExistsCommentFavoriteCount(ctx context.Context, commentID string) (bool, error)
	SetCommentFavoriteCount(ctx context.Context, commentID string, count int64, ttl time.Duration) error
	GetCommentFavoriteCount(ctx context.Context, commentID string) (int64, error)
	IncrCommentFavoriteCount(ctx context.Context, commentID string) error
	DecrCommentFavoriteCount(ctx context.Context, commentID string) error
	DelCommentFavoriteCount(ctx context.Context, commentID string) error
	MarkDirtyCommentFavoriteCount(ctx context.Context, commentID string) error
	ListDirtyCommentFavoriteCount(ctx context.Context) ([]string, error)
	RemoveDirtyCommentFavoriteCountBatch(ctx context.Context, commentIDs []interface{}) error
	TryCommentFavoriteWarmupLock(ctx context.Context, commentID string, ttl time.Duration) bool
}

type favoriteCache struct {
	rdb *redis.Client
}

const commentDirtyFavoriteCountKey = "comment:dirty_favorite_count"

func NewFavoriteCache(rdb *redis.Client) FavoriteCache {
	return &favoriteCache{rdb: rdb}
}

func buildFavoriteCountKey(vid string) string {
	return "favorite_count:video:" + vid
}

func buildFavoriteWarmupLockKey(vid string) string {
	return fmt.Sprintf("warmup_lock:favorite:%s", vid)
}

func buildCommentFavoriteCountKey(commentID string) string {
	return "comment_favorite_count:" + commentID
}

func buildCommentFavoriteWarmupLockKey(commentID string) string {
	return fmt.Sprintf("warmup_lock:comment_favorite:%s", commentID)
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

func (c *favoriteCache) ExistsCommentFavoriteCount(ctx context.Context, commentID string) (bool, error) {
	exists, err := c.rdb.Exists(ctx, buildCommentFavoriteCountKey(commentID)).Result()
	return exists > 0, err
}

func (c *favoriteCache) SetCommentFavoriteCount(ctx context.Context, commentID string, count int64, ttl time.Duration) error {
	return c.rdb.Set(ctx, buildCommentFavoriteCountKey(commentID), count, ttl).Err()
}

func (c *favoriteCache) GetCommentFavoriteCount(ctx context.Context, commentID string) (int64, error) {
	return c.rdb.Get(ctx, buildCommentFavoriteCountKey(commentID)).Int64()
}

func (c *favoriteCache) IncrCommentFavoriteCount(ctx context.Context, commentID string) error {
	return c.rdb.Incr(ctx, buildCommentFavoriteCountKey(commentID)).Err()
}

func (c *favoriteCache) DecrCommentFavoriteCount(ctx context.Context, commentID string) error {
	return c.rdb.Decr(ctx, buildCommentFavoriteCountKey(commentID)).Err()
}

func (c *favoriteCache) DelCommentFavoriteCount(ctx context.Context, commentID string) error {
	return c.rdb.Del(ctx, buildCommentFavoriteCountKey(commentID)).Err()
}

func (c *favoriteCache) MarkDirtyCommentFavoriteCount(ctx context.Context, commentID string) error {
	return c.rdb.SAdd(ctx, commentDirtyFavoriteCountKey, commentID).Err()
}

func (c *favoriteCache) ListDirtyCommentFavoriteCount(ctx context.Context) ([]string, error) {
	return c.rdb.SMembers(ctx, commentDirtyFavoriteCountKey).Result()
}

func (c *favoriteCache) RemoveDirtyCommentFavoriteCountBatch(ctx context.Context, commentIDs []interface{}) error {
	if len(commentIDs) == 0 {
		return nil
	}
	return c.rdb.SRem(ctx, commentDirtyFavoriteCountKey, commentIDs...).Err()
}

func (c *favoriteCache) TryCommentFavoriteWarmupLock(ctx context.Context, commentID string, ttl time.Duration) bool {
	ok, err := c.rdb.SetNX(ctx, buildCommentFavoriteWarmupLockKey(commentID), 1, ttl).Result()
	if err != nil {
		return false
	}
	return ok
}
