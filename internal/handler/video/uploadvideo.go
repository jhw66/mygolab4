// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video

import (
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/logic/video"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadVideoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := video.NewUploadVideoLogic(r.Context(), svcCtx)
		resp, err := l.UploadVideo(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
