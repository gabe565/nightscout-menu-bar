package tray

import (
	"context"
	"log/slog"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/tray/messages"
	"gabe565.com/nightscout-menu-bar/internal/update"
)

const (
	updateInterval = 24 * time.Hour
	updateRetry    = time.Hour
)

func (t *Tray) checkUpdates(ctx context.Context) {
	checker := update.NewChecker(t.config.Version)
	if !checker.Enabled() {
		return
	}

	ticker := time.NewTicker(updateRetry)
	defer ticker.Stop()

	var lastCheck time.Time
	var latest *update.Release
	for {
		if t.config.Data().Updates.Check {
			if time.Since(lastCheck) >= updateInterval {
				rel, err := checker.Check(ctx)
				if err != nil {
					slog.Warn("Failed to check for updates", "error", err)
				} else {
					lastCheck = time.Now()
					latest = rel
				}
			}
			if latest != nil {
				select {
				case t.bus <- messages.UpdateAvailableMsg{Release: latest}:
				case <-ctx.Done():
					return
				}
			}
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		case <-t.checkUpdate:
		}
	}
}
