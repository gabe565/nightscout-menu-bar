package ticker

import (
	"time"

	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
	"github.com/rs/zerolog/log"
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
				d := util.GetNextMinChange(properties.Bgnow.Mills.Time, t.config.Advanced.RoundAge)
				t.renderTicker.Reset(d)
				log.Debug().Stringer("in", d).Msg("Scheduled next render")
			} else {
				t.renderTicker.Reset(5 * time.Minute)
			}
		}
	}()
	return renderCh
}
