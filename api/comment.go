package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
	"github.com/jhw66/myvideo_lab4/service"
)

func Comment(c *gin.Context) {
	var comment service.Comment
	if err := c.ShouldBind(&comment); err != nil {
		c.JSON(400, serializer.Response{
			Status: 400,
			Msg:    "评论不能为空",
		})
		return
	}
	userValue, _ := c.Get("user")
	user := userValue.(*model.User)

	vid := c.Param("vid")
	_, err := service.FindVideoByVid(vid)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	comment.Uid = user.ID
	comment.Vid = vid

	com, err := comment.AddComment()
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(200, serializer.BuildComment(com))
}

func CommentList(c *gin.Context) {
	var comment service.CommentList

	if err := c.ShouldBindQuery(&comment); err != nil {
		c.JSON(400, serializer.Response{
			Status: 400,
			Msg:    "参数错误",
		})
		return
	}

	vid := c.Param("vid")
	comment.Vid = vid
	_, err := service.FindVideoByVid(vid)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	comments, total, err := comment.CommentList()
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(200, serializer.BuildCommentListResponse(comments, total, comment.Page, comment.PageSize))
}

func DelComment(c *gin.Context) {
	var comment service.Comment

	userValue, _ := c.Get("user")
	user := userValue.(*model.User)

	comment.Uid = user.ID
	comment.Cid = c.Param("cid")

	_, err := comment.DelComment()
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(200, serializer.Response{
		Status: 200,
		Msg:    "删除评论成功",
	})

}
