//go:build !windows

package systheme

import "context"

// Light is false since the macOS menu bar adapts template icons itself,
// and Linux panels don't reliably report their theme.
func Light() bool { return false }

func Watch(context.Context, func()) error { return nil }
