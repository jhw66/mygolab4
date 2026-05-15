// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package comment

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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
		return nil, errors.New("用户未登录")
	}

	comment, err := l.svcCtx.CommentRepo.FindByIDAndVideo(l.ctx, req.Cid, req.Vid)
	if err != nil {
		return nil, errors.New("该评论不存在")
	}
	if comment.UserID != user.ID {
		return nil, errors.New("不能删除别人评论")
	}

	err = l.svcCtx.TransactionRepository.WithTransaction(l.ctx, func(tx *gorm.DB) error {
		if err := l.svcCtx.CommentFavoriteRepo.DeleteByCommentIDWithTx(l.ctx, tx, comment.ID); err != nil {
			l.Errorf("delete comment favorite by comment id failed, cid=%s, err=%v", comment.ID, err)
			return err
		}
		if err := l.svcCtx.CommentRepo.DeleteWithTx(l.ctx, tx, comment); err != nil {
			l.Errorf("delete comment failed, cid=%s, err=%v", comment.ID, err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("删除评论失败")
	}

	if comment.RootID == nil {
		if err := l.svcCtx.CommentCache.InvalidateRootListByVideo(l.ctx, req.Vid); err != nil {
			l.Errorf("invalidate root comment list cache failed, vid=%s, err=%v", req.Vid, err)
		}
		if err := l.svcCtx.CommentCache.InvalidateReplyListByRoot(l.ctx, req.Vid, comment.ID); err != nil {
			l.Errorf("invalidate reply comment list cache failed, root_id=%s, err=%v", comment.ID, err)
		}
	} else {
		if err := l.svcCtx.CommentCache.InvalidateReplyListByRoot(l.ctx, req.Vid, *comment.RootID); err != nil {
			l.Errorf("invalidate reply comment list cache failed, root_id=%s, err=%v", *comment.RootID, err)
		}
	}
	if err := l.svcCtx.FavoriteCache.DelCommentFavoriteCount(l.ctx, comment.ID); err != nil {
		l.Errorf("delete comment favorite count failed, cid=%s, err=%v", comment.ID, err)
	}

	return &types.CommonRsp{Status: 200, Msg: "删除评论成功"}, nil
}
