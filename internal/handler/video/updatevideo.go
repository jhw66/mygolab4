// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video

import (
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/logic/video"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateVideoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VideoIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, &types.VideoRsp{
				Status: 400,
				Msg:    "请求参数错误",
				Error:  err.Error(),
			})
			return
		}

		l := video.NewUpdateVideoLogic(r.Context(), svcCtx)
		resp, err := l.UpdateVideo(&req, r)
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, resp)
			return
		}
		httpx.WriteJsonCtx(r.Context(), w, 200, resp)
	}
}
