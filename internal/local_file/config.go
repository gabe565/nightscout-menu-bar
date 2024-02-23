package local_file

import "github.com/spf13/viper"

const (
	EnabledKey = "local-file.enabled"
	PathKey    = "local-file.path"
	FormatKey  = "local-file.format"
	CleanupKey = "local-file.cleanup"
)

const FormatCsv = "csv"

func init() {
	viper.SetDefault(EnabledKey, false)
	viper.SetDefault(PathKey, "$TMPDIR/nightscout.csv")
	viper.SetDefault(FormatKey, FormatCsv)
	viper.SetDefault(CleanupKey, true)
}
