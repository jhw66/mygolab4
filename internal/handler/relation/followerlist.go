// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package relation

import (
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/logic/relation"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FollowerListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := relation.NewFollowerListLogic(r.Context(), svcCtx)
		resp, err := l.FollowerList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
