package storage

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestValidateStaticRoot(t *testing.T) {
	workDir := t.TempDir()
	t.Chdir(workDir)
	want := filepath.Join(workDir, "static")
	valid := []string{"static", "./static", want}
	if runtime.GOOS == "windows" {
		valid = append(valid, strings.ToUpper(want))
	}
	for _, input := range valid {
		got, err := validateStaticRoot(input)
		if err != nil || got != want {
			t.Errorf("root %q: got %q, want %q, err=%v", input, got, want, err)
		}
	}
	for _, input := range []string{"", " ", ".", "other", "static/video", filepath.Join(workDir, "other", "static")} {
		if _, err := validateStaticRoot(input); err == nil {
			t.Errorf("accepted incorrect root %q", input)
		}
	}
	if _, err := os.Stat(want); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("validation created static: %v", err)
	}
}

func TestValidateManagedPath(t *testing.T) {
	for _, input := range []string{
		"video/a.mp4", "cover/a.jpg", "avatar/a.png", "avatar/旧头像.png",
		"video/video_123_1770000000.mp4", "video/a..mp4",
	} {
		if err := validateManagedPath(input); err != nil {
			t.Errorf("rejected managed path %q: %v", input, err)
		}
	}
	for _, input := range []string{
		"", ".", "..", "video", "video/", "private/keep.txt", "/video/a.mp4",
		"video/../private/keep.txt", "video/../cover/a.jpg", "video/./a.mp4",
		"video/sub/a.mp4", "video//a.mp4", "video/..", "video/.",
		`video\a.mp4`, `video/..\cover\a.jpg`, `C:\static\video\a.mp4`,
		"video/a.mp4?x=1", "video/a#b.mp4", "video/a%2fb.mp4", "video/a:b.mp4",
		"video/a\x00.mp4", "video/a.mp4.", "video/a.mp4 ", "video/CON.mp4", "avatar/NUL",
	} {
		t.Run(input, func(t *testing.T) {
			if err := validateManagedPath(input); err == nil {
				t.Fatalf("accepted invalid path %q", input)
			}
		})
	}
}
