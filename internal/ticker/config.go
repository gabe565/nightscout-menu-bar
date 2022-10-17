package ticker

import (
	"github.com/spf13/viper"
)

func ReloadConfig() {
	timer.Reset(0)
	Fetch()

	if ticker != nil {
		ticker.Reset(viper.GetDuration("interval"))
	}
}
