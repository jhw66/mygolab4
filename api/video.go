package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jhw66/myvideo_lab4/serializer"
	"github.com/jhw66/myvideo_lab4/service"
)

func RankVideos(c *gin.Context) {
	videos, err := service.GetRankVideos(service.DefaultRankLimit)
	if err != nil {
		c.JSON(500, serializer.Response{
			Status: 500,
			Msg:    "获取热门排行榜失败",
		})
		return
	}
	c.JSON(200, serializer.BuildVideoListResponse(videos))
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
