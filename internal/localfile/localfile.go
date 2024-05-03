package localfile

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
)

func New(conf *config.Config) *LocalFile {
	l := &LocalFile{
		config: conf,
	}
	l.reloadConfig()

	conf.AddCallback(l.reloadConfig)
	return l
}

type LocalFile struct {
	config *config.Config
	path   string
}

func (l *LocalFile) Format(last *nightscout.Properties) string {
	switch l.config.LocalFile.Format {
	case config.LocalFileFormatCsv:
		return last.Bgnow.DisplayBg(l.config.Units) + "," +
			last.Bgnow.Arrow(l.config.Arrows) + "," +
			last.Delta.Display(l.config.Units) + "," +
			strconv.Itoa(int(last.Bgnow.Mills.Time.Unix())) +
			"\n"
	default:
		return ""
	}
}

func (l *LocalFile) reloadConfig() {
	var path string
	if l.config.LocalFile.Enabled {
		path = l.config.LocalFile.Path
		if strings.HasPrefix(path, "$TMPDIR") {
			path = strings.TrimPrefix(path, "$TMPDIR")
			path = filepath.Join(os.TempDir(), path)
		}
	}
	if l.path != "" && path != l.path {
		if err := l.Cleanup(); err != nil {
			slog.Error("Failed to cleanup local file", "error", err.Error())
		}
	}
	l.path = path
}

func (l *LocalFile) Write(last *nightscout.Properties) error {
	if l.path != "" {
		segment := l.Format(last)
		err := os.WriteFile(l.path, []byte(segment), 0o600)
		return err
	}
	return nil
}

func (l *LocalFile) Cleanup() error {
	if l.path != "" {
		if err := os.Remove(l.path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}
