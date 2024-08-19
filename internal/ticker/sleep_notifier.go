//go:build !darwin

package ticker

import "context"

func (t *Ticker) beginSleepNotifier(_ context.Context) {}
