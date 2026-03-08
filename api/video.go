package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jhw66/myvideo_lab4/cache"
	"github.com/jhw66/myvideo_lab4/serializer"
	"github.com/jhw66/myvideo_lab4/service"
)

func RankVideos(c *gin.Context) {
	redisKey := "video_hot_rank"

	val, err := cache.Rdb.Get(cache.Ctx, redisKey).Result()
	if err == nil {
		var res []serializer.Video
		if err := json.Unmarshal([]byte(val), &res); err == nil {
			c.JSON(http.StatusOK, serializer.Response{
				Status: 200,
				Msg:    "获取热门排行榜成功from redis",
				Data:   res,
			})
			return
		}
	}

	videos, err := service.GetRankVideos(10)
	if err != nil {
		c.JSON(500, serializer.Response{
			Status: 500,
			Msg:    "获取热门排行榜失败",
		})
		return
	}
	res := serializer.BuildVideoList(&videos)

	if data, err := json.Marshal(res); err == nil {
		cache.Rdb.Set(cache.Ctx, redisKey, data, 5*time.Minute)
	}

	c.JSON(http.StatusOK, serializer.Response{
		Status: 200,
		Msg:    "获取热门排行榜成功from mysql",
		Data:   res,
	})
}
