package local_file

import "github.com/spf13/viper"

const FormatCsv = "csv"

func init() {
	viper.SetDefault("local-file.enabled", false)
	viper.SetDefault("local-file.path", "$TMPDIR/nightscout.csv")
	viper.SetDefault("local-file.format", FormatCsv)
	viper.SetDefault("local-file.cleanup", true)
}
