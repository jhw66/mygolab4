// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package chat

import (
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/logic/chat"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AddMemberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddMemberReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, &types.CommonRsp{
				Status: 400,
				Msg:    "请求参数错误",
				Error:  err.Error(),
			})
			return
		}
		l := chat.NewAddMemberLogic(r.Context(), svcCtx)
		resp, err := l.AddMember(&req)
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, &types.CommonRsp{
				Status: 500,
				Msg:    "加入房间失败",
				Error:  err.Error(),
			})
			return
		}
		httpx.WriteJsonCtx(r.Context(), w, 200, resp)
	}
}
