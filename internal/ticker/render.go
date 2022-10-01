package ticker

import (
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/gabe565/nightscout-menu-bar/internal/tray"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
	"time"
)

var RenderCh = make(chan *nightscout.Properties)

func BeginRender() {
	go func() {
		timer := time.NewTimer(5 * time.Minute)
		var properties *nightscout.Properties
		for {
			select {
			case p := <-RenderCh:
				properties = p
			case <-timer.C:
			}
			if properties != nil {
				tray.Update <- properties
				timer.Reset(util.GetNextMinChange(properties.Bgnow.Mills.Time))
			} else {
				timer.Reset(5 * time.Minute)
			}
		}
	}()
}
