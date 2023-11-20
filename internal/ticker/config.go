package ticker

import (
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/spf13/viper"
)

func ReloadConfig() {
	if renderTimer != nil {
		renderTimer.Reset(0)
	}
	nightscout.ClearUrl()
	Fetch()

	if fetchTimer != nil {
		fetchTimer.Reset(viper.GetDuration("interval"))
	}
}
