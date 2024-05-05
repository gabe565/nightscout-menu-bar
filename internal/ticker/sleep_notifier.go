//go:build !darwin

package ticker

func (t *Ticker) beginSleepNotifier() {}
