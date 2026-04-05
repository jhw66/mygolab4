package api

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jhw66/myvideo_lab4/cache"
	"github.com/jhw66/myvideo_lab4/serializer"
	"github.com/jhw66/myvideo_lab4/service"
)

const rankResultCacheTTL = 10 * time.Second

func RankVideos(c *gin.Context) {
	cacheKey := "rank:result:hot:" + strconv.Itoa(service.DefaultRankLimit)

	if cached, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result(); err == nil {
		var res []serializer.Video
		if err := json.Unmarshal([]byte(cached), &res); err == nil {
			c.JSON(200, serializer.Response{
				Status: 200,
				Msg:    "获取热门排行榜成功",
				Data:   res,
			})
			return
		}
	}

	videos, err := service.GetRankVideos(service.DefaultRankLimit)
	if err != nil {
		c.JSON(500, serializer.Response{
			Status: 500,
			Msg:    "获取热门排行榜失败",
		})
		return
	}
	res := serializer.BuildVideoList(&videos)

	if data, err := json.Marshal(res); err == nil {
		cache.Rdb.Set(cache.Ctx, cacheKey, data, rankResultCacheTTL)
	}

	c.JSON(200, serializer.Response{
		Status: 200,
		Msg:    "获取热门排行榜成功",
		Data:   res,
	})
}

func VideoSearch(c *gin.Context) {
	var videosearch service.VideoSearch
	if err := c.ShouldBind(&videosearch); err != nil {
		c.JSON(404, serializer.Response{
			Status: 404,
			Msg:    "搜索词不能为空",
		})
		return
	}
	videos, err := videosearch.FindVideosByKeyword()
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(200, serializer.BuildVideoListResponse(videos))
}
