// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jhw66/myvideo_lab4/internal/logic/core"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/cache/cachemodel"
	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"github.com/jhw66/myvideo_lab4/pkg/serializer"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListReq) (resp *types.CommentListRsp, err error) {
	if _, err := l.svcCtx.VideoRepo.FindByID(l.ctx, req.Vid); err != nil {
		return nil, errors.New("未找到该视频")
	}

	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	if req.CommentId != "" {
		rootComment, err := l.svcCtx.CommentRepo.FindByIDAndVideoUnscoped(l.ctx, req.CommentId, req.Vid)
		if err != nil || rootComment.CommentID != nil {
			return nil, errors.New("根评论不存在")
		}
	}

	if cachedComments, err := l.svcCtx.CommentCache.GetList(l.ctx, req.Vid, req.CommentId, page, pageSize); err == nil {
		var data cachemodel.CommentListCacheData
		if err := json.Unmarshal([]byte(cachedComments), &data); err == nil {
			return serializer.CommentListRspFromModels(data.Comments, data.Total, page, pageSize), nil
		}
	}

	var comments []model.Comment
	if req.CommentId == "" {
		comments, err = l.svcCtx.CommentRepo.ListRootByVideoID(l.ctx, req.Vid, page, pageSize)
	} else {
		comments, err = l.svcCtx.CommentRepo.ListRepliesByRootID(l.ctx, req.Vid, req.CommentId, page, pageSize)
	}
	if err != nil {
		return nil, errors.New("查询评论失败")
	}

	l.setFavoriteByCache(comments)
	total := l.getCachedListTotal(req)
	cacheData := cachemodel.CommentListCacheData{
		Comments: comments,
		Total:    total,
	}
	if payload, err := json.Marshal(cacheData); err == nil {
		if err := l.svcCtx.CommentCache.SetList(l.ctx, req.Vid, req.CommentId, page, pageSize, payload, 30*time.Second); err != nil {
			l.Errorf("failed to set comment list in cache, vid=%s, comment_id=%s, err=%v", req.Vid, req.CommentId, err)
		}
	}
	return serializer.CommentListRspFromModels(comments, total, page, pageSize), nil
}

func (l *CommentListLogic) setFavoriteByCache(comments []model.Comment) {
	for i := range comments {
		commentID := comments[i].ID
		if err := core.WarmUpCommentFavoriteCount(l.ctx, l.svcCtx.FavoriteCache, l.svcCtx.CommentFavoriteRepo, commentID); err != nil {
			l.Errorf("warmup comment favorite count failed, cid=%s, err=%v", commentID, err)
			continue
		}
		count, err := l.svcCtx.FavoriteCache.GetCommentFavoriteCount(l.ctx, commentID)
		if err != nil {
			l.Errorf("get comment favorite count cache failed, cid=%s, err=%v", commentID, err)
			continue
		}
		if count < 0 {
			count = 0
		}
		comments[i].FavoriteCount = uint(count)
	}
}

func (l *CommentListLogic) getCachedListTotal(req *types.CommentListReq) int64 {
	if req.CommentId == "" {
		if err := core.WarmUpRootCommentCount(l.ctx, l.svcCtx.CommentCache, l.svcCtx.CommentRepo, req.Vid); err != nil {
			l.Errorf("warmup root comment count failed, vid=%s, err=%v", req.Vid, err)
		}
		if total, err := l.svcCtx.CommentCache.GetRootCount(l.ctx, req.Vid); err == nil {
			return total
		}
		total, err := l.svcCtx.CommentRepo.CountRootByVideoID(l.ctx, req.Vid)
		if err != nil {
			l.Errorf("query root comment count failed, vid=%s, err=%v", req.Vid, err)
			return 0
		}
		return total
	}

	if err := core.WarmUpReplyCommentCount(l.ctx, l.svcCtx.CommentCache, l.svcCtx.CommentRepo, req.Vid, req.CommentId); err != nil {
		l.Errorf("warmup reply comment count failed, root_id=%s, err=%v", req.CommentId, err)
	}
	if total, err := l.svcCtx.CommentCache.GetReplyCount(l.ctx, req.CommentId); err == nil {
		return total
	}
	total, err := l.svcCtx.CommentRepo.CountRepliesByRootID(l.ctx, req.Vid, req.CommentId)
	if err != nil {
		l.Errorf("query reply comment count failed, root_id=%s, err=%v", req.CommentId, err)
		return 0
	}
	return total
}
