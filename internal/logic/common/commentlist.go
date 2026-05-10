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
		return &types.CommentListRsp{Status: 404, Msg: "未找到该视频"}, errors.New("未找到该视频")
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
		rootComment, err := l.svcCtx.CommentRepo.FindByIDAndVideo(l.ctx, req.CommentId, req.Vid)
		if err != nil || rootComment.CommentID != nil {
			return &types.CommentListRsp{Status: 404, Msg: "根评论不存在"}, errors.New("根评论不存在")
		}
	}

	if cachedComments, err := l.svcCtx.CommentCache.GetList(l.ctx, req.Vid, req.CommentId, page, pageSize); err == nil {
		var data cachemodel.CommentListCacheData
		if err := json.Unmarshal([]byte(cachedComments), &data); err == nil {
			return serializer.CommentListRspFromModels(data.Comments, data.Total, page, pageSize), nil
		}
	}

	var comments []model.Comment
	var total int64
	if req.CommentId == "" {
		comments, err = l.svcCtx.CommentRepo.ListRootByVideoID(l.ctx, req.Vid, page, pageSize)
	} else {
		comments, err = l.svcCtx.CommentRepo.ListRepliesByRootID(l.ctx, req.Vid, req.CommentId, page, pageSize)
	}
	if err != nil {
		return &types.CommentListRsp{Status: 500, Msg: "查询评论失败"}, errors.New("查询评论失败")
	}

	if err := core.WarmUpCommentCount(l.ctx, l.svcCtx.CommentCache, l.svcCtx.VideoRepo, req.Vid); err != nil {
		l.Errorf("warmup comment count failed, vid=%s, err=%v", req.Vid, err)
	}

	if req.CommentId == "" {
		total, err = l.svcCtx.CommentRepo.CountRootByVideoID(l.ctx, req.Vid)
	} else {
		total, err = l.svcCtx.CommentRepo.CountRepliesByRootID(l.ctx, req.Vid, req.CommentId)
	}
	if err != nil {
		return &types.CommentListRsp{Status: 500, Msg: "查询评论数量失败"}, errors.New("查询评论数量失败")
	}

	cacheData := cachemodel.CommentListCacheData{
		Comments: comments,
		Total:    total,
	}
	if payload, err := json.Marshal(cacheData); err == nil {
		_ = l.svcCtx.CommentCache.SetList(l.ctx, req.Vid, req.CommentId, page, pageSize, payload, 30*time.Second)
	}
	return serializer.CommentListRspFromModels(comments, total, page, pageSize), nil
}
