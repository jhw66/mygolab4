package storage

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var ErrNotRegularFile = errors.New("managed path is not a regular file")

type StoredFile struct {
	DiskPath string
	WebPath  string
}

type LocalStore struct {
	root     *os.Root
	newToken func() (string, error)
}

func NewLocalStore(rootdir string) (*LocalStore, error) {
	absRoot, err := validateStaticRoot(rootdir)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(absRoot, 0755); err != nil {
		return nil, fmt.Errorf("create static root: %w", err)
	}
	root, err := os.OpenRoot(absRoot)
	if err != nil {
		return nil, fmt.Errorf("open static root: %w", err)
	}
	return &LocalStore{
		root: root, newToken: randomFileToken,
	}, nil
}

func (s *LocalStore) RootDir() string { return s.root.Name() }

func (s *LocalStore) Close() error { return s.root.Close() }

func (s *LocalStore) SaveMultipart(kind Kind, ownerID string, header *multipart.FileHeader) (StoredFile, error) {
	if !validKind(kind) {
		return StoredFile{}, fmt.Errorf("invalid upload kind")
	}
	if header == nil {
		return StoredFile{}, errors.New("upload file must not be nil")
	}
	if strings.TrimSpace(ownerID) == "" {
		return StoredFile{}, errors.New("upload owner must not be empty")
	}
	ext := path.Ext(strings.ReplaceAll(header.Filename, "\\", "/"))
	src, err := header.Open()
	if err != nil {
		return StoredFile{}, fmt.Errorf("open upload:%w", err)
	}
	defer src.Close()

	for attempt := 0; attempt < 8; attempt++ {
		token, err := s.newToken()
		if err != nil {
			return StoredFile{}, fmt.Errorf("generate file name: %w", err)
		}
		name := string(kind) + "_" + ownerID + "_" + token + ext
		filePath := string(kind) + "/" + name
		if err := validateManagedPath(filePath); err != nil {
			return StoredFile{}, err
		}
		if err := s.root.MkdirAll(string(kind), 0755); err != nil {
			return StoredFile{}, fmt.Errorf("create upload directory: %w", err)
		}
		out, err := s.root.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if errors.Is(err, fs.ErrExist) {
			continue
		}
		if err != nil {
			return StoredFile{}, fmt.Errorf("create upload: %w", err)
		}
		if err := s.writeCreatedUpload(filePath, out, src); err != nil {
			return StoredFile{}, err
		}
		return StoredFile{
			WebPath:  WebPath(filePath),
			DiskPath: DiskPath(filePath),
		}, nil
	}
	return StoredFile{}, errors.New("could not allocate a unique file name after 8 attempts")
}

func (s *LocalStore) writeCreatedUpload(relativePath string, out *os.File, src io.Reader) error {
	_, writeErr := io.Copy(out, src)
	if writeErr == nil {
		writeErr = out.Sync()
	}
	writeErr = errors.Join(writeErr, out.Close())
	if writeErr != nil {
		return fmt.Errorf("save upload: %w", errors.Join(writeErr, s.Remove(relativePath)))
	}
	return nil
}

func (s *LocalStore) Open(relativePath string) (*os.File, error) {
	if err := validateManagedPath(relativePath); err != nil {
		return nil, err
	}
	f, err := s.root.Open(relativePath)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	if !info.Mode().IsRegular() {
		f.Close()
		return nil, ErrNotRegularFile
	}
	return f, nil
}

func (s *LocalStore) Remove(relativePath string) error {
	if err := validateManagedPath(relativePath); err != nil {
		return err
	}
	info, err := s.root.Lstat(relativePath)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return ErrNotRegularFile
	}
	err = s.root.Remove(relativePath)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	return err
}

func (s *LocalStore) RootPathFromDiskPath(diskPath string) (string, error) {
	if err := validateDiskPathInput(diskPath); err != nil {
		return "", err
	}
	nativePath := filepath.FromSlash(strings.ReplaceAll(diskPath, "\\", "/"))
	absPath, err := filepath.Abs(nativePath)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(s.RootDir(), absPath)
	if err != nil {
		return "", fmt.Errorf("invalid managed file key: %v", err)
	}
	relativePath := filepath.ToSlash(rel)
	if err := validateManagedPath(relativePath); err != nil {
		return "", err
	}
	return relativePath, nil
}

func (s *LocalStore) RootPathFromWebPath(webPath string) (string, error) {
	if !strings.HasPrefix(webPath, "/static/") {
		return "", fmt.Errorf("invalid managed file key: expected /static/ path")
	}
	relativePath := strings.TrimPrefix(webPath, "/static/")
	if err := validateManagedPath(relativePath); err != nil {
		return "", err
	}
	return relativePath, nil
}

func randomFileToken() (string, error) {
	var token [16]byte
	if _, err := rand.Read(token[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(token[:]), nil
}
