package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jhw66/myvideo_lab4/cache"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
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

func (com Comment) AddComment() (*model.Comment, *serializer.Response) {
	comment := &model.Comment{
		UserID:  com.Uid,
		VideoID: com.Vid,
		Content: com.Content,
	}

	if err := model.Db.Create(comment).Error; err != nil {
		return nil, &serializer.Response{Status: 500, Msg: "评论保存失败"}
	}
	model.Db.Preload("User").Take(comment)

	redisKey := buildCommentCountKey(com.Vid)
	warmUpCommentCount(com.Vid, redisKey)
	cache.Rdb.Incr(cache.Ctx, redisKey)

	UpdateRankScore(com.Vid)
	invalidateCommentListCache(com.Vid)

	return comment, nil
}

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

	cacheKey := buildCommentListCacheKey(com.Vid, page, pageSize)
	if cached, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result(); err == nil {
		var data commentListCache
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			return &data.Comments, data.Total, nil
		}
	}

	offset := (page - 1) * pageSize
	var commentList []model.Comment
	if err := model.Db.Preload("User").Where("video_id = ?", com.Vid).
		Order("created_at desc").Limit(pageSize).Offset(offset).
		Find(&commentList).Error; err != nil {
		return nil, 0, &serializer.Response{Status: 500, Msg: "查询评论失败"}
	}

	total := getCommentCount(com.Vid)

	cacheData := commentListCache{Comments: commentList, Total: total}
	if data, err := json.Marshal(cacheData); err == nil {
		cache.Rdb.Set(cache.Ctx, cacheKey, data, commentListCacheTTL)
	}

	return &commentList, total, nil
}

func (com Comment) DelComment() (*model.Comment, *serializer.Response) {
	var comment model.Comment
	if err := model.Db.Where("id = ?", com.Cid).Take(&comment).Error; err != nil {
		return nil, &serializer.Response{Status: 404, Msg: "该评论不存在"}
	}
	if comment.UserID != com.Uid {
		return nil, &serializer.Response{Status: 403, Msg: "不能删除别人评论"}
	}

	if err := model.Db.Delete(&comment).Error; err != nil {
		return nil, &serializer.Response{Status: 500, Msg: "删除评论失败"}
	}

	redisKey := buildCommentCountKey(comment.VideoID)
	warmUpCommentCount(comment.VideoID, redisKey)
	cache.Rdb.Decr(cache.Ctx, redisKey)

	UpdateRankScore(comment.VideoID)
	invalidateCommentListCache(comment.VideoID)

	return &comment, nil
}

func warmUpCommentCount(vid string, redisKey string) {
	if exists, _ := cache.Rdb.Exists(cache.Ctx, redisKey).Result(); exists > 0 {
		return
	}

	lockKey := fmt.Sprintf("warmup_lock:comment_count:%s", vid)
	if !cache.TryWarmupLock(lockKey, 5*time.Second) {
		time.Sleep(100 * time.Millisecond)
		return
	}

	if exists, _ := cache.Rdb.Exists(cache.Ctx, redisKey).Result(); exists > 0 {
		return
	}

	var video model.Video
	model.Db.Where("id = ?", vid).Take(&video)
	cache.Rdb.Set(cache.Ctx, redisKey, video.CommentCount, 24*time.Hour)
}

func getCommentCount(vid string) int64 {
	redisKey := buildCommentCountKey(vid)
	warmUpCommentCount(vid, redisKey)
	count, err := cache.Rdb.Get(cache.Ctx, redisKey).Int64()
	if err != nil {
		var video model.Video
		model.Db.Where("id = ?", vid).Take(&video)
		return int64(video.CommentCount)
	}
	return count
}

func invalidateCommentListCache(vid string) {
	pattern := fmt.Sprintf("comment_list:video:%s:*", vid)
	iter := cache.Rdb.Scan(cache.Ctx, 0, pattern, 100).Iterator()
	var keys []string
	for iter.Next(cache.Ctx) {
		keys = append(keys, iter.Val())
	}
	if len(keys) > 0 {
		cache.Rdb.Del(cache.Ctx, keys...)
	}
}
