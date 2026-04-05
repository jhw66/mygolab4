package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
	"github.com/jhw66/myvideo_lab4/service"
)

func Favorite(c *gin.Context) {
	var favorite service.Favorite

	userValue, _ := c.Get("user")
	user := userValue.(*model.User)

	vid := c.Param("vid")

	_, err := service.FindVideoByVid(vid)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	favorite.Uid = user.ID
	favorite.Vid = vid

	res := favorite.Favorite()
	c.JSON(res.Status, res)
}

func FavoriteList(c *gin.Context) {
	var favorite service.Favorite

	userValue, _ := c.Get("user")
	user := userValue.(*model.User)

	favorite.Uid = user.ID

	videos, err := favorite.GetUserFavorite()
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(200, serializer.BuildVideoListResponse(videos))
}
