package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jhw66/myvideo_lab4/api"
	sessions "github.com/jhw66/myvideo_lab4/middleware"
)

func NewRouter(r *gin.Engine) {
	r.Use(sessions.Session("secret"))
	r.Use(sessions.CurrentAccount())

	v1 := r.Group("/api/v1")
	{
		//用户注册
		v1.POST("/user/register", api.UserRegister)

		//用户登录
		v1.POST("/user/login", api.UserLogin)

		//搜索视频
		v1.GET("/videos/search", api.VideoSearch)

		//排行榜
		v1.GET("/video/rank", api.RankVideos)

		//保护接口
		p := v1.Group("/p")
		p.Use(sessions.AuthLogin())
		{
			//用户详情
			p.GET("/user/me", api.UserMe)
			//登出
			p.POST("/user/logout", api.UserLogout)
			//上传头像
			p.POST("/user/avatar", api.UserAvatar)

			//上传视频
			p.POST("/video", api.UploadVideo)
			//查看视频
			p.GET("/video/me", api.MyVideo)
			//更新视频
			p.PUT("video/:id", api.UpdateVideo)
			//删除视频
			p.DELETE("video/:id", api.DeleteVideo)

			//点赞操作
			p.POST("/favorite/:vid", api.Favorite)
			//用户点赞列表
			p.GET("/favoriteList/me", api.FavoriteList)

			//评论操作
			p.POST("/comment/:vid", api.Comment)
			//视频评论列表
			p.GET("/comment/:vid", api.CommentList)
			//删除评论
			p.DELETE("/comment/:cid", api.DelComment)
		}

	}
}
