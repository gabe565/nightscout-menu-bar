package ticker

import (
	"context"
	"log/slog"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/nightscout"
	"gabe565.com/nightscout-menu-bar/internal/tray/messages"
	"gabe565.com/nightscout-menu-bar/internal/util"
)

func (t *Ticker) beginRender(ctx context.Context) chan<- *nightscout.Properties {
	renderCh := make(chan *nightscout.Properties)
	go func() {
		defer close(renderCh)
		t.renderTicker = time.NewTicker(5 * time.Minute)
		defer t.renderTicker.Stop()
		var properties *nightscout.Properties
		for {
			var renderType messages.RenderType
			select {
			case <-ctx.Done():
				return
			case p := <-renderCh:
				renderType = messages.RenderTypeFetch
				properties = p
			case <-t.renderTicker.C:
				renderType = messages.RenderTypeTimestamp
			}
			if properties != nil {
				t.bus <- messages.RenderMessage{
					Type:       renderType,
					Properties: properties,
				}
				d := util.GetNextMinChange(properties.Bgnow.Mills.Time, t.config.Data().Advanced.RoundAge)
				t.renderTicker.Reset(d)
				slog.Debug("Scheduled next render", "in", d)
			} else {
				t.renderTicker.Reset(5 * time.Minute)
			}
		}
	}()
	return renderCh
}
