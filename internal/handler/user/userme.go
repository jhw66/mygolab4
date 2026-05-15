// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/logic/user"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserMeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserMeLogic(r.Context(), svcCtx)
		resp, err := l.UserMe()
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, &types.CommonRsp{
				Status: 500,
				Msg:    "查询用户信息失败",
				Error:  err.Error(),
			})
			return
		}
		httpx.WriteJsonCtx(r.Context(), w, 200, resp)
	}
}
