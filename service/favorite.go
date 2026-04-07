package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/jhw66/myvideo_lab4/cache"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
	"gorm.io/gorm"
)

type Favorite struct {
	Uid string
	Vid string
}

const favoriteCountTTL = 24 * time.Hour

// 构建点赞数缓存key
func buildFavoriteCountKey(vid string) string {
	return "favorite_count:video:" + vid
}

func (favorite Favorite) Favorite() *serializer.Response {
	if err := warmUpFavoriteCount(favorite.Vid); err != nil {
		return &serializer.Response{
			Status: 500,
			Msg:    "缓存预热失败",
		}
	}

	// 点赞操作
	liked, err, res := favorite.DoFavoriteAction()
	if err != nil {
		return res
	}

	// 更新排行榜
	err = UpdateRankScore(favorite.Vid)
	if err != nil {
		return &serializer.Response{
			Status: 500,
			Msg:    "更新排行榜失败",
		}
	}

	msg := "取消点赞成功"
	if liked {
		msg = "点赞成功"
	}

	return &serializer.Response{
		Status: 200,
		Msg:    msg,
	}
}

func (favorite Favorite) GetUserFavorite() (*[]model.Video, *serializer.Response) {
	var videos []model.Video
	err := model.Db.
		Joins("JOIN favorite ON favorite.video_id = video.id").
		Where("favorite.user_id = ?", favorite.Uid).
		Find(&videos).Error
	if err != nil {
		return nil, &serializer.Response{Status: 500, Msg: "查询失败"}
	}
	return &videos, nil
}

// 预热点赞数缓存
func warmUpFavoriteCount(vid string) error {
	countKey := buildFavoriteCountKey(vid)
	if exists, _ := cache.Rdb.Exists(cache.Ctx, countKey).Result(); exists > 0 {
		return nil
	}

	// 分布式锁 + 双重检查，防止并发预热导致重复 DB 查询
	lockKey := fmt.Sprintf("warmup_lock:favorite:%s", vid)
	if !cache.TryWarmupLock(lockKey, 5*time.Second) {
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	if exists, _ := cache.Rdb.Exists(cache.Ctx, countKey).Result(); exists > 0 {
		return nil
	}

	var video model.Video
	model.Db.Where("id = ?", vid).Take(&video)
	return cache.Rdb.Set(cache.Ctx, countKey, video.FavoriteCount, favoriteCountTTL).Err()
}

// 点赞操作
func (favorite Favorite) DoFavoriteAction() (bool, error, *serializer.Response) {
	var liked bool
	err := model.Db.Transaction(func(tx *gorm.DB) error {
		var isexist model.Favorite
		err := tx.Where("user_id = ? AND video_id = ?", favorite.Uid, favorite.Vid).Take(&isexist).Error

		switch {
		case err == nil:
			liked = false
			return tx.Where("user_id = ? AND video_id = ?", favorite.Uid, favorite.Vid).Delete(&model.Favorite{}).Error

		case errors.Is(err, gorm.ErrRecordNotFound):
			liked = true
			return tx.Create(&model.Favorite{
				UserID:  favorite.Uid,
				VideoID: favorite.Vid,
			}).Error

		default:
			return err
		}
	})
	if err != nil {
		return false, err, &serializer.Response{Status: 500, Msg: "点赞操作失败", Error: err.Error()}
	}

	if liked {
		cache.Rdb.Incr(cache.Ctx, buildFavoriteCountKey(favorite.Vid))
	} else {
		cache.Rdb.Decr(cache.Ctx, buildFavoriteCountKey(favorite.Vid))
	}

	return liked, nil, &serializer.Response{Status: 200, Msg: "点赞操作成功"}
}
