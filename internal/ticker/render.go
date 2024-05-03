package ticker

import (
	"time"

	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
)

func (t *Ticker) beginRender() chan<- *nightscout.Properties {
	renderCh := make(chan *nightscout.Properties)
	go func() {
		defer close(renderCh)
		t.renderTicker = time.NewTicker(5 * time.Minute)
		var properties *nightscout.Properties
		for {
			select {
			case <-t.ctx.Done():
				return
			case p := <-renderCh:
				properties = p
			case <-t.renderTicker.C:
			}
			if properties != nil {
				t.bus <- properties
				t.renderTicker.Reset(util.GetNextMinChange(properties.Bgnow.Mills.Time))
			} else {
				t.renderTicker.Reset(5 * time.Minute)
			}
		}
	}()
	return renderCh
}
