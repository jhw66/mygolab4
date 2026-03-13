package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
	"github.com/jhw66/myvideo_lab4/utils"
)

// 这里要带空格！
//const prefix = "Bearer "

func AccessAuthProtect() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if tokenString == "" || err != nil {
			c.JSON(401, serializer.Response{
				Status: 401,
				Msg:    "无access令牌",
			})
			c.Abort()
			return
		}

		// if len(tokenHeader) <= len(prefix) || tokenHeader[:len(prefix)] != prefix {
		// 	c.JSON(401, serializer.Response{
		// 		Status: 401,
		// 		Msg:    "令牌不合法",
		// 	})
		// 	c.Abort()
		// 	return
		// }

		//tokenString := tokenHeader[len(prefix):]
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(401, serializer.Response{
				Status: 401,
				Msg:    "access令牌过期或者不合法",
			})
			c.Abort()
			return
		}

		if claims.TokenType != "access" {
			c.JSON(401, serializer.Response{
				Status: 401,
				Msg:    "access令牌类型错误",
			})
			c.Abort()
			return
		}

		user, err := model.GetUserByID(claims.UserID)
		if err != nil || user == nil {
			c.JSON(401, serializer.Response{
				Status: 401,
				Msg:    "用户不存在或者access令牌无效",
			})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
