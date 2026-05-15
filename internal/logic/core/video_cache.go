package core

import (
	"context"
	"time"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/pkg/cache/cacherepo"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/comment"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/commentfavorite"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/favorite"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	FavoriteWeight       uint = 100
	CommentWeight        uint = 20
	DefaultRankLimit          = 10
	DefaultSyncBatchSize      = 100
	DefaultSyncInterval       = 3 * time.Second

	favoriteCountTTL = 24 * time.Hour
	commentCountTTL  = 24 * time.Hour
)

func CalculateHotScore(favoriteCount, commentCount uint) uint {
	return favoriteCount*FavoriteWeight + commentCount*CommentWeight
}

func UpdateRankScore(ctx context.Context, svcCtx *svc.ServiceContext, vid string) error {
	if err := WarmUpFavoriteCount(ctx, svcCtx.FavoriteCache, svcCtx.FavoriteRepo, vid); err != nil {
		logx.WithContext(ctx).Errorf("warmup favorite count failed, vid=%s, err=%v", vid, err)
	}
	if err := WarmUpCommentCount(ctx, svcCtx.CommentCache, svcCtx.CommentRepo, vid); err != nil {
		logx.WithContext(ctx).Errorf("warmup comment count failed, vid=%s, err=%v", vid, err)
	}

	favCount := uint64(0)
	if v, err := svcCtx.FavoriteCache.GetCountUint64(ctx, vid); err == nil {
		favCount = v
	}
	comCount := uint64(0)
	if v, err := svcCtx.CommentCache.GetCountUint64(ctx, vid); err == nil {
		comCount = v
	}

	score := float64(CalculateHotScore(uint(favCount), uint(comCount)))
	if err := svcCtx.RankCache.ZAddScore(ctx, vid, score); err != nil {
		logx.WithContext(ctx).Errorf("zadd rank score failed, vid=%s, err=%v", vid, err)
	}
	if err := svcCtx.RankCache.MarkDirty(ctx, vid); err != nil {
		logx.WithContext(ctx).Errorf("mark rank dirty failed, vid=%s, err=%v", vid, err)
	}
	return nil
}

func WarmUpRankZSet(ctx context.Context, svcCtx *svc.ServiceContext) error {
	videos, err := svcCtx.VideoRepo.ListRankSeed(ctx)
	if err != nil {
		return err
	}
	if len(videos) == 0 {
		return nil
	}

	entries := make([]redis.Z, 0, len(videos))
	for i := range videos {
		score := float64(CalculateHotScore(videos[i].FavoriteCount, videos[i].CommentCount))
		entries = append(entries, redis.Z{Member: videos[i].ID, Score: score})
	}
	return svcCtx.RankCache.ZBatchAdd(ctx, entries)
}

func SyncDirtyVideoStats(ctx context.Context, svcCtx *svc.ServiceContext, batchSize int) error {
	if batchSize <= 0 {
		batchSize = DefaultSyncBatchSize
	}

	dirtyVids, err := svcCtx.RankCache.ListDirty(ctx)
	if err != nil {
		return err
	}
	if len(dirtyVids) == 0 {
		return nil
	}
	if len(dirtyVids) > batchSize {
		dirtyVids = dirtyVids[:batchSize]
	}

	success := make([]interface{}, 0, len(dirtyVids))
	for _, vid := range dirtyVids {
		if err := WarmUpFavoriteCount(ctx, svcCtx.FavoriteCache, svcCtx.FavoriteRepo, vid); err != nil {
			logx.WithContext(ctx).Errorf("sync warmup favorite failed, vid=%s, err=%v", vid, err)
			continue
		}
		if err := WarmUpCommentCount(ctx, svcCtx.CommentCache, svcCtx.CommentRepo, vid); err != nil {
			logx.WithContext(ctx).Errorf("sync warmup comment failed, vid=%s, err=%v", vid, err)
			continue
		}

		favCount, err := svcCtx.FavoriteCache.GetCountUint64(ctx, vid)
		if err != nil {
			logx.WithContext(ctx).Errorf("sync get favorite cache failed, vid=%s, err=%v", vid, err)
			continue
		}
		commentCount, err := svcCtx.CommentCache.GetCountUint64(ctx, vid)
		if err != nil {
			logx.WithContext(ctx).Errorf("sync get comment cache failed, vid=%s, err=%v", vid, err)
			continue
		}

		hotScore := CalculateHotScore(uint(favCount), uint(commentCount))
		if err := svcCtx.VideoRepo.UpdateByID(ctx, vid, favCount, commentCount, hotScore); err != nil {
			logx.WithContext(ctx).Errorf("sync update video stats failed, vid=%s, err=%v", vid, err)
			continue
		}
		if err := svcCtx.RankCache.ZAddScore(ctx, vid, float64(hotScore)); err != nil {
			logx.WithContext(ctx).Errorf("sync refresh rank score failed, vid=%s, err=%v", vid, err)
			continue
		}

		success = append(success, vid)
	}

	return svcCtx.RankCache.RemoveDirtyBatch(ctx, success)
}

func StartVideoStatSync(ctx context.Context, svcCtx *svc.ServiceContext, interval time.Duration, batchSize int) {
	if interval <= 0 {
		interval = DefaultSyncInterval
	}
	if batchSize <= 0 {
		batchSize = DefaultSyncBatchSize
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := SyncDirtyVideoStats(ctx, svcCtx, batchSize); err != nil {
				logx.WithContext(ctx).Errorf("sync dirty video stats failed, err=%v", err)
			}
		}
	}
}

func SyncDirtyCommentFavoriteStats(ctx context.Context, svcCtx *svc.ServiceContext, batchSize int) error {
	if batchSize <= 0 {
		batchSize = DefaultSyncBatchSize
	}

	dirtyCommentIDs, err := svcCtx.FavoriteCache.ListDirtyCommentFavoriteCount(ctx)
	if err != nil {
		return err
	}
	if len(dirtyCommentIDs) == 0 {
		return nil
	}
	if len(dirtyCommentIDs) > batchSize {
		dirtyCommentIDs = dirtyCommentIDs[:batchSize]
	}

	success := make([]interface{}, 0, len(dirtyCommentIDs))
	for _, commentID := range dirtyCommentIDs {
		if err := WarmUpCommentFavoriteCount(ctx, svcCtx.FavoriteCache, svcCtx.CommentFavoriteRepo, commentID); err != nil {
			logx.WithContext(ctx).Errorf("sync warmup comment favorite failed, cid=%s, err=%v", commentID, err)
			continue
		}
		favoriteCount, err := svcCtx.FavoriteCache.GetCommentFavoriteCount(ctx, commentID)
		if err != nil {
			logx.WithContext(ctx).Errorf("sync get comment favorite cache failed, cid=%s, err=%v", commentID, err)
			continue
		}
		if err := svcCtx.CommentRepo.UpdateFavoriteCountByID(ctx, commentID, favoriteCount); err != nil {
			logx.WithContext(ctx).Errorf("sync update comment favorite count failed, cid=%s, err=%v", commentID, err)
			continue
		}
		success = append(success, commentID)
	}

	return svcCtx.FavoriteCache.RemoveDirtyCommentFavoriteCountBatch(ctx, success)
}

func StartCommentFavoriteStatSync(ctx context.Context, svcCtx *svc.ServiceContext, interval time.Duration, batchSize int) {
	if interval <= 0 {
		interval = DefaultSyncInterval
	}
	if batchSize <= 0 {
		batchSize = DefaultSyncBatchSize
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := SyncDirtyCommentFavoriteStats(ctx, svcCtx, batchSize); err != nil {
				logx.WithContext(ctx).Errorf("sync dirty comment favorite stats failed, err=%v", err)
			}
		}
	}
}

func WarmUpFavoriteCount(ctx context.Context, cc cacherepo.FavoriteCache, favoriteRepo favorite.FavoriteRepository, vid string) error {
	exists, _ := cc.ExistsCount(ctx, vid)
	if exists {
		return nil
	}
	if !cc.TryWarmupLock(ctx, vid, 5*time.Second) {
		time.Sleep(100 * time.Millisecond)
		return nil
	}
	exists, _ = cc.ExistsCount(ctx, vid)
	if exists {
		return nil
	}

	count, err := favoriteRepo.CountByVideoID(ctx, vid)
	if err != nil {
		return err
	}
	return cc.SetCount(ctx, vid, count, favoriteCountTTL)
}

func WarmUpCommentCount(ctx context.Context, cc cacherepo.CommentCache, commentRepo comment.CommentRepository, vid string) error {
	exists, _ := cc.ExistsCount(ctx, vid)
	if exists {
		return nil
	}
	if !cc.TryWarmupLock(ctx, "comment_count:"+vid, 5*time.Second) {
		time.Sleep(100 * time.Millisecond)
		return nil
	}
	exists, _ = cc.ExistsCount(ctx, vid)
	if exists {
		return nil
	}

	count, err := commentRepo.CountByVideoID(ctx, vid)
	if err != nil {
		return err
	}
	return cc.SetCount(ctx, vid, count, commentCountTTL)
}

func WarmUpRootCommentCount(ctx context.Context, cc cacherepo.CommentCache, commentRepo comment.CommentRepository, vid string) error {
	exists, _ := cc.ExistsRootCount(ctx, vid)
	if exists {
		return nil
	}
	if !cc.TryWarmupLock(ctx, "comment_root_count:"+vid, 5*time.Second) {
		time.Sleep(100 * time.Millisecond)
		return nil
	}
	exists, _ = cc.ExistsRootCount(ctx, vid)
	if exists {
		return nil
	}
	count, err := commentRepo.CountRootByVideoID(ctx, vid)
	if err != nil {
		return err
	}
	return cc.SetRootCount(ctx, vid, count, commentCountTTL)
}

func WarmUpReplyCommentCount(ctx context.Context, cc cacherepo.CommentCache, commentRepo comment.CommentRepository, vid string, rootID string) error {
	exists, _ := cc.ExistsReplyCount(ctx, rootID)
	if exists {
		return nil
	}
	if !cc.TryWarmupLock(ctx, "comment_reply_count:"+rootID, 5*time.Second) {
		time.Sleep(100 * time.Millisecond)
		return nil
	}
	exists, _ = cc.ExistsReplyCount(ctx, rootID)
	if exists {
		return nil
	}
	count, err := commentRepo.CountRepliesByRootID(ctx, vid, rootID)
	if err != nil {
		return err
	}
	return cc.SetReplyCount(ctx, rootID, count, commentCountTTL)
}

func WarmUpCommentFavoriteCount(ctx context.Context, cc cacherepo.FavoriteCache, commentFavoriteRepo commentfavorite.CommentFavoriteRepository, commentID string) error {
	exists, _ := cc.ExistsCommentFavoriteCount(ctx, commentID)
	if exists {
		return nil
	}
	if !cc.TryCommentFavoriteWarmupLock(ctx, commentID, 5*time.Second) {
		time.Sleep(100 * time.Millisecond)
		return nil
	}
	exists, _ = cc.ExistsCommentFavoriteCount(ctx, commentID)
	if exists {
		return nil
	}
	count, err := commentFavoriteRepo.CountByCommentID(ctx, commentID)
	if err != nil {
		return err
	}
	return cc.SetCommentFavoriteCount(ctx, commentID, count, commentCountTTL)
}
