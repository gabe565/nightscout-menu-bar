package ticker

import (
	"time"

	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/gabe565/nightscout-menu-bar/internal/tray"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
)

var (
	renderTimer = time.NewTimer(5 * time.Minute)
	RenderCh    = make(chan *nightscout.Properties)
)

func BeginRender() {
	go func() {
		var properties *nightscout.Properties
		for {
			select {
			case p := <-RenderCh:
				properties = p
			case <-renderTimer.C:
			}
			if properties != nil {
				tray.Update <- properties
				renderTimer.Reset(util.GetNextMinChange(properties.Bgnow.Mills.Time))
			} else {
				renderTimer.Reset(5 * time.Minute)
			}
		}
	}()
}
