package service

import (
	"log"
	"time"

	"github.com/jhw66/myvideo_lab4/cache"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/redis/go-redis/v9"
)

const (
	FavoriteWeight uint64 = 100
	CommentWeight  uint64 = 1

	DefaultRankLimit  = 10
	RankZSetKey       = "rank:video:hot"
	rankDirtyVideoKey = "rank:dirty_videos"
	SyncInterval      = 30 * time.Second
)

func CalculateHotScore(favoriteCount, commentCount uint) uint64 {
	return uint64(favoriteCount)*FavoriteWeight + uint64(commentCount)*CommentWeight
}

func UpdateRankScore(vid string) {
	favCount := uint64(0)
	if val, err := cache.Rdb.Get(cache.Ctx, buildFavoriteCountKey(vid)).Uint64(); err == nil {
		favCount = val
	}
	comCount := uint64(0)
	if val, err := cache.Rdb.Get(cache.Ctx, buildCommentCountKey(vid)).Uint64(); err == nil {
		comCount = val
	}
	score := float64(CalculateHotScore(uint(favCount), uint(comCount)))
	cache.Rdb.ZAdd(cache.Ctx, RankZSetKey, redis.Z{
		Score:  score,
		Member: vid,
	})
	cache.Rdb.SAdd(cache.Ctx, rankDirtyVideoKey, vid)
}

func WarmUpRankZSet() {
	var videos []model.Video
	if err := model.Db.Select("id, favorite_count, comment_count, hot_score").Find(&videos).Error; err != nil {
		log.Println("排行榜ZSET预热失败:", err)
		return
	}
	pipe := cache.Rdb.Pipeline()
	for _, v := range videos {
		pipe.ZAdd(cache.Ctx, RankZSetKey, redis.Z{
			Score:  float64(v.HotScore),
			Member: v.ID,
		})
	}
	if _, err := pipe.Exec(cache.Ctx); err != nil {
		log.Println("排行榜ZSET预热Pipeline执行失败:", err)
		return
	}
	log.Printf("排行榜ZSET预热完成, 视频数: %d\n", len(videos))
}

func GetTopRankVideoIDs(limit int) ([]string, error) {
	results, err := cache.Rdb.ZRevRangeWithScores(cache.Ctx, RankZSetKey, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(results))
	for _, z := range results {
		memberStr, ok := z.Member.(string)
		if !ok {
			continue
		}
		ids = append(ids, memberStr)
	}
	return ids, nil
}

func GetRankVideos(limit int) ([]model.Video, error) {
	ids, err := GetTopRankVideoIDs(limit)
	if err != nil || len(ids) == 0 {
		var videos []model.Video
		err := model.Db.Order("hot_score desc").Order("created_at desc").Order("id desc").
			Limit(limit).Find(&videos).Error
		return videos, err
	}
	var videos []model.Video
	if err := model.Db.Where("id IN ?", ids).Find(&videos).Error; err != nil {
		return nil, err
	}
	videoMap := make(map[string]model.Video, len(videos))
	for _, v := range videos {
		videoMap[v.ID] = v
	}
	ordered := make([]model.Video, 0, len(ids))
	for _, id := range ids {
		if v, ok := videoMap[id]; ok {
			ordered = append(ordered, v)
		}
	}
	return ordered, nil
}

func SyncDirtyToMySQL() {
	for {
		rawVID, err := cache.Rdb.SPop(cache.Ctx, rankDirtyVideoKey).Result()
		if err != nil || rawVID == "" {
			break
		}
		vid := rawVID
		favCount := uint(0)
		if val, err := cache.Rdb.Get(cache.Ctx, buildFavoriteCountKey(vid)).Uint64(); err == nil {
			favCount = uint(val)
		}
		comCount := uint(0)
		if val, err := cache.Rdb.Get(cache.Ctx, buildCommentCountKey(vid)).Uint64(); err == nil {
			comCount = uint(val)
		}
		model.Db.Model(&model.Video{}).Where("id = ?", vid).Updates(map[string]interface{}{
			"favorite_count": favCount,
			"comment_count":  comCount,
			"hot_score":      CalculateHotScore(favCount, comCount),
		})
	}
}
