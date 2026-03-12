package service

import (
	"strconv"

	"github.com/jhw66/myvideo_lab4/cache"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
)

type Comment struct {
	Uid     uint   `form:"uid" json:"uid"`
	Vid     uint   `form:"vid" json:"vid"`
	Content string `form:"content" json:"content" binding:"required,max=50"`
	Cid     uint   `form:"cid" json:"cid"`
}

type CommentList struct {
	Uid      uint `form:"uid" json:"uid"`
	Vid      uint `form:"vid" json:"vid"`
	Page     int  `form:"page"`
	PageSize int  `form:"page_size"`
}

func (com Comment) AddComment() (*model.Comment, *serializer.Response) {
	redisKey := "comment_count:video:" + strconv.Itoa(int(com.Vid))

	comment := &model.Comment{
		UserID:  com.Uid,
		VideoID: com.Vid,
		Content: com.Content,
	}

	if err := com.CacheWarmUp(redisKey); err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "缓存预热失败",
		}
	}

	if err := cache.Rdb.Incr(cache.Ctx, redisKey).Err(); err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "评论保存失败(redis)",
		}
	}

	if err := model.Db.Create(&comment).Error; err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "评论保存失败(mysql)",
		}
	}
	return comment, nil
}

func (com *CommentList) CommentList() (*[]model.Comment, int64, *serializer.Response) {
	var commentList []model.Comment
	var video model.Video

	page := com.Page
	PageSize := com.PageSize
	if page <= 0 {
		page = 1
		com.Page = 1
	}
	if PageSize <= 0 {
		PageSize = 10
		com.PageSize = 10
	}
	offset := (page - 1) * PageSize

	if err := model.Db.Preload("User").Where("video_id = ?", com.Vid).
		Order("created_at desc").Limit(PageSize).Offset(offset).
		Find(&commentList).Error; err != nil {
		return nil, 0, &serializer.Response{
			Status: 500,
			Msg:    "查询评论失败",
		}
	}

	if err := model.Db.Model(&model.Video{}).Where("id = ?", com.Vid).
		Find(&video).Error; err != nil {
		return nil, 0, &serializer.Response{
			Status: 500,
			Msg:    "查询评论失败",
		}
	}

	return &commentList, int64(video.CommentCount), nil
}

func (com Comment) DelComment() (*model.Comment, *serializer.Response) {
	var comment model.Comment

	if err := model.Db.Model(&comment).Where("id = ?", com.Cid).Take(&comment).Error; err != nil {
		return nil, &serializer.Response{
			Status: 404,
			Msg:    "该评论不存在",
		}
	}

	redisKey := "comment_count:video:" + strconv.Itoa(int(comment.VideoID))
	if err := com.CacheWarmUp(redisKey); err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "缓存预热失败",
		}
	}

	if comment.UserID != com.Uid {
		return nil, &serializer.Response{
			Status: 403,
			Msg:    "不能删除别人评论",
		}
	}

	if err := cache.Rdb.Decr(cache.Ctx, redisKey).Err(); err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "删除评论失败(redis)",
		}
	}

	if err := model.Db.Delete(&comment).Error; err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "删除评论失败(mysql)",
		}
	}

	return &comment, nil
}

func (com Comment) CacheWarmUp(redisKey string) error {
	if exists, _ := cache.Rdb.Exists(cache.Ctx, redisKey).Result(); exists == 0 {
		var video model.Video
		model.Db.Where("id = ?", com.Vid).Find(&video)
		if video.FavoriteCount == 0 {
			return nil
		}
		counts := int64(video.CommentCount)
		if err := cache.Rdb.IncrBy(cache.Ctx, redisKey, counts).Err(); err != nil {
			return err
		}
	}
	return nil
}
