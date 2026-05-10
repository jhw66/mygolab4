// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package totp

import (
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/logic/totp"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TotpDisableHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Disable2FAReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, types.CommonRsp{
				Status: 400,
				Msg:    "请求参数错误",
				Error:  err.Error(),
			})
			return
		}

		l := totp.NewTotpDisableLogic(r.Context(), svcCtx)
		resp, err := l.TotpDisable(&req)
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, types.CommonRsp{
				Status: 400,
				Msg:    "禁用2FA失败",
				Error:  err.Error(),
			})
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
