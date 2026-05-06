package utils

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ParseMultipartForm(r *http.Request, maxMemory int64) error {
	if r.MultipartForm != nil {
		return nil
	}
	return r.ParseMultipartForm(maxMemory)
}

func GetFormValue(r *http.Request, key string) (string, bool) {
	if r.MultipartForm == nil {
		return "", false
	}
	values, exists := r.MultipartForm.Value[key]
	if !exists {
		return "", false
	}
	if len(values) == 0 {
		return "", true
	}
	return values[0], true
}

func GetFormFile(r *http.Request, name string) (*multipart.FileHeader, error) {
	if r.MultipartForm == nil {
		return nil, http.ErrMissingFile
	}
	files, exists := r.MultipartForm.File[name]
	if !exists || len(files) == 0 {
		return nil, http.ErrMissingFile
	}
	return files[0], nil
}

func OptionalGetFormFile(r *http.Request, name string) (*multipart.FileHeader, bool) {
	header, err := GetFormFile(r, name)
	if err != nil {
		return nil, false
	}
	return header, true
}

func StoreUploadedFile(dir, prefix, ownerID string, file *multipart.FileHeader) (diskPath, webPath string, err error) {
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", "", err
	}
	diskPath = buildStoredFilePath(dir, prefix, ownerID, file.Filename)
	if err = saveMultipartFile(file, diskPath); err != nil {
		_ = RemoveIfExists(diskPath)
		return "", "", err
	}
	return diskPath, toWebPath(diskPath), nil
}

func ReplaceStoredFile(dir, prefix, ownerID string, file *multipart.FileHeader, oldWebPath string) (newDiskPath, newWebPath string, cleanupOld func() error, err error) {
	newDiskPath, newWebPath, err = StoreUploadedFile(dir, prefix, ownerID, file)
	if err != nil {
		return "", "", nil, err
	}
	cleanupOld = func() error {
		return RemoveIfExists(toDiskPath(oldWebPath))
	}
	return newDiskPath, newWebPath, cleanupOld, nil
}

func buildStoredFilePath(dir, prefix, ownerID, originalFilename string) string {
	filename := prefix + "_" + ownerID + "_" + strconvUnix() + filepath.Ext(originalFilename)
	return filepath.Join(dir, filename)
}

func saveMultipartFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func toWebPath(path string) string {
	if path == "" {
		return ""
	}
	webpath := filepath.ToSlash(path)
	if strings.HasPrefix(webpath, "/") {
		return webpath
	}
	return "/" + webpath
}

func toDiskPath(webpath string) string {
	if webpath == "" {
		return ""
	}
	diskpath := strings.TrimPrefix(webpath, "/")
	return filepath.FromSlash(diskpath)
}

func RemoveIfExists(diskpath string) error {
	if diskpath == "" {
		return nil
	}
	if err := os.Remove(diskpath); err != nil {
		return err
	}
	return nil
}

func strconvUnix() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
