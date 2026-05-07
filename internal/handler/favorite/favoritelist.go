// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package favorite

import (
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/logic/favorite"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FavoriteListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := favorite.NewFavoriteListLogic(r.Context(), svcCtx)
		resp, err := l.FavoriteList()
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, resp)
			return
		}
		httpx.WriteJsonCtx(r.Context(), w, 200, resp)
	}
}
