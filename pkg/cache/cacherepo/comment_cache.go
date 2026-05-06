package cacherepo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type CommentCache interface {
	ExistsCount(ctx context.Context, vid string) (bool, error)
	SetCount(ctx context.Context, vid string, count uint, ttl time.Duration) error
	GetCount(ctx context.Context, vid string) (int64, error)
	GetCountUint64(ctx context.Context, vid string) (uint64, error)
	IncrCount(ctx context.Context, vid string) error
	DecrCount(ctx context.Context, vid string) error
	DelCount(ctx context.Context, vid string) error
	TryWarmupLock(ctx context.Context, vid string, ttl time.Duration) bool
	GetList(ctx context.Context, vid string, page int, pageSize int) (string, error)
	SetList(ctx context.Context, vid string, page int, pageSize int, payload []byte, ttl time.Duration) error
	InvalidateListByVideo(ctx context.Context, vid string) error
}

type commentCache struct {
	rdb *redis.Client
}

func NewCommentCache(rdb *redis.Client) CommentCache {
	return &commentCache{rdb: rdb}
}

func buildCommentCountKey(vid string) string {
	return "comment_count:video:" + vid
}

func buildCommentWarmupLockKey(vid string) string {
	return fmt.Sprintf("warmup_lock:comment_count:%s", vid)
}

func buildCommentListCacheKey(vid string, page int, pageSize int) string {
	return "comment_list:video:" + vid + ":p" + strconv.Itoa(page) + ":s" + strconv.Itoa(pageSize)
}

func (c *commentCache) ExistsCount(ctx context.Context, vid string) (bool, error) {
	exists, err := c.rdb.Exists(ctx, buildCommentCountKey(vid)).Result()
	return exists > 0, err
}

func (c *commentCache) SetCount(ctx context.Context, vid string, count uint, ttl time.Duration) error {
	return c.rdb.Set(ctx, buildCommentCountKey(vid), count, ttl).Err()
}

func (c *commentCache) GetCount(ctx context.Context, vid string) (int64, error) {
	return c.rdb.Get(ctx, buildCommentCountKey(vid)).Int64()
}

func (c *commentCache) GetCountUint64(ctx context.Context, vid string) (uint64, error) {
	return c.rdb.Get(ctx, buildCommentCountKey(vid)).Uint64()
}

func (c *commentCache) IncrCount(ctx context.Context, vid string) error {
	return c.rdb.Incr(ctx, buildCommentCountKey(vid)).Err()
}

func (c *commentCache) DecrCount(ctx context.Context, vid string) error {
	return c.rdb.Decr(ctx, buildCommentCountKey(vid)).Err()
}

func (c *commentCache) DelCount(ctx context.Context, vid string) error {
	return c.rdb.Del(ctx, buildCommentCountKey(vid)).Err()
}

func (c *commentCache) TryWarmupLock(ctx context.Context, vid string, ttl time.Duration) bool {
	ok, err := c.rdb.SetNX(ctx, buildCommentWarmupLockKey(vid), 1, ttl).Result()
	if err != nil {
		return false
	}
	return ok
}

func (c *commentCache) GetList(ctx context.Context, vid string, page int, pageSize int) (string, error) {
	return c.rdb.Get(ctx, buildCommentListCacheKey(vid, page, pageSize)).Result()
}

func (c *commentCache) SetList(ctx context.Context, vid string, page int, pageSize int, payload []byte, ttl time.Duration) error {
	return c.rdb.Set(ctx, buildCommentListCacheKey(vid, page, pageSize), payload, ttl).Err()
}

func (c *commentCache) InvalidateListByVideo(ctx context.Context, vid string) error {
	match := fmt.Sprintf("comment_list:video:%s:*", vid)
	iter := c.rdb.Scan(ctx, 0, match, 100).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if len(keys) > 0 {
		return c.rdb.Del(ctx, keys...).Err()
	}
	return nil
}
