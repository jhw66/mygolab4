package utils

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/jhw66/myvideo_lab4/pkg/storage"
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
	if strings.TrimSpace(dir) == "" || filepath.Base(filepath.Clean(dir)) != prefix {
		return "", "", fmt.Errorf("upload directory or kind is wrong")
	}
	store, err := storage.NewLocalStore(filepath.Dir(filepath.Clean(dir)))
	if err != nil {
		return "", "", err
	}
	defer store.Close()
	stored, err := store.SaveMultipart(storage.Kind(prefix), ownerID, file)
	if err != nil {
		return "", "", err
	}
	return stored.DiskPath, stored.WebPath, nil
}

func ReplaceStoredFile(dir, prefix, ownerID string, file *multipart.FileHeader, oldWebPath string) (newDiskPath, newWebPath string, cleanupOld func() error, err error) {
	newDiskPath, newWebPath, err = StoreUploadedFile(dir, prefix, ownerID, file)
	if err != nil {
		return "", "", nil, err
	}
	cleanupOld = func() error {
		return RemoveIfExistsWithUrl(oldWebPath)
	}
	return newDiskPath, newWebPath, cleanupOld, nil
}

func RemoveIfExists(diskpath string) error {
	if diskpath == "" {
		return nil
	}
	store, err := storage.NewLocalStore("static")
	if err != nil {
		return err
	}
	defer store.Close()
	path, err := store.RootPathFromDiskPath(diskpath)
	if err != nil {
		return err
	}
	return store.Remove(path)
}

func RemoveIfExistsWithUrl(url string) error {
	if url == "" {
		return nil
	}
	store, err := storage.NewLocalStore("static")
	if err != nil {
		return err
	}
	defer store.Close()
	path, err := store.RootPathFromWebPath(url)
	if err != nil {
		return err
	}
	return store.Remove(path)
}
