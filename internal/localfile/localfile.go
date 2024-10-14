package localfile

import (
	"errors"
	"log/slog"
	"os"
	"strconv"

	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/nightscout"
	"gabe565.com/nightscout-menu-bar/internal/util"
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
		path = util.ResolvePath(l.config.LocalFile.Path)
	}
	if l.path != "" && path != l.path {
		if err := l.Cleanup(); err != nil {
			slog.Error("Failed to cleanup local file", "error", err)
		}
	}
	l.path = path
}

func (l *LocalFile) Write(last *nightscout.Properties) error {
	if l.path != "" {
		data := l.Format(last)
		slog.Debug("Writing local file", "data", data)
		err := os.WriteFile(l.path, []byte(data), 0o600)
		return err
	}
	return nil
}

func (l *LocalFile) Cleanup() error {
	if l.path != "" {
		slog.Debug("Removing local file", "path", l.path)
		if err := os.Remove(l.path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}
