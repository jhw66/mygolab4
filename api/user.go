package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
	"github.com/jhw66/myvideo_lab4/service"
	"github.com/jhw66/myvideo_lab4/utils"
)

func UserRegister(c *gin.Context) {
	var register service.UserRegister
	if err := c.ShouldBind(&register); err != nil {
		c.JSON(400, serializer.Response{
			Status: 400,
			Msg:    "输入格式错误",
		})
		return
	}

	user, res := register.Register()
	if res != nil {
		c.JSON(res.Status, res)
		return
	}

	c.JSON(200, serializer.BuildUserResponse(user))
}

func UserLogin(c *gin.Context) {
	var login service.UserLogin

	if err := c.ShouldBind(&login); err != nil {
		c.JSON(400, serializer.Response{
			Status: 400,
			Msg:    "输入格式错误",
		})
		return
	}

	if userVal, exists := c.Get("user"); exists {
		if user, ok := userVal.(*model.User); ok && user.UserName == login.UserName {
			c.JSON(409, serializer.Response{
				Status: 409,
				Msg:    "请勿重复登录",
			})
			return
		}
	}

	user, res := login.Login()
	if res != nil {
		c.JSON(res.Status, res)
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(500, serializer.Response{
			Status: 500,
			Msg:    "生成accessToken失败",
		})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(500, serializer.Response{
			Status: 500,
			Msg:    "生成freshToken失败",
		})
		return
	}

	c.SetCookieData(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/refresh",
		MaxAge:   int((24 * time.Hour).Seconds()),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	})

	c.SetCookieData(&http.Cookie{
		Name:   "access_token",
		Value:  accessToken,
		Path:   "/",
		MaxAge: int((15 * time.Minute).Seconds()),
	})

	c.JSON(200, serializer.BuildUserResponse(user))
}

func RefreshHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if refreshToken == "" || err != nil {
		c.JSON(401, serializer.Response{
			Status: 401,
			Msg:    "无refresh令牌",
		})
		return
	}

	claims, err := utils.ParseToken(refreshToken)
	if err != nil {
		c.JSON(401, serializer.Response{
			Status: 401,
			Msg:    err.Error(),
		})
		return
	}

	if claims.TokenType != "refresh" {
		c.JSON(401, serializer.Response{
			Status: 401,
			Msg:    "refresh令牌类型错误",
		})
		return
	}

	newAccessToken, err := utils.GenerateAccessToken(claims.UserID)
	if err != nil {
		c.JSON(500, serializer.Response{
			Status: 500,
			Msg:    "生成新accesstoken失败",
		})
		return
	}

	c.SetCookieData(&http.Cookie{
		Name:   "access_token",
		Value:  newAccessToken,
		Path:   "/",
		MaxAge: int((15 * time.Minute).Seconds()),
	})

	c.JSON(200, serializer.Response{
		Status: 200,
		Msg:    "重新生成accesstoken成功",
	})
}

func UserLogout(c *gin.Context) {
	c.SetCookieData(&http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1/refresh",
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	})

	c.SetCookieData(&http.Cookie{
		Name:   "access_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	c.JSON(200, serializer.Response{
		Status: 200,
		Msg:    "登出成功",
	})
}

func UserMe(c *gin.Context) {
	if user, _ := c.Get("user"); user != nil {
		if u, ok := user.(*model.User); ok {
			res := serializer.BuildUserResponse(u)
			c.JSON(200, res)
			return
		}
	}
	c.JSON(404, serializer.Response{
		Status: 404,
		Msg:    "资源不存在",
	})
}

func UserAvatar(c *gin.Context) {
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(401, serializer.Response{
			Status: 401,
			Msg:    "未登录",
		})
		return
	}

	user, _ := userValue.(*model.User)
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(400, serializer.Response{
			Status: 400,
			Msg:    "头像文件不能为空",
		})
		return
	}
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("avatar_%d_%d%s", user.ID, time.Now().Unix(), ext)

	dir := "static/avatar"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		c.JSON(500, serializer.Response{
			Status: 500,
			Msg:    "创建目录失败",
		})
		return
	}

	//删除旧头像
	if user.Avatar != "" {
		oldpath := strings.TrimPrefix(user.Avatar, "/")
		if _, err := os.Stat(oldpath); err == nil {
			os.Remove(oldpath)
		}
	}

	savepath := filepath.Join(dir, filename)
	if err := c.SaveUploadedFile(file, savepath); err != nil {
		c.JSON(500, serializer.Response{
			Status: 500,
			Msg:    "保存头像失败",
		})
		return
	}

	user.Avatar = "/" + strings.ReplaceAll(savepath, "\\", "/")
	if _, res := service.UploadAvatar(user); res != nil {
		c.JSON(res.Status, res)
	}
	c.JSON(200, serializer.BuildUserResponse(user))
}

// func UserRegister(c *gin.Context) {
// 	var useregister service.UserRegister
// 	err := c.ShouldBind(&useregister)
// 	if err == nil {
// 		if user, err := useregister.Register(); err != nil {
// 			c.JSON(err.Status, err)
// 		} else {
// 			res := serializer.BuildUserResponse(user)
// 			c.JSON(200, res)
// 		}
// 	} else {
// 		c.JSON(400, serializer.Response{
// 			Status: 400,
// 			Msg:    "输入格式错误",
// 		})
// 	}
// }

// func UserLogin(c *gin.Context) {
// 	var userlogin service.UserLogin
// 	if err := c.ShouldBind(&userlogin); err == nil {
// 		if user, err := c.Get("user"); err != false && user.(*model.User).UserName == userlogin.UserName {
// 			c.JSON(409, serializer.Response{
// 				Status: 409,
// 				Msg:    "请勿重复登录",
// 			})
// 			return
// 		}
// 		if user, err := userlogin.Login(); err != nil {
// 			c.JSON(err.Status, err)
// 		} else {
// 			s := sessions.Default(c)
// 			s.Clear()
// 			s.Set("user_id", user.ID)
// 			s.Save()
// 			res := serializer.BuildUserResponse(user)
// 			c.JSON(200, res)
// 		}
// 	} else {
// 		c.JSON(400, serializer.Response{
// 			Status: 400,
// 			Msg:    "输入格式错误",
// 		})
// 	}
// }

// func UserLogout(c *gin.Context) {
// 	s := sessions.Default(c)
// 	s.Clear()
// 	s.Save()
// 	c.JSON(200, serializer.Response{
// 		Status: 200,
// 		Msg:    "登出成功",
// 	})
// }
