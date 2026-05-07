// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"net/http"

	"github.com/jhw66/myvideo_lab4/internal/logic/common"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CommentListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CommentListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, &types.CommentListRsp{
				Status: 400,
				Msg:    "请求参数错误",
				Error:  err.Error(),
			})
			return
		}

		l := common.NewCommentListLogic(r.Context(), svcCtx)
		resp, err := l.CommentList(&req)
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, resp)
			return
		}
		httpx.WriteJsonCtx(r.Context(), w, 200, resp)
	}
}
