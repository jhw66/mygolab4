package service

import (
	"strconv"
	"strings"
	"time"

	"github.com/jhw66/myvideo_lab4/cache"
	"github.com/jhw66/myvideo_lab4/model"
)

func SyncFavoirte() {
	for {
		var redisChangeKey = "favorite:change_videos:"
		var favoritesAdd []model.Favorite
		var favoriteDel []model.Favorite
		for i := 0; i < 10; i++ {
			res, _ := cache.Rdb.RPop(cache.Ctx, redisChangeKey).Result()
			if res == "" {
				time.Sleep(1 * time.Second)
				continue
			}

			arr := strings.Split(res, ":")
			action := arr[0]
			uid, _ := strconv.Atoi(arr[1])
			vid, _ := strconv.Atoi(arr[2])

			if action == "add" {
				favoritesAdd = append(favoritesAdd, model.Favorite{
					UserID:  uint(uid),
					VideoID: uint(vid),
				})
			}
			if action == "delete" {
				favoriteDel = append(favoriteDel, model.Favorite{
					UserID:  uint(uid),
					VideoID: uint(vid),
				})
			}
		}
		//点击过快会导致更新失败
		//2026/03/11 17:33:21 D:/myvideo_lab4/service/favorite.go:51 Error 1062 (23000): Duplicate entry '9-2' for key 'favorite.PRIMARY'
		//[0.520ms] [rows:0] INSERT INTO `favorite` (`user_id`,`video_id`) VALUES (9,2),(9,2),(9,2),(9,2),(9,2)
		if len(favoritesAdd) > 0 {
			model.Db.Create(&favoritesAdd)
		}
		if len(favoriteDel) > 0 {
			model.Db.Delete(&favoriteDel)
		}
	}
}

func SyncFavoriteCount() {
	keys, err := cache.Rdb.Keys(cache.Ctx, "favorite_count:video:*").Result()
	if err != nil {
		panic(err)
	}
	for _, key := range keys {
		parts := strings.Split(key, ":")
		favoriteCount, _ := cache.Rdb.Get(cache.Ctx, key).Int()
		vid, _ := strconv.Atoi(parts[2])

		model.Db.Model(&model.Video{}).Where("id = ?", vid).Update("favorite_count", favoriteCount)
	}
}

func SyncCommentCount() {
	keys, err := cache.Rdb.Keys(cache.Ctx, "comment_count:video:*").Result()
	if err != nil {
		panic(err)
	}
	for _, key := range keys {
		parts := strings.Split(key, ":")
		commentCount, _ := cache.Rdb.Get(cache.Ctx, key).Int()
		vid, _ := strconv.Atoi(parts[2])

		model.Db.Model(&model.Video{}).Where("id = ?", vid).Update("comment_count", commentCount)
	}
}
