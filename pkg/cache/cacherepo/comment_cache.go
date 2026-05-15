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
	SetCount(ctx context.Context, vid string, count int64, ttl time.Duration) error
	GetCount(ctx context.Context, vid string) (int64, error)
	GetCountUint64(ctx context.Context, vid string) (uint64, error)
	IncrCount(ctx context.Context, vid string) error
	DecrCount(ctx context.Context, vid string) error
	IncrCountBy(ctx context.Context, vid string, num int64) error
	DecrCountBy(ctx context.Context, vid string, num int64) error
	DelCount(ctx context.Context, vid string) error
	ExistsRootCount(ctx context.Context, vid string) (bool, error)
	SetRootCount(ctx context.Context, vid string, count int64, ttl time.Duration) error
	GetRootCount(ctx context.Context, vid string) (int64, error)
	IncrRootCount(ctx context.Context, vid string) error
	DecrRootCount(ctx context.Context, vid string) error
	DelRootCount(ctx context.Context, vid string) error
	ExistsReplyCount(ctx context.Context, rootID string) (bool, error)
	SetReplyCount(ctx context.Context, rootID string, count int64, ttl time.Duration) error
	GetReplyCount(ctx context.Context, rootID string) (int64, error)
	IncrReplyCount(ctx context.Context, rootID string) error
	DecrReplyCount(ctx context.Context, rootID string) error
	DelReplyCount(ctx context.Context, rootID string) error
	TryWarmupLock(ctx context.Context, key string, ttl time.Duration) bool
	GetList(ctx context.Context, vid string, commentID string, page int, pageSize int) (string, error)
	SetList(ctx context.Context, vid string, commentID string, page int, pageSize int, payload []byte, ttl time.Duration) error
	InvalidateRootListByVideo(ctx context.Context, vid string) error
	InvalidateReplyListByRoot(ctx context.Context, vid string, rootID string) error
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

func buildCommentRootCountKey(vid string) string {
	return "comment_root_count:video:" + vid
}

func buildCommentReplyCountKey(rootID string) string {
	return "comment_reply_count:root:" + rootID
}

func buildCommentWarmupLockKey(key string) string {
	return "warmup_lock:" + key
}

func buildCommentListCacheKey(vid string, commentID string, page int, pageSize int) string {
	if commentID == "" {
		commentID = "root"
	}
	return "comment_list:video:" + vid + ":comment:" + commentID + ":p" + strconv.Itoa(page) + ":s" + strconv.Itoa(pageSize)
}

func (c *commentCache) ExistsCount(ctx context.Context, vid string) (bool, error) {
	exists, err := c.rdb.Exists(ctx, buildCommentCountKey(vid)).Result()
	return exists > 0, err
}

func (c *commentCache) SetCount(ctx context.Context, vid string, count int64, ttl time.Duration) error {
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

func (c *commentCache) IncrCountBy(ctx context.Context, vid string, num int64) error {
	return c.rdb.IncrBy(ctx, buildCommentCountKey(vid), num).Err()
}

func (c *commentCache) DecrCountBy(ctx context.Context, vid string, num int64) error {
	return c.rdb.DecrBy(ctx, buildCommentCountKey(vid), num).Err()
}

func (c *commentCache) DelCount(ctx context.Context, vid string) error {
	return c.rdb.Del(ctx, buildCommentCountKey(vid)).Err()
}

func (c *commentCache) ExistsRootCount(ctx context.Context, vid string) (bool, error) {
	exists, err := c.rdb.Exists(ctx, buildCommentRootCountKey(vid)).Result()
	return exists > 0, err
}

func (c *commentCache) SetRootCount(ctx context.Context, vid string, count int64, ttl time.Duration) error {
	return c.rdb.Set(ctx, buildCommentRootCountKey(vid), count, ttl).Err()
}

func (c *commentCache) GetRootCount(ctx context.Context, vid string) (int64, error) {
	return c.rdb.Get(ctx, buildCommentRootCountKey(vid)).Int64()
}

func (c *commentCache) IncrRootCount(ctx context.Context, vid string) error {
	return c.rdb.Incr(ctx, buildCommentRootCountKey(vid)).Err()
}

func (c *commentCache) DecrRootCount(ctx context.Context, vid string) error {
	return c.rdb.Decr(ctx, buildCommentRootCountKey(vid)).Err()
}

func (c *commentCache) DelRootCount(ctx context.Context, vid string) error {
	return c.rdb.Del(ctx, buildCommentRootCountKey(vid)).Err()
}

func (c *commentCache) ExistsReplyCount(ctx context.Context, rootID string) (bool, error) {
	exists, err := c.rdb.Exists(ctx, buildCommentReplyCountKey(rootID)).Result()
	return exists > 0, err
}

func (c *commentCache) SetReplyCount(ctx context.Context, rootID string, count int64, ttl time.Duration) error {
	return c.rdb.Set(ctx, buildCommentReplyCountKey(rootID), count, ttl).Err()
}

func (c *commentCache) GetReplyCount(ctx context.Context, rootID string) (int64, error) {
	return c.rdb.Get(ctx, buildCommentReplyCountKey(rootID)).Int64()
}

func (c *commentCache) IncrReplyCount(ctx context.Context, rootID string) error {
	return c.rdb.Incr(ctx, buildCommentReplyCountKey(rootID)).Err()
}

func (c *commentCache) DecrReplyCount(ctx context.Context, rootID string) error {
	return c.rdb.Decr(ctx, buildCommentReplyCountKey(rootID)).Err()
}

func (c *commentCache) DelReplyCount(ctx context.Context, rootID string) error {
	return c.rdb.Del(ctx, buildCommentReplyCountKey(rootID)).Err()
}

func (c *commentCache) TryWarmupLock(ctx context.Context, key string, ttl time.Duration) bool {
	ok, err := c.rdb.SetNX(ctx, buildCommentWarmupLockKey(key), 1, ttl).Result()
	if err != nil {
		return false
	}
	return ok
}

func (c *commentCache) GetList(ctx context.Context, vid string, commentID string, page int, pageSize int) (string, error) {
	return c.rdb.Get(ctx, buildCommentListCacheKey(vid, commentID, page, pageSize)).Result()
}

func (c *commentCache) SetList(ctx context.Context, vid string, commentID string, page int, pageSize int, payload []byte, ttl time.Duration) error {
	return c.rdb.Set(ctx, buildCommentListCacheKey(vid, commentID, page, pageSize), payload, ttl).Err()
}

func (c *commentCache) InvalidateRootListByVideo(ctx context.Context, vid string) error {
	return c.invalidateListByPattern(ctx, fmt.Sprintf("comment_list:video:%s:comment:root:*", vid))
}

func (c *commentCache) InvalidateReplyListByRoot(ctx context.Context, vid string, rootID string) error {
	return c.invalidateListByPattern(ctx, fmt.Sprintf("comment_list:video:%s:comment:%s:*", vid, rootID))
}

func (c *commentCache) InvalidateListByVideo(ctx context.Context, vid string) error {
	return c.invalidateListByPattern(ctx, fmt.Sprintf("comment_list:video:%s:*", vid))
}

func (c *commentCache) invalidateListByPattern(ctx context.Context, match string) error {
	iter := c.rdb.Scan(ctx, 0, match, 100).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	if len(keys) > 0 {
		return c.rdb.Del(ctx, keys...).Err()
	}
	return nil
}
