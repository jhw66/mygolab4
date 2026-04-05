package service

import (
	"errors"
	"fmt"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
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

func buildFavoriteCountKey(vid string) string {
	return "favorite_count:video:" + vid
}

func (favorite Favorite) Favorite() *serializer.Response {
	if err := favorite.WarmUpFavoriteCount(); err != nil {
		return &serializer.Response{Status: 500, Msg: "缓存预热失败"}
	}

	liked, res := favorite.toggleByMySQL()
	if res != nil {
		return res
	}

	UpdateRankScore(favorite.Vid)

	msg := "取消点赞成功"
	if liked {
		msg = "点赞成功"
	}
	return &serializer.Response{Status: 200, Msg: msg}
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

func (favorite Favorite) WarmUpFavoriteCount() error {
	countKey := buildFavoriteCountKey(favorite.Vid)

	if exists, _ := cache.Rdb.Exists(cache.Ctx, countKey).Result(); exists > 0 {
		return nil
	}

	lockKey := fmt.Sprintf("warmup_lock:favorite:%s", favorite.Vid)
	if !cache.TryWarmupLock(lockKey, 5*time.Second) {
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	if exists, _ := cache.Rdb.Exists(cache.Ctx, countKey).Result(); exists > 0 {
		return nil
	}

	var video model.Video
	model.Db.Where("id = ?", favorite.Vid).Take(&video)
	return cache.Rdb.Set(cache.Ctx, countKey, video.FavoriteCount, favoriteCountTTL).Err()
}

func (favorite Favorite) toggleByMySQL() (bool, *serializer.Response) {
	liked := false

	err := model.Db.Transaction(func(tx *gorm.DB) error {
		var record model.Favorite
		err := tx.Where("user_id = ? AND video_id = ?", favorite.Uid, favorite.Vid).Take(&record).Error
		if err == nil {
			result := tx.Where("user_id = ? AND video_id = ?", favorite.Uid, favorite.Vid).Delete(&model.Favorite{})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				liked = false
				return nil
			}
			liked = false
			return changeFavoriteCount(favorite.Vid, -1)
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		createResult := tx.Create(&model.Favorite{
			UserID:  favorite.Uid,
			VideoID: favorite.Vid,
		})
		if createResult.Error != nil {
			if isDuplicateFavoriteError(createResult.Error) {
				liked = true
				return nil
			}
			return createResult.Error
		}

		liked = true
		return changeFavoriteCount(favorite.Vid, 1)
	})
	if err != nil {
		return false, &serializer.Response{Status: 500, Msg: "点赞操作失败"}
	}

	return liked, nil
}

func changeFavoriteCount(vid string, delta int64) error {
	countKey := buildFavoriteCountKey(vid)
	pipe := cache.Rdb.Pipeline()
	pipe.IncrBy(cache.Ctx, countKey, delta)
	pipe.Expire(cache.Ctx, countKey, favoriteCountTTL)
	_, err := pipe.Exec(cache.Ctx)
	return err
}

func isDuplicateFavoriteError(err error) bool {
	var mysqlErr *mysqlDriver.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
