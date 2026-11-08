package systheme

import "context"

// Light is false since macOS recolors template icons for the menu bar itself.
func Light() bool { return false }

// MenuLight is true, but unused since macOS recolors template icons in menus too.
func MenuLight() bool { return true }

func Watch(context.Context, func()) error { return nil }
