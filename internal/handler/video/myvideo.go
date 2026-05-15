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

func MyVideoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := video.NewMyVideoLogic(r.Context(), svcCtx)
		resp, err := l.MyVideo()
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, &types.CommonRsp{
				Status: 400,
				Msg:    "请求处理失败",
				Error:  err.Error(),
			})
			return
		}
		httpx.WriteJsonCtx(r.Context(), w, 200, resp)
	}
}
