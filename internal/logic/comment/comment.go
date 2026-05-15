// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package comment

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/logic/core"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"
	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"github.com/jhw66/myvideo_lab4/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentLogic {
	return &CommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentLogic) Comment(req *types.CommentAddReq) (resp *types.CommonRsp, err error) {
	if err := utils.ValidateRuneLength(req.Content, 1, 50); err != nil {
		return nil, err
	}

	user, ok := auth.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return nil, errors.New("用户未登录")
	}
	if _, err := l.svcCtx.VideoRepo.FindByID(l.ctx, req.Vid); err != nil {
		return nil, errors.New("未找到该视频")
	}

	var commentID *string
	var rootID *string
	if req.CommentId != "" {
		parent, err := l.svcCtx.CommentRepo.FindByIDAndVideo(l.ctx, req.CommentId, req.Vid)
		if err != nil {
			return nil, errors.New("父评论不存在")
		}
		commentID = &parent.ID
		if parent.RootID != nil {
			rootID = parent.RootID
		} else {
			rootID = &parent.ID
		}
	}

	if err := core.WarmUpCommentCount(l.ctx, l.svcCtx.CommentCache, l.svcCtx.CommentRepo, req.Vid); err != nil {
		l.Errorf("warmup comment count failed, vid=%s, err=%v", req.Vid, err)
	}
	if rootID == nil {
		if err := core.WarmUpRootCommentCount(l.ctx, l.svcCtx.CommentCache, l.svcCtx.CommentRepo, req.Vid); err != nil {
			l.Errorf("warmup root comment count failed, vid=%s, err=%v", req.Vid, err)
		}
	} else {
		if err := core.WarmUpReplyCommentCount(l.ctx, l.svcCtx.CommentCache, l.svcCtx.CommentRepo, req.Vid, *rootID); err != nil {
			l.Errorf("warmup reply comment count failed, root_id=%s, err=%v", *rootID, err)
		}
	}

	if err = l.svcCtx.CommentRepo.Create(l.ctx, &model.Comment{
		UserID:    user.ID,
		VideoID:   req.Vid,
		CommentID: commentID,
		RootID:    rootID,
		Content:   req.Content,
	}); err != nil {
		return nil, errors.New("评论保存失败")
	}

	if err := l.svcCtx.CommentCache.IncrCount(l.ctx, req.Vid); err != nil {
		l.Errorf("incr comment count failed, vid=%s, err=%v", req.Vid, err)
	}

	if rootID == nil {
		if err := l.svcCtx.CommentCache.IncrRootCount(l.ctx, req.Vid); err != nil {
			l.Errorf("incr root comment count failed, vid=%s, err=%v", req.Vid, err)
		}
		if err := l.svcCtx.CommentCache.InvalidateRootListByVideo(l.ctx, req.Vid); err != nil {
			l.Errorf("invalidate root comment list cache failed, vid=%s, err=%v", req.Vid, err)
		}
	} else {
		if err := l.svcCtx.CommentCache.IncrReplyCount(l.ctx, *rootID); err != nil {
			l.Errorf("incr reply comment count failed, root_id=%s, err=%v", *rootID, err)
		}
		if err := l.svcCtx.CommentCache.InvalidateReplyListByRoot(l.ctx, req.Vid, *rootID); err != nil {
			l.Errorf("invalidate reply comment list cache failed, root_id=%s, err=%v", *rootID, err)
		}
	}

	if err := core.UpdateRankScore(l.ctx, l.svcCtx, req.Vid); err != nil {
		l.Errorf("update rank score failed, vid=%s, err=%v", req.Vid, err)
	}

	return &types.CommonRsp{Status: 200, Msg: "评论成功"}, nil
}
