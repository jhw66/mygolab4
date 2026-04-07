package service

import (
	"log"
	"time"

	"github.com/jhw66/myvideo_lab4/cache"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/redis/go-redis/v9"
)

const (
	favoriteWeight uint = 100
	commentWeight  uint = 20

	DefaultRankLimit    = 10
	rankZSetKey         = "rank:video:hot"
	rankDirtyVideoKey   = "rank:dirty_videos"
	SyncInterval        = 30 * time.Second
	syncDirtyBatchLimit = 200
)

// 计算热度分数
func CalculateHotScore(favoriteCount, commentCount uint) uint {
	return favoriteCount*favoriteWeight + commentCount*commentWeight
}

// 更新排行榜分数
func UpdateRankScore(vid string) error {
	// 预热两个计数 key：点赞/评论只触发一方的预热，另一方可能还不在 Redis 中
	warmUpFavoriteCount(vid)
	warmUpCommentCount(vid)

	favCount := uint64(0)
	if val, err := cache.Rdb.Get(cache.Ctx, buildFavoriteCountKey(vid)).Uint64(); err == nil {
		favCount = val
	}
	comCount := uint64(0)
	if val, err := cache.Rdb.Get(cache.Ctx, buildCommentCountKey(vid)).Uint64(); err == nil {
		comCount = val
	}
	score := float64(CalculateHotScore(uint(favCount), uint(comCount)))
	cache.Rdb.ZAdd(cache.Ctx, rankZSetKey, redis.Z{
		Score:  score,
		Member: vid,
	})
	cache.Rdb.SAdd(cache.Ctx, rankDirtyVideoKey, vid)
	return nil
}

// 预热排行榜ZSET
func WarmUpRankZSet() {
	var videos []model.Video
	if err := model.Db.Select("id,favorite_count,comment_count").Find(&videos).Error; err != nil {
		log.Println("排行榜ZSET预热失败:", err)
		return
	}

	pipe := cache.Rdb.Pipeline()
	for _, v := range videos {
		pipe.ZAdd(cache.Ctx, rankZSetKey, redis.Z{
			Score:  float64(CalculateHotScore(v.FavoriteCount, v.CommentCount)),
			Member: v.ID,
		})
	}
	if _, err := pipe.Exec(cache.Ctx); err != nil {
		log.Println("排行榜ZSET预热Pipeline执行失败:", err)
		return
	}

	log.Printf("排行榜ZSET预热完成,视频数:%d", len(videos))
}

// 获取排行榜
func GetRankVideos(limit int) (*[]model.Video, error) {
	vids, err := GetTopRankVideoIDs(limit)
	if err != nil || len(vids) == 0 {
		return nil, err
	}

	//可能mysql中的点赞/评论数与redis中的不一致，因为redis中的数据是异步更新的
	//但是主要是从mysql中获取url等信息，不影响排名
	var videos []model.Video
	if err := model.Db.Where("id IN ?", vids).Find(&videos).Error; err != nil {
		return nil, err
	}

	videoMap := make(map[string]model.Video, len(videos))
	for _, v := range videos {
		videoMap[v.ID] = v
	}
	orderedVideos := make([]model.Video, 0, len(videos))
	for _, vid := range vids {
		// 检查 key 是否存在，防止已删除视频产生零值对象
		if v, ok := videoMap[vid]; ok {
			orderedVideos = append(orderedVideos, v)
		}
	}

	return &orderedVideos, nil
}

// 获取排行榜视频ID
func GetTopRankVideoIDs(limit int) ([]string, error) {
	results, err := cache.Rdb.ZRangeArgs(cache.Ctx, redis.ZRangeArgs{
		Key:   rankZSetKey,
		Start: 0,
		Stop:  int64(limit - 1),
		Rev:   true,
	}).Result()
	if err != nil {
		return nil, err
	}

	vids := make([]string, 0, len(results))
	for _, member := range results {
		vids = append(vids, member)
	}
	return vids, nil
}

// 将dirty集合中的视频ID同步到MySQL
func SyncDirtyToMySQL() {
	allVids, err := cache.Rdb.SMembers(cache.Ctx, rankDirtyVideoKey).Result()
	if err != nil {
		log.Printf("SyncDirtyToMySQL: SMEMBERS failed with error: %v", err)
		return
	}
	if len(allVids) == 0 {
		return
	}

	cleanVids := allVids
	if len(cleanVids) > syncDirtyBatchLimit {
		cleanVids = cleanVids[:syncDirtyBatchLimit]
	}

	successVIDs := make([]interface{}, 0, len(cleanVids))
	for _, vid := range cleanVids {
		favCount, err := cache.Rdb.Get(cache.Ctx, buildFavoriteCountKey(vid)).Uint64()
		if err != nil {
			// key 已过期（TTL 到期或被驱逐），无法获取准确值
			// 将 vid 加入 successVIDs 以清除 dirty 标记，避免该条目永久积压在 dirty 集合中
			log.Printf("SyncDirtyToMySQL: 点赞数 key 已过期，跳过同步 vid=%s", vid)
			successVIDs = append(successVIDs, vid)
			continue
		}
		comCount, err := cache.Rdb.Get(cache.Ctx, buildCommentCountKey(vid)).Uint64()
		if err != nil {
			// 同上，评论数 key 已过期
			log.Printf("SyncDirtyToMySQL: 评论数 key 已过期，跳过同步 vid=%s", vid)
			successVIDs = append(successVIDs, vid)
			continue
		}
		// 使用 map 而非 struct，避免零值字段被 GORM 跳过不更新
		if err := model.Db.Model(&model.Video{}).
			Where("id = ?", vid).
			Updates(map[string]interface{}{
				"favorite_count": favCount,
				"comment_count":  comCount,
				"hot_score":      CalculateHotScore(uint(favCount), uint(comCount)),
			}).Error; err != nil {
			log.Printf("SyncDirtyToMySQL: 同步排行榜统计失败 vid=%s err=%v", vid, err)
			continue
		}

		successVIDs = append(successVIDs, vid)
	}

	if len(successVIDs) > 0 {
		if err := cache.Rdb.SRem(cache.Ctx, rankDirtyVideoKey, successVIDs...).Err(); err != nil {
			log.Printf("SyncDirtyToMySQL: 批量删除 dirty 标记失败 err=%v", err)
		}
	}
}
