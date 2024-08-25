package ticker

import (
	"context"
	"log/slog"
	"time"

	"github.com/prashantgupta24/mac-sleep-notifier/notifier"
)

func (t *Ticker) beginSleepNotifier(ctx context.Context) {
	go func() {
		notify := &notifier.Notifier{}
		notifyCh := notify.Start()
		defer close(notifyCh)
		defer notify.Quit()

		for {
			select {
			case <-ctx.Done():
				return
			case activity := <-notifyCh:
				logger := slog.With("reason", activity.Type)
				switch activity.Type {
				case notifier.Awake:
					logger.Info("Starting timers")
					t.renderTicker.Reset(time.Second)
					t.fetchTicker.Reset(time.Second)
				case notifier.Sleep:
					logger.Info("Stopping timers")
					t.fetchTicker.Stop()
					t.renderTicker.Stop()
				}
			}
		}
	}()
}
