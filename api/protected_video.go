package api

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
	"github.com/jhw66/myvideo_lab4/service"
	"gorm.io/gorm"
)

func UploadVideo(c *gin.Context) {
	model.Db.Transaction(func(tx *gorm.DB) error {
		userValue, _ := c.Get("user")
		user := userValue.(*model.User)

		title := c.PostForm("title")
		info := c.PostForm("info")
		if title == "" {
			c.JSON(400, serializer.Response{
				Status: 400,
				Msg:    "标题不能为空",
			})
			return errors.New("标题不能为空")
		}

		// 接收视频文件
		videoFile, err := c.FormFile("video")
		if err != nil {
			c.JSON(400, serializer.Response{
				Status: 400,
				Msg:    "视频文件不能为空",
			})
			return err
		}

		// 接收封面文件
		coverFile, err := c.FormFile("cover")
		if err != nil {
			c.JSON(400, serializer.Response{
				Status: 400,
				Msg:    "封面文件不能为空",
			})
			return err
		}

		videoDir := "static/video"
		coverDir := "static/cover"
		if err := os.MkdirAll(videoDir, os.ModePerm); err != nil {
			c.JSON(500, serializer.Response{
				Status: 500,
				Msg:    "创建视频文件失败",
			})
			return err
		}
		if err := os.MkdirAll(coverDir, os.ModePerm); err != nil {
			c.JSON(500, serializer.Response{
				Status: 500,
				Msg:    "创建封面文件失败",
			})
			return err
		}

		videoName := fmt.Sprintf("video_%s_%d%s", user.ID, time.Now().Unix(), filepath.Ext(videoFile.Filename))
		videoPath := filepath.Join(videoDir, videoName)
		coverName := fmt.Sprintf("cover_%s_%d%s", user.ID, time.Now().Unix(), filepath.Ext(coverFile.Filename))
		coverPath := filepath.Join(coverDir, coverName)
		c.SaveUploadedFile(videoFile, videoPath)
		c.SaveUploadedFile(coverFile, coverPath)

		video := model.Video{
			UserID:        user.ID,
			Title:         title,
			Info:          info,
			URL:           "/" + videoPath,
			Cover:         "/" + coverPath,
			FavoriteCount: 0,
			CommentCount:  0,
			HotScore:      service.CalculateHotScore(0, 0),
		}

		if _, err := service.UploadVideo(tx, &video); err != nil {
			c.JSON(err.Status, err)
			return errors.New(err.Msg)
		}

		c.JSON(200, serializer.BuildVideoResponse(&video))
		return nil
	})
}

func MyVideo(c *gin.Context) {
	userValue, _ := c.Get("user")
	user := userValue.(*model.User)

	videos, err := service.FindVideoByUser(user)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(200, serializer.BuildVideoListResponse(videos))
}

func UpdateVideo(c *gin.Context) {
	model.Db.Transaction(func(tx *gorm.DB) error {
		userValue, _ := c.Get("user")
		user := userValue.(*model.User)

		vid := c.Param("id")
		if vid == "" {
			c.JSON(400, serializer.Response{
				Status: 400,
				Msg:    "请传入视频id",
			})
			return errors.New("请传入视频id")
		}

		if !service.CompareVidAndUid(user.ID, vid) {
			c.JSON(403, serializer.Response{
				Status: 403,
				Msg:    "没有修改视频权限或者不存在该视频",
			})
			return errors.New("没有修改视频权限或者不存在该视频")
		}

		video, res := service.FindVideoByVid(vid)
		if res != nil {
			c.JSON(res.Status, res)
			return errors.New(res.Msg)
		}

		title := c.PostForm("title")
		info := c.PostForm("info")
		if title == "" {
			c.JSON(400, serializer.Response{
				Status: 400,
				Msg:    "标题不能为空",
			})
			return errors.New("标题不能为空")
		}
		video.Title = title
		video.Info = info

		if videoFile, err := c.FormFile("video"); err == nil {
			videoDir := "static/video"
			if err := os.MkdirAll(videoDir, os.ModePerm); err != nil {
				c.JSON(500, serializer.Response{
					Status: 500,
					Msg:    "创建视频目录失败",
				})
				return err
			}
			if video.URL != "" {
				oldPath := strings.TrimPrefix(video.URL, "/")
				if _, err := os.Stat(oldPath); err == nil {
					os.Remove(oldPath)
				}
			}
			videoName := fmt.Sprintf("video_%s_%d%s", user.ID, time.Now().Unix(), filepath.Ext(videoFile.Filename))
			videoPath := filepath.Join(videoDir, videoName)
			if err := c.SaveUploadedFile(videoFile, videoPath); err != nil {
				c.JSON(500, serializer.Response{
					Status: 500,
					Msg:    "保存视频文件失败",
				})
				return err
			}
			video.URL = "/" + videoPath
		}

		if coverFile, err := c.FormFile("cover"); err == nil {
			coverDir := "static/cover"
			if err := os.MkdirAll(coverDir, os.ModePerm); err != nil {
				c.JSON(500, serializer.Response{
					Status: 500,
					Msg:    "创建封面目录失败",
				})
				return err
			}
			if video.Cover != "" {
				oldPath := strings.TrimPrefix(video.Cover, "/")
				if _, err := os.Stat(oldPath); err == nil {
					os.Remove(oldPath)
				}
			}
			coverName := fmt.Sprintf("cover_%s_%d%s", user.ID, time.Now().Unix(), filepath.Ext(coverFile.Filename))
			coverPath := filepath.Join(coverDir, coverName)
			if err := c.SaveUploadedFile(coverFile, coverPath); err != nil {
				c.JSON(500, serializer.Response{
					Status: 500,
					Msg:    "保存封面文件失败",
				})
				return err
			}
			video.Cover = "/" + coverPath
		}

		if _, err := service.UpdateVideo(tx, video); err != nil {
			c.JSON(err.Status, err)
			return errors.New(err.Msg)
		}

		c.JSON(200, serializer.BuildVideoResponse(video))
		return nil
	})
}

func DeleteVideo(c *gin.Context) {
	model.Db.Transaction(func(tx *gorm.DB) error {
		userValue, _ := c.Get("user")
		user := userValue.(*model.User)

		vid := c.Param("id")
		if vid == "" {
			c.JSON(400, serializer.Response{
				Status: 400,
				Msg:    "请传入视频id",
			})
			return errors.New("请传入视频id")
		}

		if !service.CompareVidAndUid(user.ID, vid) {
			c.JSON(404, serializer.Response{
				Status: 403,
				Msg:    "没有修改视频权限或者不存在该视频",
			})
			return errors.New("没有修改视频权限或者不存在该视频")
		}

		if video, err := service.DeleteVideo(tx, vid); err != nil {
			c.JSON(err.Status, err)
			return errors.New(err.Msg)
		} else {
			if video.URL != "" {
				oldpath := strings.TrimPrefix(video.URL, "/")
				if _, err := os.Stat(oldpath); err == nil {
					os.Remove(oldpath)
				}
			}
			if video.Cover != "" {
				oldpath := strings.TrimPrefix(video.Cover, "/")
				if _, err := os.Stat(oldpath); err == nil {
					os.Remove(oldpath)
				}
			}
			c.JSON(200, serializer.Response{
				Status: 200,
				Msg:    "视频删除成功",
			})
			return nil

		}

	})

}
