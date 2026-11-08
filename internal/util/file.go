package util

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// WriteFileAtomic writes a file via a temporary file and rename,
// so a crash can't leave it partially written.
func WriteFileAtomic(path string, data []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o777); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(tmp.Name())
	}()
	if err := tmp.Chmod(perm); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// WatchFile calls onChange after the file is created, written, or removed,
// until ctx is canceled. Events are debounced.
func WatchFile(ctx context.Context, path string, onChange func()) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o777); err != nil {
		return err
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	// Files written atomically are replaced, so the directory must be watched
	if err := w.Add(dir); err != nil {
		_ = w.Close()
		return err
	}

	go func() {
		defer func() {
			_ = w.Close()
		}()

		var debounce <-chan time.Time
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-w.Events:
				if !ok {
					return
				}
				if filepath.Base(event.Name) == filepath.Base(path) {
					debounce = time.After(100 * time.Millisecond)
				}
			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				slog.Error("File watcher failed", "path", path, "error", err)
			case <-debounce:
				onChange()
			}
		}
	}()
	return nil
}
