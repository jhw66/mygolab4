package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
	"github.com/jhw66/myvideo_lab4/service"
)

func RelationAction(c *gin.Context) {
	var relation service.Relation
	targetID := c.Param("uid")

	userValue, _ := c.Get("user")
	user := userValue.(*model.User)

	_, err := model.GetUserByID(targetID)
	if err != nil {
		c.JSON(404, serializer.Response{
			Status: 404,
			Msg:    "目标用户不存在",
		})
		return
	}

	if user.ID == targetID {
		c.JSON(400, serializer.Response{
			Status: 400,
			Msg:    "不可关注自己",
		})
		return
	}

	res := relation.RelationAction(targetID, user.ID)

	c.JSON(res.Status, res)
}

func FollowingList(c *gin.Context) {
	var relation service.Relation

	userValue, _ := c.Get("user")
	user := userValue.(*model.User)

	res := relation.FollowingList(user.ID)

	c.JSON(res.Status, res)
}

func FollowerList(c *gin.Context) {
	var relation service.Relation

	userValue, _ := c.Get("user")
	user := userValue.(*model.User)

	res := relation.FollowerList(user.ID)

	c.JSON(res.Status, res)
}

func FriendList(c *gin.Context) {
	var relation service.Relation

	userValue, _ := c.Get("user")
	user := userValue.(*model.User)

	res := relation.FriendList(user.ID)

	c.JSON(res.Status, res)
}
