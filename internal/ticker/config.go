package ticker

import (
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
)

func init() {
	config.AddReloader(reloadConfig)
}

func reloadConfig() {
	if renderTimer != nil {
		renderTimer.Reset(0)
	}
	nightscout.ClearUrl()
	Fetch()

	if fetchTimer != nil {
		fetchTimer.Reset(config.Default.Interval.Duration)
	}
}
