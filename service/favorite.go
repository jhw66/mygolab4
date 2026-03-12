package service

import (
	"fmt"
	"strconv"

	"github.com/jhw66/myvideo_lab4/cache"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
)

type Favorite struct {
	Uid int
	Vid int
}

const redisChangeKey = "favorite:change_videos:"

func (favorite Favorite) Favorite() *serializer.Response {
	redisKey := "favorite:video:" + strconv.Itoa(favorite.Vid)
	redisCountKey := "favorite_count:video:" + strconv.Itoa(favorite.Vid)

	if err := favorite.CacheWarmUp(redisKey, redisCountKey); err != nil {
		return &serializer.Response{
			Status: 500,
			Msg:    "缓存预热失败",
		}
	}

	if cache.Rdb.SIsMember(cache.Ctx, redisKey, favorite.Uid).Val() {
		if res := favorite.DeleteFavorite(); res != nil {
			return &serializer.Response{
				Status: res.Status,
				Msg:    res.Msg,
			}
		}
		change := fmt.Sprintf("delete:%d:%d", favorite.Uid, favorite.Vid)
		cache.Rdb.LPush(cache.Ctx, redisChangeKey, change)
		return &serializer.Response{
			Status: 200,
			Msg:    "取消点赞成功",
		}

	}

	if _, err := cache.Rdb.SAdd(cache.Ctx, redisKey, favorite.Uid).Result(); err != nil {
		return &serializer.Response{
			Status: 500,
			Msg:    "点赞失败",
		}
	}
	change := fmt.Sprintf("add:%d:%d", favorite.Uid, favorite.Vid)
	cache.Rdb.LPush(cache.Ctx, redisChangeKey, change)
	cache.Rdb.Incr(cache.Ctx, redisCountKey)
	return &serializer.Response{
		Status: 200,
		Msg:    "点赞成功",
	}
}

func (favorite Favorite) DeleteFavorite() *serializer.Response {
	redisKey := "favorite:video:" + strconv.Itoa(favorite.Vid)
	redisCountKey := "favorite_count:video:" + strconv.Itoa(favorite.Vid)

	if _, err := cache.Rdb.SRem(cache.Ctx, redisKey, favorite.Uid).Result(); err != nil {
		return &serializer.Response{
			Status: 500,
			Msg:    "取消点赞失败",
		}
	}
	cache.Rdb.Decr(cache.Ctx, redisCountKey)

	return nil
}

func (favorite Favorite) GetFavorite() (uint, *serializer.Response) {
	redisCountKey := "favorite_count:video:" + strconv.Itoa(favorite.Vid)
	count, err := strconv.Atoi(cache.Rdb.Get(cache.Ctx, redisCountKey).Val())
	if err != nil {
		return 0, &serializer.Response{
			Status: 500,
			Msg:    "获取点赞数失败",
		}
	}
	model.Db.Model(&model.Video{}).Where("id = ?", favorite.Vid).Update("favorite_count", uint(count))
	return uint(count), nil
}

func (favorite Favorite) GetUserFavorite() (*[]model.Video, *serializer.Response) {
	var videos []model.Video

	err := model.Db.
		Joins("JOIN favorite ON favorite.video_id = video.id").
		Where("favorite.user_id = ?", favorite.Uid).
		Find(&videos).Error

	if err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "查询失败",
		}
	}
	return &videos, nil
}

// 缓存预热
func (favorite Favorite) CacheWarmUp(redisKey string, redisCountKey string) error {
	if exists, _ := cache.Rdb.Exists(cache.Ctx, redisKey).Result(); exists == 0 {
		var favorites []model.Favorite
		var counts int64
		model.Db.Where("video_id = ?", favorite.Vid).Find(&favorites).Count(&counts)
		if counts == 0 {
			return nil
		}
		for _, fav := range favorites {
			if err := cache.Rdb.SAdd(cache.Ctx, redisKey, fav.UserID).Err(); err != nil {
				return err
			}
		}
	}

	if exists, _ := cache.Rdb.Exists(cache.Ctx, redisCountKey).Result(); exists == 0 {
		var video model.Video
		model.Db.Where("id = ?", favorite.Vid).Find(&video)
		if video.FavoriteCount == 0 {
			return nil
		}
		counts := int64(video.FavoriteCount)
		if err := cache.Rdb.IncrBy(cache.Ctx, redisCountKey, counts).Err(); err != nil {
			return err
		}
	}
	return nil
}
