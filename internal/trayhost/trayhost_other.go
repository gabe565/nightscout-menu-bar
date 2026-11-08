//go:build darwin || windows

package trayhost

import (
	"context"
	"errors"
)

// Available is always true, since macOS and Windows always have a tray.
func Available() bool { return true }

func IsGNOME() bool { return false }

func EnableExtension(context.Context) error { return errors.ErrUnsupported }

func Register(context.Context) error { return nil }
