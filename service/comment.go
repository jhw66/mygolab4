package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jhw66/myvideo_lab4/cache"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
	"gorm.io/gorm"
)

type Comment struct {
	Uid     string `form:"uid" json:"uid"`
	Vid     string `form:"vid" json:"vid"`
	Content string `form:"content" json:"content" binding:"required,max=50"`
	Cid     string `form:"cid" json:"cid"`
}

type CommentList struct {
	Uid      string `form:"uid" json:"uid"`
	Vid      string `form:"vid" json:"vid"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

const commentListCacheTTL = 30 * time.Second
const commentCountTTL = 24 * time.Hour

type commentListCache struct {
	Comments []model.Comment `json:"comments"`
	Total    int64           `json:"total"`
}

func buildCommentCountKey(vid string) string {
	return "comment_count:video:" + vid
}

func buildCommentListCacheKey(vid string, page, pageSize int) string {
	return fmt.Sprintf("comment_list:video:%s:p%d:s%d", vid, page, pageSize)
}

// 添加评论
func (com Comment) AddComment() *serializer.Response {
	// 预热缓存移到事务外，避免事务回滚后 Redis 不一致
	if err := warmUpCommentCount(com.Vid); err != nil {
		return &serializer.Response{Status: 500, Msg: "预热评论数缓存失败"}
	}

	err := model.Db.Transaction(func(tx *gorm.DB) error {
		var video model.Video
		if err := tx.Where("id = ?", com.Vid).Take(&video).Error; err != nil {
			return errors.New("视频不存在")
		}

		comment := &model.Comment{
			UserID:  com.Uid,
			VideoID: com.Vid,
			Content: com.Content,
		}

		if err := tx.Create(comment).Error; err != nil {
			return errors.New("评论保存失败")
		}

		return nil
	})

	if err != nil {
		return &serializer.Response{Status: 500, Msg: err.Error()}
	}

	// 更新评论数缓存和排行榜
	cache.Rdb.Incr(cache.Ctx, buildCommentCountKey(com.Vid))
	UpdateRankScore(com.Vid)
	invalidateCommentListCache(com.Vid)

	return &serializer.Response{Status: 200, Msg: "评论成功"}
}

// 删除评论
func (com Comment) DelComment() *serializer.Response {
	err := model.Db.Transaction(func(tx *gorm.DB) error {
		var comment model.Comment
		if err := tx.Where("id = ? AND video_id = ?", com.Cid, com.Vid).Take(&comment).Error; err != nil {
			return errors.New("该评论不存在")
		}

		if comment.UserID != com.Uid {
			return errors.New("不能删除别人评论")
		}

		if err := tx.Delete(&comment).Error; err != nil {
			return errors.New("删除评论失败")
		}

		return nil
	})

	if err != nil {
		return &serializer.Response{Status: 500, Msg: err.Error()}
	}

	//由于在api层没有判断视频是否存在，将这个放在事务上面会导致即使视频不存在，也会预热评论数缓存
	// 预热缓存移到事务外，避免事务回滚后 Redis 不一致
	if err := warmUpCommentCount(com.Vid); err != nil {
		return &serializer.Response{Status: 500, Msg: "预热评论数缓存失败"}
	}

	// 更新评论数缓存和排行榜
	cache.Rdb.Decr(cache.Ctx, buildCommentCountKey(com.Vid))
	UpdateRankScore(com.Vid)
	invalidateCommentListCache(com.Vid)

	return &serializer.Response{Status: 200, Msg: "删除评论成功"}
}

// 获取评论列表
func (com *CommentList) CommentList() (*[]model.Comment, int64, *serializer.Response) {
	page := com.Page
	pageSize := com.PageSize
	if page <= 0 {
		page = 1
		com.Page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
		com.PageSize = 10
	}

	// 查询评论列表和评论数 from redis
	cacheKey := buildCommentListCacheKey(com.Vid, page, pageSize)
	if cached, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result(); err == nil {
		var data commentListCache
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			return &data.Comments, data.Total, nil
		}
	}

	// 查询评论列表 from mysql
	offset := (page - 1) * pageSize
	var commentList []model.Comment
	if err := model.Db.Preload("User").Where("video_id = ?", com.Vid).
		Order("created_at desc").Limit(pageSize).Offset(offset).
		Find(&commentList).Error; err != nil {
		return nil, 0, &serializer.Response{Status: 500, Msg: "查询评论失败"}
	}

	//优先从 Redis 获取评论数（实时更新），而非 MySQL（有 30s 同步延迟）
	total := getCommentCount(com.Vid)

	cacheData := commentListCache{Comments: commentList, Total: total}
	if data, err := json.Marshal(cacheData); err == nil {
		cache.Rdb.Set(cache.Ctx, cacheKey, data, commentListCacheTTL)
	}

	return &commentList, total, nil
}

// 从 Redis 获取评论数
func getCommentCount(vid string) int64 {
	redisKey := buildCommentCountKey(vid)
	if err := warmUpCommentCount(vid); err != nil {
		return 0
	}

	count, err := cache.Rdb.Get(cache.Ctx, redisKey).Int64()
	if err != nil {
		return 0
	}
	return count
}

// 预热评论数缓存
func warmUpCommentCount(vid string) error {
	commkey := buildCommentCountKey(vid)
	if exists, _ := cache.Rdb.Exists(cache.Ctx, commkey).Result(); exists > 0 {
		return nil
	}

	// 分布式锁+双重检查，防止并发预热导致重复 DB 查询
	lockKey := fmt.Sprintf("warmup_lock:comment_count:%s", vid)
	if !cache.TryWarmupLock(lockKey, 5*time.Second) {
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	if exists, _ := cache.Rdb.Exists(cache.Ctx, commkey).Result(); exists > 0 {
		return nil
	}

	var video model.Video
	model.Db.Where("id = ?", vid).Take(&video)
	return cache.Rdb.Set(cache.Ctx, commkey, video.CommentCount, commentCountTTL).Err()
}

// 清除评论列表缓存
func invalidateCommentListCache(vid string) {
	match := fmt.Sprintf("comment_list:video:%s:*", vid)

	iter := cache.Rdb.Scan(cache.Ctx, 0, match, 100).Iterator()
	var keys []string
	for iter.Next(cache.Ctx) {
		keys = append(keys, iter.Val())
	}
	if len(keys) > 0 {
		cache.Rdb.Del(cache.Ctx, keys...)
	}
}
