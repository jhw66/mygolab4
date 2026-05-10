// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video

import (
	"context"
	"errors"
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"
	"github.com/jhw66/myvideo_lab4/pkg/serializer"
	"github.com/jhw66/myvideo_lab4/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateVideoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateVideoLogic {
	return &UpdateVideoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateVideoLogic) UpdateVideo(req *types.VideoIdReq, r *http.Request) (resp *types.VideoRsp, err error) {
	user, ok := auth.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return &types.VideoRsp{Status: 401, Msg: "用户未登录"}, errors.New("用户未登录")
	}
	if req.Id == "" {
		return &types.VideoRsp{Status: 400, Msg: "请传入视频id"}, errors.New("请传入视频id")
	}

	video, err := l.svcCtx.VideoRepo.FindByID(l.ctx, req.Id)
	if err != nil {
		return &types.VideoRsp{Status: 404, Msg: "未找到该视频"}, errors.New("未找到该视频")
	}
	if video.UserID != user.ID {
		return &types.VideoRsp{Status: 403, Msg: "没有修改视频权限或者不存在该视频"}, errors.New("没有修改视频权限或者不存在该视频")
	}

	if err = utils.ParseMultipartForm(r, 64<<20); err != nil {
		return &types.VideoRsp{Status: 400, Msg: "请求格式错误"}, errors.New("请求格式错误")
	}

	if title, exists := utils.GetFormValue(r, "title"); exists {
		if title == "" {
			return &types.VideoRsp{Status: 400, Msg: "标题不能为空"}, errors.New("标题不能为空")
		}
		video.Title = title
	}
	if info, exists := utils.GetFormValue(r, "info"); exists {
		video.Info = info
	}

	oldVideoURL := video.URL
	oldCoverURL := video.Cover
	newVideoDiskPath := ""
	newCoverDiskPath := ""
	var cleanupOldVideo func() error
	var cleanupOldCover func() error

	videoFile, hasVideo := utils.OptionalGetFormFile(r, "video")
	if hasVideo {
		var newVideoWebPath string
		newVideoDiskPath, newVideoWebPath, cleanupOldVideo, err = utils.ReplaceStoredFile("static/video", "video", user.ID, videoFile, oldVideoURL)
		if err != nil {
			return &types.VideoRsp{Status: 500, Msg: "保存视频文件失败"}, errors.New("保存视频文件失败")
		}
		video.URL = newVideoWebPath
	}

	coverFile, hasCover := utils.OptionalGetFormFile(r, "cover")
	if hasCover {
		var newCoverWebPath string
		newCoverDiskPath, newCoverWebPath, cleanupOldCover, err = utils.ReplaceStoredFile("static/cover", "cover", user.ID, coverFile, oldCoverURL)
		if err != nil {
			_ = utils.RemoveIfExists(newVideoDiskPath)
			return &types.VideoRsp{Status: 500, Msg: "保存封面文件失败"}, err
		}
		video.Cover = newCoverWebPath
	}

	err = l.svcCtx.VideoRepo.Update(l.ctx, video)
	if err != nil {
		_ = utils.RemoveIfExists(newVideoDiskPath)
		_ = utils.RemoveIfExists(newCoverDiskPath)
		return &types.VideoRsp{Status: 500, Msg: "更新视频失败"}, errors.New("更新视频失败")
	}
	if cleanupOldVideo != nil {
		_ = cleanupOldVideo()
	}
	if cleanupOldCover != nil {
		_ = cleanupOldCover()
	}

	l.svcCtx.RankCache.ZAddScore(l.ctx, video.ID, float64(video.HotScore))

	return serializer.VideoRspFromModel(video), nil
}
