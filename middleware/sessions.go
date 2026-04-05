package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
)

func Session(secret string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 3600, Path: "/"})
	return sessions.Sessions("my_cookie", store)
}

func CurrentAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		rawID := session.Get("user_id")
		if rawID != nil {
			userID, ok := rawID.(string)
			if ok {
				user, err := model.GetUserByID(userID)
				if err == nil && user != nil {
					c.Set("user", user)
				}
			}
		}
		c.Next()
	}
}

func AuthLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(401, serializer.Response{
				Status: 401,
				Msg:    "需要登录",
			})
			c.Abort()
			return
		}

		if _, ok := user.(*model.User); !ok {
			c.JSON(401, serializer.Response{
				Status: 401,
				Msg:    "需要登录",
			})
			c.Abort()
			return
		}
		c.Next()
	}

}
