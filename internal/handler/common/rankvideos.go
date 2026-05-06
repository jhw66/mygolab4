// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/logic/common"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RankVideosHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := common.NewRankVideosLogic(r.Context(), svcCtx)
		resp, err := l.RankVideos()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
