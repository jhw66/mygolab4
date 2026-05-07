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
		// 删除视频的点赞
		if err := l.svcCtx.FavoriteRepo.DeleteByVideoIDWithTx(l.ctx, tx, req.Id); err != nil {
			l.Errorf("delete favorite by video failed, vid=%s, err=%v", req.Id, err)
			return err
		}
		// 删除视频的评论
		if err := l.svcCtx.CommentRepo.DeleteByVideoIDWithTx(l.ctx, tx, req.Id); err != nil {
			l.Errorf("delete comment by video failed, vid=%s, err=%v", req.Id, err)
			return err
		}
		// 删除视频
		deleted, err := l.svcCtx.VideoRepo.DeleteByIDWithTx(l.ctx, tx, req.Id)
		if err != nil {
			l.Errorf("delete video failed, vid=%s, err=%v", req.Id, err)
			return err
		}
		// 删除视频文件
		if err := utils.RemoveIfExistsWithUrl(deleted.URL); err != nil {
			l.Errorf("remove video file failed, url=%s, err=%v", deleted.URL, err)
			return err
		}
		// 删除视频封面文件
		if err := utils.RemoveIfExistsWithUrl(deleted.Cover); err != nil {
			l.Errorf("remove video cover file failed, url=%s, err=%v", deleted.Cover, err)
			return err
		}
		return nil
	})
	if err != nil {
		return &types.CommonRsp{Status: 500, Msg: "视频删除失败"}, errors.New("视频删除失败")
	}

	if err := deleteVideoCaches(l, req.Id); err != nil {
		l.Errorf("delete video caches failed, vid=%s, err=%v", req.Id, err)
	}

	return &types.CommonRsp{Status: 200, Msg: "视频删除成功"}, nil
}

func deleteVideoCaches(l *DeleteVideoLogic, vid string) error {
	if err := l.svcCtx.RankCache.RemoveVideo(l.ctx, vid); err != nil {
		l.Errorf("remove video from rank cache failed, vid=%s, err=%v", vid, err)
		return err
	}
	if err := l.svcCtx.FavoriteCache.DelCount(l.ctx, vid); err != nil {
		l.Errorf("delete favorite count failed, vid=%s, err=%v", vid, err)
		return err
	}
	if err := l.svcCtx.CommentCache.DelCount(l.ctx, vid); err != nil {
		l.Errorf("delete comment count failed, vid=%s, err=%v", vid, err)
		return err
	}
	if err := l.svcCtx.CommentCache.InvalidateListByVideo(l.ctx, vid); err != nil {
		l.Errorf("invalidate comment list failed, vid=%s, err=%v", vid, err)
		return err
	}
	return nil
}
