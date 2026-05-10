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

	"github.com/zeromicro/go-zero/core/logx"
)

type DelCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelCommentLogic {
	return &DelCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelCommentLogic) DelComment(req *types.DelCommentReq) (resp *types.CommonRsp, err error) {
	user, ok := auth.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return &types.CommonRsp{Status: 401, Msg: "用户未登录"}, errors.New("用户未登录")
	}

	comment, err := l.svcCtx.CommentRepo.FindByIDAndVideo(l.ctx, req.Cid, req.Vid)
	if err != nil {
		return &types.CommonRsp{Status: 404, Msg: "该评论不存在"}, errors.New("该评论不存在")
	}
	if comment.UserID != user.ID {
		return &types.CommonRsp{Status: 403, Msg: "不能删除别人评论"}, errors.New("不能删除别人评论")
	}

	if err := core.WarmUpCommentCount(l.ctx, l.svcCtx.CommentCache, l.svcCtx.VideoRepo, req.Vid); err != nil {
		l.Errorf("warmup comment count failed, vid=%s, err=%v", req.Vid, err)
	}

	if err = l.svcCtx.CommentRepo.Delete(l.ctx, comment); err != nil {
		return &types.CommonRsp{Status: 500, Msg: "删除评论失败"}, errors.New("删除评论失败")
	}

	if err := l.svcCtx.CommentCache.DecrCount(l.ctx, req.Vid); err != nil {
		l.Errorf("decr comment count failed, vid=%s, err=%v", req.Vid, err)
	}
	if err := core.UpdateRankScore(l.ctx, l.svcCtx, req.Vid); err != nil {
		l.Errorf("update rank score failed, vid=%s, err=%v", req.Vid, err)
	}
	if err := l.svcCtx.CommentCache.InvalidateListByVideo(l.ctx, req.Vid); err != nil {
		l.Errorf("invalidate comment list cache failed, vid=%s, err=%v", req.Vid, err)
	}
	return &types.CommonRsp{Status: 200, Msg: "删除评论成功"}, nil
}
