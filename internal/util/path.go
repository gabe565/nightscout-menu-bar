package util

import (
	"os"
	"path/filepath"
	"strings"
)

func ResolvePath(path string) string {
	var prefixHome bool
	if after, ok := strings.CutPrefix(path, "$HOME"); ok {
		path = after
		prefixHome = true
	} else if after, ok := strings.CutPrefix(path, "~"); ok {
		path = after
		prefixHome = true
	}

	if prefixHome {
		if home, err := os.UserHomeDir(); err == nil {
			path = filepath.Join(home, path)
		}
	}

	if after, ok := strings.CutPrefix(path, "$TMPDIR"); ok {
		path = after
		path = filepath.Join(os.TempDir(), path)
	}

	return path
}
