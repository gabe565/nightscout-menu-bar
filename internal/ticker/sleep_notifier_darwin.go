package ticker

import (
	"log/slog"
	"time"

	"github.com/prashantgupta24/mac-sleep-notifier/notifier"
)

func (t *Ticker) beginSleepNotifier() {
	go func() {
		notify := &notifier.Notifier{}
		notifyCh := notify.Start()
		defer notify.Quit()

		for {
			select {
			case <-t.ctx.Done():
				return
			case activity := <-notifyCh:
				switch activity.Type {
				case notifier.Awake:
					slog.Info("Starting timers for awake mode")
					t.fetch.Reset()
					t.renderTicker.Reset(time.Second)
					t.fetchTicker.Reset(time.Second)
				case notifier.Sleep:
					slog.Info("Stopping timers for sleep mode")
					t.fetchTicker.Stop()
					t.renderTicker.Stop()
				}
			}
		}
	}()
}
