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
		v1.POST("user/register", api.UserRegister)

		//用户登录
		v1.POST("user/login", api.UserLogin)

		//保护接口
		p := v1.Group("/")
		p.Use(sessions.AuthLogin())
		{
			p.GET("user/me", api.UserMe)
			p.POST("user/logout", api.UserLogout)
		}
	}
}
