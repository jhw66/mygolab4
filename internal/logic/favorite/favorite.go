// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package favorite

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/logic/core"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"
	"github.com/jhw66/myvideo_lab4/pkg/db/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type FavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteLogic {
	return &FavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteLogic) Favorite(req *types.FavoriteReq) (resp *types.CommonRsp, err error) {
	user, ok := auth.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return &types.CommonRsp{Status: 401, Msg: "用户未登录"}, errors.New("用户未登录")
	}
	if _, err := l.svcCtx.VideoRepo.FindByID(l.ctx, req.Vid); err != nil {
		return &types.CommonRsp{Status: 404, Msg: "未找到该视频"}, errors.New("未找到该视频")
	}

	if req.CommentId != "" {
		return l.favoriteComment(req, user.ID)
	}
	return l.favoriteVideo(req, user.ID)
}

func (l *FavoriteLogic) favoriteVideo(req *types.FavoriteReq, userID string) (*types.CommonRsp, error) {
	var liked bool
	_, err := l.svcCtx.FavoriteRepo.FindByUserAndVideo(l.ctx, userID, req.Vid)
	switch {
	case err == nil:
		liked = false
		err = l.svcCtx.FavoriteRepo.DeleteByUserAndVideo(l.ctx, userID, req.Vid)
		if err != nil {
			return &types.CommonRsp{Status: 500, Msg: "取消点赞失败"}, errors.New("取消点赞失败")
		}
	case errors.Is(err, gorm.ErrRecordNotFound):
		liked = true
		err = l.svcCtx.FavoriteRepo.Create(l.ctx, &model.Favorite{
			UserID:  userID,
			VideoID: req.Vid,
		})
		if err != nil {
			return &types.CommonRsp{Status: 500, Msg: "点赞失败"}, errors.New("点赞失败")
		}
	default:
		return &types.CommonRsp{Status: 500, Msg: "点赞操作失败"}, errors.New("点赞操作失败")
	}

	if err := core.WarmUpFavoriteCount(l.ctx, l.svcCtx.FavoriteCache, l.svcCtx.FavoriteRepo, req.Vid); err != nil {
		l.Errorf("warmup favorite count failed, vid=%s, err=%v", req.Vid, err)
	}

	if liked {
		if err := l.svcCtx.FavoriteCache.IncrCount(l.ctx, req.Vid); err != nil {
			l.Errorf("incr favorite count failed, vid=%s, err=%v", req.Vid, err)
		}
	} else {
		if err := l.svcCtx.FavoriteCache.DecrCount(l.ctx, req.Vid); err != nil {
			l.Errorf("decr favorite count failed, vid=%s, err=%v", req.Vid, err)
		}
	}

	if err := core.UpdateRankScore(l.ctx, l.svcCtx, req.Vid); err != nil {
		l.Errorf("update rank score failed, vid=%s, err=%v", req.Vid, err)
	}

	msg := "取消点赞成功"
	if liked {
		msg = "点赞成功"
	}
	return &types.CommonRsp{Status: 200, Msg: msg}, nil
}

func (l *FavoriteLogic) favoriteComment(req *types.FavoriteReq, userID string) (*types.CommonRsp, error) {
	comment, err := l.svcCtx.CommentRepo.FindByIDAndVideo(l.ctx, req.CommentId, req.Vid)
	if err != nil {
		return &types.CommonRsp{Status: 404, Msg: "评论不存在"}, errors.New("评论不存在")
	}
	if err := core.WarmUpCommentFavoriteCount(l.ctx, l.svcCtx.CommentCache, l.svcCtx.CommentFavoriteRepo, req.CommentId); err != nil {
		l.Errorf("warmup comment favorite count failed, cid=%s, err=%v", req.CommentId, err)
	}

	liked := false
	err = l.svcCtx.TransactionRepository.WithTransaction(l.ctx, func(tx *gorm.DB) error {
		_, err := l.svcCtx.CommentFavoriteRepo.FindByUserAndCommentWithTx(l.ctx, tx, userID, req.CommentId)
		switch {
		case err == nil:
			liked = false
			return l.svcCtx.CommentFavoriteRepo.DeleteByUserAndCommentWithTx(l.ctx, tx, userID, req.CommentId)
		case errors.Is(err, gorm.ErrRecordNotFound):
			liked = true
			return l.svcCtx.CommentFavoriteRepo.CreateWithTx(l.ctx, tx, &model.CommentFavorite{
				UserID:    userID,
				CommentID: req.CommentId,
			})
		default:
			return err
		}
	})
	if err != nil {
		return &types.CommonRsp{Status: 500, Msg: "评论点赞操作失败"}, errors.New("评论点赞操作失败")
	}

	if liked {
		if err := l.svcCtx.CommentCache.IncrFavoriteCount(l.ctx, req.CommentId); err != nil {
			l.Errorf("incr comment favorite count failed, cid=%s, err=%v", req.CommentId, err)
		}
	} else {
		if err := l.svcCtx.CommentCache.DecrFavoriteCount(l.ctx, req.CommentId); err != nil {
			l.Errorf("decr comment favorite count failed, cid=%s, err=%v", req.CommentId, err)
		}
	}
	if err := l.svcCtx.CommentCache.MarkDirtyFavoriteCount(l.ctx, req.CommentId); err != nil {
		l.Errorf("mark dirty comment favorite count failed, cid=%s, err=%v", req.CommentId, err)
	}
	if comment.CommentID == nil {
		if err := l.svcCtx.CommentCache.InvalidateRootListByVideo(l.ctx, req.Vid); err != nil {
			l.Errorf("invalidate root comment list cache failed, vid=%s, err=%v", req.Vid, err)
		}
	} else if comment.RootID != nil {
		if err := l.svcCtx.CommentCache.InvalidateReplyListByRoot(l.ctx, req.Vid, *comment.RootID); err != nil {
			l.Errorf("invalidate reply comment list cache failed, root_id=%s, err=%v", *comment.RootID, err)
		}
	}

	msg := "取消评论点赞成功"
	if liked {
		msg = "评论点赞成功"
	}
	return &types.CommonRsp{Status: 200, Msg: msg}, nil
}
