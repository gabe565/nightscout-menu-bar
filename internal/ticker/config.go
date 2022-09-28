package ticker

import (
	"github.com/spf13/viper"
)

func ReloadConfig() {
	if ticker != nil {
		ticker.Reset(viper.GetDuration("interval"))
	}

	Tick()
}
