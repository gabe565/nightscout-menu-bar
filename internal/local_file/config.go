package local_file

import (
	"path/filepath"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	EnabledFlag = "local-file"
	EnabledKey  = "local-file.enabled"

	PathFlag = "local-file-path"
	PathKey  = "local-file.path"

	FormatFlag = "local-file-format"
	FormatKey  = "local-file.format"

	CleanupFlag = "local-file-cleanup"
	CleanupKey  = "local-file.cleanup"
)

const FormatCsv = "csv"

func init() {
	flag.Bool(EnabledFlag, false, "Write blood sugar to a local file")
	if err := viper.BindPFlag(EnabledKey, flag.Lookup(EnabledFlag)); err != nil {
		panic(err)
	}

	defaultPath := filepath.Join("$TMPDIR", "nightscout.csv")
	flag.String(PathFlag, defaultPath, "Write blood sugar to a local file")
	if err := viper.BindPFlag(PathKey, flag.Lookup(PathFlag)); err != nil {
		panic(err)
	}

	flag.String(FormatFlag, FormatCsv, "Local file format (one of: csv)")
	if err := viper.BindPFlag(FormatKey, flag.Lookup(FormatFlag)); err != nil {
		panic(err)
	}

	flag.Bool(CleanupFlag, true, "Clean up local file on exit")
	if err := viper.BindPFlag(CleanupKey, flag.Lookup(CleanupFlag)); err != nil {
		panic(err)
	}
}
