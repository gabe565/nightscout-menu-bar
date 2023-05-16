package ticker

import (
	"github.com/spf13/viper"
)

func ReloadConfig() {
	if renderTimer != nil {
		renderTimer.Reset(0)
	}
	Fetch()

	if fetchTimer != nil {
		fetchTimer.Reset(viper.GetDuration("interval"))
	}
}
