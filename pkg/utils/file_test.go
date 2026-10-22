package utils

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func testUpload(t *testing.T, content string) *multipart.FileHeader {
	t.Helper()
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, err := w.CreateFormFile("file", "same.jpg")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(part, content); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("POST", "/", &body)
	r.Header.Set("Content-Type", w.FormDataContentType())
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { r.MultipartForm.RemoveAll() })
	return r.MultipartForm.File["file"][0]
}

func TestUploadRejectsWrongDirectoryWithoutCreatingIt(t *testing.T) {
	workDir := t.TempDir()
	t.Chdir(workDir)
	for _, dir := range []string{"other/video", filepath.Join(workDir, "other", "static", "video")} {
		if _, _, err := StoreUploadedFile(dir, "video", "123", testUpload(t, "content")); err == nil {
			t.Errorf("accepted incorrect directory %q", dir)
		}
	}
	if _, err := os.Stat("other"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("incorrect root was created: %v", err)
	}
	if _, err := os.Stat("static"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("rejected uploads created static: %v", err)
	}
}

func TestUploadReplacementAndRepeatedCleanup(t *testing.T) {
	t.Chdir(t.TempDir())
	oldDisk, oldWeb, err := StoreUploadedFile("static/avatar", "avatar", "123", testUpload(t, "old"))
	if err != nil {
		t.Fatal(err)
	}
	newDisk, newWeb, cleanup, err := ReplaceStoredFile("static/avatar", "avatar", "123", testUpload(t, "new"), oldWeb)
	if err != nil {
		t.Fatal(err)
	}
	if oldDisk == newDisk || oldWeb == newWeb {
		t.Fatal("same-name replacement reused the old path")
	}
	if content, err := os.ReadFile(oldDisk); err != nil || string(content) != "old" {
		t.Fatalf("old file changed before cleanup: %q, %v", content, err)
	}
	for i := 0; i < 2; i++ {
		if err := cleanup(); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Stat(oldDisk); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("old file was not removed: %v", err)
	}
	if content, err := os.ReadFile(newDisk); err != nil || string(content) != "new" {
		t.Fatalf("replacement was damaged: %q, %v", content, err)
	}
	if err := RemoveIfExists(newDisk); err != nil {
		t.Fatal(err)
	}
	if err := RemoveIfExists(newDisk); err != nil {
		t.Fatal(err)
	}
}

func TestDeletionWrappersPreserveInvalidTargets(t *testing.T) {
	t.Chdir(t.TempDir())
	for _, dir := range []string{"private", "video", "cover"} {
		if err := os.MkdirAll(filepath.Join("static", dir), 0755); err != nil {
			t.Fatal(err)
		}
	}
	private := filepath.Join("static", "private", "keep.txt")
	cover := filepath.Join("static", "cover", "keep.jpg")
	for _, filename := range []string{private, cover} {
		if err := os.WriteFile(filename, []byte("keep"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	for _, input := range []string{"/static/private/keep.txt", "/static/video/../private/keep.txt", "/static/video/../cover/keep.jpg"} {
		if err := RemoveIfExistsWithUrl(input); err == nil {
			t.Errorf("accepted invalid web path %q", input)
		}
	}
	for _, input := range []string{private, "static/video/../private/keep.txt", `static\video\..\cover\keep.jpg`} {
		if err := RemoveIfExists(input); err == nil {
			t.Errorf("accepted invalid disk path %q", input)
		}
	}
	for _, filename := range []string{private, cover} {
		if content, err := os.ReadFile(filename); err != nil || string(content) != "keep" {
			t.Fatalf("invalid deletion changed %q: %q, %v", filename, content, err)
		}
	}
}
