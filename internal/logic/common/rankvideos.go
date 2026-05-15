// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"context"
	"errors"

	"github.com/jhw66/myvideo_lab4/internal/logic/core"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"github.com/jhw66/myvideo_lab4/pkg/serializer"

	"github.com/zeromicro/go-zero/core/logx"
)

type RankVideosLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRankVideosLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RankVideosLogic {
	return &RankVideosLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RankVideosLogic) RankVideos() (resp *types.VideoListRsp, err error) {
	vids, err := l.svcCtx.RankCache.GetTopVideoIDs(l.ctx, core.DefaultRankLimit)
	if err != nil || len(vids) == 0 {
		return nil, errors.New("获取热门排行榜失败")
	}
	videos, err := l.svcCtx.VideoRepo.FindByIDs(l.ctx, vids)
	if err != nil {
		return nil, errors.New("获取热门排行榜失败")
	}

	ordered := orderVideos(vids, videos)
	return serializer.VideoListRspFromModels(ordered), nil
}

func orderVideos(vids []string, videos []model.Video) []model.Video {
	// 将videos按id映射到map中
	videoMap := make(map[string]int, len(videos))
	for i := range videos {
		videoMap[videos[i].ID] = i
	}
	// 按照vids的顺序，将videos中的视频添加到ordered中
	ordered := make([]model.Video, 0, len(videos))
	for _, id := range vids {
		if idx, ok := videoMap[id]; ok {
			ordered = append(ordered, videos[idx])
		}
	}
	return ordered
}
