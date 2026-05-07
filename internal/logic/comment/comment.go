// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package comment

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/logic/core"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"github.com/jhw66/myvideo_lab4/pkg/utlcontext"
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
	if err := utils.ValidateRuneLength(req.Content, 1, 50, "评论内容长度需在1-50个字符"); err != nil {
		return &types.CommonRsp{Status: 400, Msg: err.Error()}, err
	}

	user, ok := utlcontext.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return &types.CommonRsp{Status: 401, Msg: "用户未登录"}, errors.New("用户未登录")
	}
	if _, err := l.svcCtx.VideoRepo.FindByID(l.ctx, req.Vid); err != nil {
		return &types.CommonRsp{Status: 404, Msg: "未找到该视频"}, errors.New("未找到该视频")
	}

	if err := core.WarmUpCommentCount(l.ctx, l.svcCtx.CommentCache, l.svcCtx.VideoRepo, req.Vid); err != nil {
		l.Errorf("warmup comment count failed, vid=%s, err=%v", req.Vid, err)
	}

	if err = l.svcCtx.CommentRepo.Create(l.ctx, &model.Comment{
		UserID:  user.ID,
		VideoID: req.Vid,
		Content: req.Content,
	}); err != nil {
		return &types.CommonRsp{Status: 500, Msg: "评论保存失败"}, errors.New("评论保存失败")
	}

	if err := l.svcCtx.CommentCache.IncrCount(l.ctx, req.Vid); err != nil {
		l.Errorf("incr comment count failed, vid=%s, err=%v", req.Vid, err)
	}
	if err := core.UpdateRankScore(l.ctx, l.svcCtx, req.Vid); err != nil {
		l.Errorf("update rank score failed, vid=%s, err=%v", req.Vid, err)
	}
	if err := l.svcCtx.CommentCache.InvalidateListByVideo(l.ctx, req.Vid); err != nil {
		l.Errorf("invalidate comment list cache failed, vid=%s, err=%v", req.Vid, err)
	}

	return &types.CommonRsp{Status: 200, Msg: "评论成功"}, nil
}
