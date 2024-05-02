package localfile

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
)

func Format(format string, last *nightscout.Properties) string {
	switch format {
	case config.LocalFileFormatCsv:
		return last.Bgnow.DisplayBg() + "," +
			last.Bgnow.Arrow() + "," +
			last.Delta.Display() + "," +
			strconv.Itoa(int(last.Bgnow.Mills.Time.Unix())) +
			"\n"
	default:
		return ""
	}
}

var path string

func ReloadConfig() {
	var newPath string
	if config.Default.LocalFile.Enabled {
		newPath = config.Default.LocalFile.Path
		if strings.HasPrefix(newPath, "$TMPDIR") {
			newPath = strings.Replace(newPath, "$TMPDIR"+string(os.PathSeparator), os.TempDir(), 1)
		}
	}
	if newPath != path {
		if err := Cleanup(); err != nil {
			slog.Error("Failed to cleanup local file", "error", err.Error())
		}
	}
	path = newPath
}

func init() {
	config.AddReloader(ReloadConfig)
}

func Write(last *nightscout.Properties) error {
	if path == "" {
		ReloadConfig()
		if path == "" {
			return nil
		}
	}

	segment := Format(config.Default.LocalFile.Format, last)
	err := os.WriteFile(path, []byte(segment), 0o600)
	return err
}

func Cleanup() error {
	if config.Default.LocalFile.Cleanup && path != "" {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}
