package trayhost

import "errors"

var (
	ErrCanceled   = errors.New("the AppIndicator extension wasn't installed")
	ErrNotEnabled = errors.New("the AppIndicator extension couldn't be enabled")
)
