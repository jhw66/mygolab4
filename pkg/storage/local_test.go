package storage

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func testStore(t *testing.T) *LocalStore {
	t.Helper()
	t.Chdir(t.TempDir())
	store, err := NewLocalStore("static")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func uploadHeader(t *testing.T, name, content string) *multipart.FileHeader {
	t.Helper()
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, err := w.CreateFormFile("file", name)
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

func assertContent(t *testing.T, filename, want string) {
	t.Helper()
	got, err := os.ReadFile(filename)
	if err != nil || string(got) != want {
		t.Fatalf("file %q: content=%q, want=%q, err=%v", filename, got, want, err)
	}
}

func TestNewLocalStoreRejectsWrongRootsWithoutCreatingThem(t *testing.T) {
	workDir := t.TempDir()
	t.Chdir(workDir)
	for _, input := range []string{"other", "static/video", ".", filepath.Join(workDir, "other", "static")} {
		if store, err := NewLocalStore(input); err == nil {
			store.Close()
			t.Fatalf("accepted incorrect root %q", input)
		}
	}
	for _, input := range []string{"other", "static"} {
		if _, err := os.Stat(input); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("incorrect constructor created %q: %v", input, err)
		}
	}
	store, err := NewLocalStore("static")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	other, err := NewLocalStore(filepath.Join(workDir, "static"))
	if err != nil {
		t.Fatal(err)
	}
	defer other.Close()
	if other.RootDir() != store.RootDir() {
		t.Fatalf("absolute and relative constructors disagree: %q, %q", other.RootDir(), store.RootDir())
	}
}

func TestNormalUploadsAndPathConversion(t *testing.T) {
	store := testStore(t)
	for _, kind := range []Kind{Video, Cover, Avatar} {
		stored, err := store.SaveMultipart(kind, "123", uploadHeader(t, "sample.mp4", "content"))
		if err != nil {
			t.Fatalf("%s upload failed: %v", kind, err)
		}
		wantPrefix := string(kind) + "/" + string(kind) + "_123_"
		relativePath, err := store.RootPathFromWebPath(stored.WebPath)
		if err != nil || !strings.HasPrefix(relativePath, wantPrefix) {
			t.Fatalf("web path conversion failed: %q, %v", relativePath, err)
		}
		token := strings.TrimSuffix(strings.TrimPrefix(relativePath, wantPrefix), ".mp4")
		if len(token) != 32 || strings.ContainsAny(token, ":-/") {
			t.Fatalf("expected 128-bit random name, got %q", token)
		}
		if !filepath.IsAbs(stored.DiskPath) || stored.DiskPath != filepath.Join(store.RootDir(), filepath.FromSlash(relativePath)) {
			t.Fatalf("upload escaped its managed directory: %+v", stored)
		}
		for _, diskPath := range []string{stored.DiskPath, filepath.Join("static", filepath.FromSlash(relativePath)), strings.ReplaceAll(stored.DiskPath, "/", "\\")} {
			got, err := store.RootPathFromDiskPath(diskPath)
			if err != nil || got != relativePath {
				t.Fatalf("disk path conversion %q: got %q, want %q, err=%v", diskPath, got, relativePath, err)
			}
		}
		f, err := store.Open(relativePath)
		if err != nil {
			t.Fatal(err)
		}
		content, err := io.ReadAll(f)
		f.Close()
		if err != nil || string(content) != "content" {
			t.Fatalf("read uploaded content: %q, %v", content, err)
		}
	}
}

func TestConcurrentSameNameUploads(t *testing.T) {
	store := testStore(t)
	header := uploadHeader(t, "same.mp4", "complete content")
	const count = 32
	results := make(chan StoredFile, count)
	errs := make(chan error, count)
	var workers sync.WaitGroup
	for i := 0; i < count; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			stored, err := store.SaveMultipart(Video, "123", header)
			if err != nil {
				errs <- err
				return
			}
			results <- stored
		}()
	}
	workers.Wait()
	close(results)
	close(errs)
	for err := range errs {
		t.Error(err)
	}
	seen := make(map[string]bool)
	for result := range results {
		if seen[result.DiskPath] {
			t.Errorf("duplicate upload path: %q", result.DiskPath)
		}
		seen[result.DiskPath] = true
		assertContent(t, result.DiskPath, "complete content")
	}
	if len(seen) != count {
		t.Fatalf("got %d unique uploads, want %d", len(seen), count)
	}
}

func TestCollisionAndFailedUploadsPreserveExistingFiles(t *testing.T) {
	store := testStore(t)
	store.newToken = func() (string, error) { return strings.Repeat("a", 32), nil }
	original, err := store.SaveMultipart(Video, "123", uploadHeader(t, "same.mp4", "original"))
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	store.newToken = func() (string, error) {
		calls++
		if calls == 1 {
			return strings.Repeat("a", 32), nil
		}
		return strings.Repeat("b", 32), nil
	}
	replacement, err := store.SaveMultipart(Video, "123", uploadHeader(t, "same.mp4", "replacement"))
	if err != nil || replacement.DiskPath == original.DiskPath || calls != 2 {
		t.Fatalf("collision was not retried correctly: %+v, calls=%d, err=%v", replacement, calls, err)
	}
	store.newToken = func() (string, error) { return strings.Repeat("a", 32), nil }
	if _, err := store.SaveMultipart(Video, "123", uploadHeader(t, "same.mp4", "never saved")); err == nil {
		t.Fatal("eight collisions must return an error")
	}
	if _, err := store.SaveMultipart(Video, "123", nil); err == nil {
		t.Fatal("nil file must return an error")
	}
	if _, err := store.SaveMultipart(Video, "123", &multipart.FileHeader{Filename: "missing.mp4"}); err == nil {
		t.Fatal("missing upload source must return an error")
	}
	sourceErr := errors.New("random source unavailable")
	store.newToken = func() (string, error) { return "", sourceErr }
	if _, err := store.SaveMultipart(Video, "123", uploadHeader(t, "same.mp4", "never saved")); !errors.Is(err, sourceErr) {
		t.Fatalf("random source error was lost: %v", err)
	}
	assertContent(t, original.DiskPath, "original")
	assertContent(t, replacement.DiskPath, "replacement")
}

type failedReader struct{ err error }

func (r failedReader) Read([]byte) (int, error) { return 0, r.err }

func TestFailedWriteRemovesOnlyCreatedFile(t *testing.T) {
	store := testStore(t)
	if err := os.Mkdir(filepath.Join(store.RootDir(), "video"), 0755); err != nil {
		t.Fatal(err)
	}
	keep := filepath.Join(store.RootDir(), "video", "keep.mp4")
	if err := os.WriteFile(keep, []byte("keep"), 0644); err != nil {
		t.Fatal(err)
	}
	sourceErr := errors.New("source read failed")
	for _, closed := range []bool{false, true} {
		filePath := "video/partial.mp4"
		out, err := store.root.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			t.Fatal(err)
		}
		var src io.Reader = io.MultiReader(strings.NewReader("partial"), failedReader{sourceErr})
		if closed {
			out.Close()
			src = bytes.NewReader(nil) // Exercise synchronization/close failure too.
		}
		if err := store.writeCreatedUpload(filePath, out, src); err == nil {
			t.Fatal("failed write must return an error")
		}
		if _, err := os.Stat(filepath.Join(store.RootDir(), filepath.FromSlash(filePath))); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("incomplete file was not removed: %v", err)
		}
		assertContent(t, keep, "keep")
	}
}

func TestRemoveIsIdempotentAndRejectsDirectories(t *testing.T) {
	store := testStore(t)
	stored, err := store.SaveMultipart(Avatar, "123", uploadHeader(t, "photo.png", "avatar"))
	if err != nil {
		t.Fatal(err)
	}
	key, _ := store.RootPathFromWebPath(stored.WebPath)
	for i := 0; i < 2; i++ {
		if err := store.Remove(key); err != nil {
			t.Fatalf("remove attempt %d: %v", i, err)
		}
	}
	if _, err := store.Open(key); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("removed file still accessible: %v", err)
	}
	if err := os.Mkdir(filepath.Join(store.RootDir(), "avatar", "directory.png"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := store.Remove("avatar/directory.png"); !errors.Is(err, ErrNotRegularFile) {
		t.Fatalf("directory deletion was not rejected: %v", err)
	}
}

func TestInvalidPathsCannotDeleteOrOpenFiles(t *testing.T) {
	store := testStore(t)
	for _, dir := range []string{"private", "video", "cover"} {
		if err := os.Mkdir(filepath.Join(store.RootDir(), dir), 0755); err != nil {
			t.Fatal(err)
		}
	}
	private := filepath.Join(store.RootDir(), "private", "keep.txt")
	cover := filepath.Join(store.RootDir(), "cover", "keep.jpg")
	for _, target := range []string{private, cover} {
		if err := os.WriteFile(target, []byte("keep"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	for _, input := range []string{"", ".", "private/keep.txt", "video/../private/keep.txt", "video/../cover/keep.jpg", `video\..\cover\keep.jpg`, cover} {
		if err := store.Remove(input); err == nil {
			t.Errorf("Remove accepted invalid path %q", input)
		}
		if f, err := store.Open(input); err == nil {
			f.Close()
			t.Errorf("Open accepted invalid path %q", input)
		}
		assertContent(t, private, "keep")
		assertContent(t, cover, "keep")
	}
	for _, input := range []string{"/static/private/keep.txt", "/static/video/../private/keep.txt", "/static/video/../cover/keep.jpg", "/static/cover/keep.jpg?x=1", "https://example.com/static/cover/keep.jpg"} {
		if _, err := store.RootPathFromWebPath(input); err == nil {
			t.Errorf("accepted invalid web path %q", input)
		}
		assertContent(t, private, "keep")
		assertContent(t, cover, "keep")
	}
	for _, input := range []string{
		"", private, filepath.Dir(store.RootDir()),
		store.RootDir() + "/video/../cover/keep.jpg",
		store.RootDir() + `\video\..\cover\keep.jpg`,
		"static/video/../cover/keep.jpg", `static\video\..\cover\keep.jpg`,
	} {
		if _, err := store.RootPathFromDiskPath(input); err == nil {
			t.Errorf("accepted invalid disk path %q", input)
		}
		assertContent(t, private, "keep")
		assertContent(t, cover, "keep")
	}
}

func TestStoreRemainsPinnedAfterWorkingDirectoryChanges(t *testing.T) {
	store := testStore(t)
	originalRoot := store.RootDir()
	t.Chdir(t.TempDir())
	stored, err := store.SaveMultipart(Video, "123", uploadHeader(t, "video.mp4", "pinned"))
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(filepath.Dir(stored.DiskPath)) != originalRoot {
		t.Fatalf("store followed the new working directory: %+v", stored)
	}
	key, err := store.RootPathFromDiskPath(stored.DiskPath)
	if err != nil {
		t.Fatal(err)
	}
	f, err := store.Open(key)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	if err := store.Remove(key); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat("static"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("old instance created static under the new cwd: %v", err)
	}
}

func TestSymlinkDeletionIsRejected(t *testing.T) {
	store := testStore(t)
	stored, err := store.SaveMultipart(Video, "123", uploadHeader(t, "video.mp4", "keep"))
	if err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(store.RootDir(), "video", "link.mp4")
	if err := os.Symlink(stored.DiskPath, link); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	if err := store.Remove("video/link.mp4"); !errors.Is(err, ErrNotRegularFile) {
		t.Fatalf("symlink deletion was not rejected: %v", err)
	}
	assertContent(t, stored.DiskPath, "keep")
}
