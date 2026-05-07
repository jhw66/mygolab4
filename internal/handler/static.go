package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func staticFileHandler(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name string `path:"name"`
		}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.WriteJsonCtx(r.Context(), w, 200, &types.CommonRsp{
				Status: 400,
				Msg:    "请求参数错误",
				Error:  err.Error(),
			})
			return
		}

		name := filepath.Base(strings.TrimSpace(req.Name))
		if name == "" || name == "." || name == ".." {
			httpx.WriteJsonCtx(r.Context(), w, 200, &types.CommonRsp{
				Status: 400,
				Msg:    "文件名无效",
			})
			return
		}

		filePath := filepath.Join(baseDir, name)
		if info, err := os.Stat(filePath); err != nil || info.IsDir() {
			httpx.WriteJsonCtx(r.Context(), w, 200, &types.CommonRsp{
				Status: 404,
				Msg:    "文件不存在",
			})
			return
		}
		http.ServeFile(w, r, filePath)
	}
}
