package ticker

import (
	"github.com/spf13/viper"
)

func ReloadConfig() {
	Fetch()

	if ticker != nil {
		ticker.Reset(viper.GetDuration("interval"))
	}
}
