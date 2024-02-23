package local_file

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/spf13/viper"
)

func Format(format string, last *nightscout.Properties) string {
	switch format {
	case FormatCsv:
		return last.Bgnow.DisplayBg() + "," +
			last.Bgnow.Arrow() + "," +
			last.Delta.Display() + "," +
			strconv.Itoa(int(last.Bgnow.Mills.Time.Unix())) +
			"\n"
	default:
		return ""
	}
}

var (
	Enabled bool
	path    string
	format  string
	cleanup bool
)

func ReloadConfig() {
	Enabled = viper.GetBool(EnabledKey)
	format = viper.GetString(FormatKey)
	cleanup = viper.GetBool(CleanupKey)
	var newPath string
	if Enabled {
		newPath = viper.GetString(PathKey)
		newPath = strings.ReplaceAll(newPath, "$TMPDIR/", os.TempDir())
	}
	if newPath != path {
		if err := Cleanup(); err != nil {
			log.Println(err)
		}
	}
	path = newPath
}

func Write(last *nightscout.Properties) error {
	if path == "" {
		ReloadConfig()
		if path == "" {
			return nil
		}
	}

	segment := Format(format, last)
	err := os.WriteFile(path, []byte(segment), 0o600)
	return err
}

func Cleanup() error {
	if cleanup && path != "" {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}
