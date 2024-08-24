package util

import (
	"os"
	"path/filepath"
	"strings"
)

func ResolvePath(path string) string {
	var prefixHome bool
	if strings.HasPrefix(path, "$HOME") {
		path = strings.TrimPrefix(path, "$HOME")
		prefixHome = true
	} else if strings.HasPrefix(path, "~") {
		path = strings.TrimPrefix(path, "~")
		prefixHome = true
	}

	if prefixHome {
		if home, err := os.UserHomeDir(); err == nil {
			path = filepath.Join(home, path)
		}
	}

	if strings.HasPrefix(path, "$TMPDIR") {
		path = strings.TrimPrefix(path, "$TMPDIR")
		path = filepath.Join(os.TempDir(), path)
	}

	return path
}
