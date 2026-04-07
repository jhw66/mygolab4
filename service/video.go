package service

import (
	"github.com/jhw66/myvideo_lab4/cache"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type VideoSearch struct {
	KeyWord string `form:"key_word" json:"key_word" binding:"required,max=10"`
}

func (videosearch VideoSearch) FindVideosByKeyword() (*[]model.Video, *serializer.Response) {
	var videos []model.Video
	var count int64
	if err := model.Db.Where("title like ?", "%"+videosearch.KeyWord+"%").Find(&videos).
		Count(&count).Error; err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "查找失败",
		}
	}
	if count == 0 {
		return nil, &serializer.Response{
			Status: 404,
			Msg:    "未找到相关视频",
		}
	}
	return &videos, nil
}

func DeleteVideo(tx *gorm.DB, vid string) (*model.Video, *serializer.Response) {
	var video model.Video
	tx.Where("id = ?", vid).Take(&video)
	if err := tx.Delete(&video).Error; err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "删除视频失败",
		}
	}
	//删除redis缓存
	deleteVideoFromCache(vid)
	return &video, nil
}

func deleteVideoFromCache(vid string) {
	// 清理排行榜 ZSet 中该视频的记录(rankZSetKey= "rank:video:hot")
	cache.Rdb.ZRem(cache.Ctx, rankZSetKey, vid)
	// 清理 dirty 集合，防止 SyncDirtyToMySQL 尝试同步已删除的视频(rankDirtyVideoKey= "rank:dirty_videos")
	cache.Rdb.SRem(cache.Ctx, rankDirtyVideoKey, vid)
	// 清理点赞数缓存(buildFavoriteCountKey= "favorite_count:video:" + vid)
	cache.Rdb.Del(cache.Ctx, buildFavoriteCountKey(vid))
	// 清理评论数缓存(buildCommentCountKey= "comment_count:video:" + vid)
	cache.Rdb.Del(cache.Ctx, buildCommentCountKey(vid))
	// 清理所有分页评论列表缓存(invalidateCommentListCache= "comment_list:video:" + vid + ":p" + strconv.Itoa(page) + ":s" + strconv.Itoa(pageSize))
	invalidateCommentListCache(vid)
}

func FindVideoByUser(user *model.User) (*[]model.Video, *serializer.Response) {
	var videos []model.Video
	if err := model.Db.Model(user).Association("Videos").Find(&videos); err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "查找失败",
		}
	}
	return &videos, nil
}

func FindVideoByVid(vid string) (*model.Video, *serializer.Response) {
	var video model.Video
	if err := model.Db.Where("id = ?", vid).Take(&video).Error; err != nil {
		return nil, &serializer.Response{
			Status: 404,
			Msg:    "未找到该视频",
		}
	}
	return &video, nil
}

func UploadVideo(tx *gorm.DB, video *model.Video) (*model.Video, *serializer.Response) {
	if err := tx.Create(video).Error; err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "上传视频失败",
		}
	}
	cache.Rdb.ZAdd(cache.Ctx, rankZSetKey, redis.Z{
		Score:  float64(video.HotScore),
		Member: video.ID,
	})
	return video, nil
}

func UpdateVideo(tx *gorm.DB, video *model.Video) (*model.Video, *serializer.Response) {
	video.HotScore = uint64(CalculateHotScore(video.FavoriteCount, video.CommentCount))
	if err := tx.Save(video).Error; err != nil {
		return nil, &serializer.Response{
			Status: 500,
			Msg:    "更新视频失败",
		}
	}
	cache.Rdb.ZAdd(cache.Ctx, rankZSetKey, redis.Z{
		Score:  float64(video.HotScore),
		Member: video.ID,
	})
	return video, nil
}
