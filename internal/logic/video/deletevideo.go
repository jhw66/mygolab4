// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/utils"
	"github.com/jhw66/myvideo_lab4/pkg/utlcontext"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteVideoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteVideoLogic {
	return &DeleteVideoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteVideoLogic) DeleteVideo(req *types.VideoIdReq) (resp *types.CommonRsp, err error) {
	user, ok := utlcontext.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return &types.CommonRsp{Status: 401, Msg: "用户未登录"}, errors.New("用户未登录")
	}
	if req.Id == "" {
		return &types.CommonRsp{Status: 400, Msg: "请传入视频id"}, errors.New("请传入视频id")
	}

	video, err := l.svcCtx.VideoRepo.FindByID(l.ctx, req.Id)
	if err != nil {
		return &types.CommonRsp{Status: 404, Msg: "未找到该视频"}, errors.New("未找到该视频")
	}
	if video.UserID != user.ID {
		return &types.CommonRsp{Status: 403, Msg: "没有修改视频权限或者不存在该视频"}, errors.New("没有修改视频权限或者不存在该视频")
	}

	err = l.svcCtx.TransactionRepository.WithTransaction(l.ctx, func(tx *gorm.DB) error {
		deleted, err := l.svcCtx.VideoRepo.DeleteByIDWithTx(l.ctx, tx, req.Id)
		if err != nil {
			return err
		}
		if err := utils.RemoveIfExistsWithUrl(deleted.URL); err != nil {
			return err
		}
		if err := utils.RemoveIfExistsWithUrl(deleted.Cover); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &types.CommonRsp{Status: 500, Msg: "视频删除失败"}, errors.New("视频删除失败")
	}

	deleteVideoCaches(l.ctx, l.svcCtx, req.Id)

	return &types.CommonRsp{Status: 200, Msg: "视频删除成功"}, nil
}

func deleteVideoCaches(ctx context.Context, svcCtx *svc.ServiceContext, vid string) error {
	_ = svcCtx.RankCache.RemoveVideo(ctx, vid)
	_ = svcCtx.FavoriteCache.DelCount(ctx, vid)
	_ = svcCtx.CommentCache.DelCount(ctx, vid)
	_ = svcCtx.CommentCache.InvalidateListByVideo(ctx, vid)
	return nil
}
