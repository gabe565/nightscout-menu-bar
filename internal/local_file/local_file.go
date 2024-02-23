package local_file

import (
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
)

func ReloadConfig() {
	Enabled = viper.GetBool("local-file.enabled")
	format = viper.GetString("local-file.format")
	var newPath string
	if Enabled {
		newPath = viper.GetString("local-file.path")
		newPath = strings.ReplaceAll(newPath, "$TMPDIR/", os.TempDir())
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
