// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video

import (
	"context"
	"errors"
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/logic/core"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"github.com/jhw66/myvideo_lab4/pkg/serializer"
	"github.com/jhw66/myvideo_lab4/pkg/utils"
	"github.com/jhw66/myvideo_lab4/pkg/utlcontext"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadVideoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadVideoLogic {
	return &UploadVideoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadVideoLogic) UploadVideo(r *http.Request) (resp *types.VideoRsp, err error) {
	user, ok := utlcontext.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return &types.VideoRsp{Status: 401, Msg: "用户未登录"}, errors.New("用户未登录")
	}

	if err = utils.ParseMultipartForm(r, 64<<20); err != nil {
		return &types.VideoRsp{Status: 400, Msg: "请求格式错误"}, errors.New("请求格式错误")
	}

	title, exists := utils.GetFormValue(r, "title")
	if !exists || title == "" {
		return &types.VideoRsp{Status: 400, Msg: "标题不能为空"}, errors.New("标题不能为空")
	}
	info, _ := utils.GetFormValue(r, "info")

	videoFile, err := utils.GetFormFile(r, "video")
	if err != nil {
		return &types.VideoRsp{Status: 400, Msg: "视频文件不能为空"}, errors.New("视频文件不能为空")
	}
	coverFile, err := utils.GetFormFile(r, "cover")
	if err != nil {
		return &types.VideoRsp{Status: 400, Msg: "封面文件不能为空"}, errors.New("封面文件不能为空")
	}

	videoDiskPath, videoWebPath, err := utils.StoreUploadedFile("static/video", "video", user.ID, videoFile)
	if err != nil {
		return &types.VideoRsp{Status: 500, Msg: "保存视频文件失败"}, errors.New("保存视频文件失败")
	}
	coverDiskPath, coverWebPath, err := utils.StoreUploadedFile("static/cover", "cover", user.ID, coverFile)
	if err != nil {
		_ = utils.RemoveIfExists(videoDiskPath)
		return &types.VideoRsp{Status: 500, Msg: "保存封面文件失败"}, errors.New("保存封面文件失败")
	}

	video := model.Video{
		UserID:        user.ID,
		Title:         title,
		Info:          info,
		URL:           videoWebPath,
		Cover:         coverWebPath,
		FavoriteCount: 0,
		CommentCount:  0,
		HotScore:      uint64(core.CalculateHotScore(0, 0)),
	}

	err = l.svcCtx.VideoRepo.Create(l.ctx, &video)
	if err != nil {
		_ = utils.RemoveIfExists(videoDiskPath)
		_ = utils.RemoveIfExists(coverDiskPath)
		return &types.VideoRsp{Status: 500, Msg: "上传视频失败"}, errors.New("上传视频失败")
	}

	l.svcCtx.RankCache.ZAddScore(l.ctx, video.ID, float64(video.HotScore))

	return serializer.VideoRspFromModel(&video), nil
}
