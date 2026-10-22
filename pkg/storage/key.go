package storage

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Kind string

const (
	Video  Kind = "video"
	Cover  Kind = "cover"
	Avatar Kind = "avatar"
)

// validateStaticRoot only accepts the static directory under the current working directory. This function validates without creating any directories.
func validateStaticRoot(rootDir string) (string, error) {
	if strings.TrimSpace(rootDir) == "" {
		return "", errors.New("static root must not be empty")
	}
	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return "", fmt.Errorf("resolve static root: %w", err)
	}
	expectedRoot, err := filepath.Abs("static")
	if err != nil {
		return "", fmt.Errorf("resolve expected static root: %w", err)
	}
	rel, err := filepath.Rel(expectedRoot, absRoot)
	if err != nil || rel != "." {
		return "", fmt.Errorf("static root must be the working directory's static directory: %q", rootDir)
	}
	return expectedRoot, nil
}

// validateManagedPath checks the original relative path.
func validateManagedPath(relativePath string) error {
	kind, name, ok := strings.Cut(relativePath, "/")
	if !ok || !validKind(Kind(kind)) || !validFileName(name) {
		return fmt.Errorf("invalid managed file path: %q", relativePath)
	}
	return nil
}

// Disk paths may use native Windows separators, but parent segments must be rejected before Abs/Rel normalize them, even if they would stay inside static.
func validateDiskPathInput(diskPath string) error {
	if strings.TrimSpace(diskPath) == "" {
		return errors.New("disk path must not be empty")
	}
	for _, segment := range strings.FieldsFunc(diskPath, func(c rune) bool { return c == '/' || c == '\\' }) {
		if segment == ".." {
			return fmt.Errorf("disk path must not contain parent segments: %q", diskPath)
		}
	}
	return nil
}

func WebPath(key string) string {
	if key == "" {
		return ""
	}
	return "/static/" + key
}

func DiskPath(key string) string {
	if key == "" {
		return ""
	}
	return filepath.Join("static", filepath.FromSlash(key))
}

func validKind(kind Kind) bool {
	return kind == Video || kind == Cover || kind == Avatar
}

func validFileName(name string) bool {
	if name == "" || name == "." || name == ".." || len(name) > 255 || !utf8.ValidString(name) || strings.TrimRight(name, " .") != name || strings.ContainsAny(name, "/\\:<>\"|?*#%") {
		return false
	}
	if strings.IndexFunc(name, unicode.IsControl) >= 0 {
		return false
	}
	base, _, _ := strings.Cut(strings.ToUpper(name), ".")
	base = strings.TrimRight(base, " .")
	switch base {
	case "CON", "PRN", "AUX", "NUL", "CONIN$", "CONOUT$":
		return false
	}
	for _, prefix := range []string{"COM", "LPT"} {
		if strings.HasPrefix(base, prefix) && strings.Contains("123456789¹²³", strings.TrimPrefix(base, prefix)) &&
			utf8.RuneCountInString(strings.TrimPrefix(base, prefix)) == 1 {
			return false
		}
	}
	return true
}
