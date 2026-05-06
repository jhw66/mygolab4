package handler

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func staticFileHandler(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name string `path:"name"`
		}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		name := filepath.Base(strings.TrimSpace(req.Name))
		if name == "" || name == "." || name == ".." {
			httpx.ErrorCtx(r.Context(), w, errors.New("文件名无效"))
			return
		}

		filePath := filepath.Join(baseDir, name)
		if info, err := os.Stat(filePath); err != nil || info.IsDir() {
			httpx.ErrorCtx(r.Context(), w, errors.New("文件不存在"))
			return
		}
		http.ServeFile(w, r, filePath)
	}
}
