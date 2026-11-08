//go:build !windows

package items

// UpdateIcons does nothing, since macOS adapts template icons itself.
func (i Items) UpdateIcons() {}
