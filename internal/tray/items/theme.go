//go:build !darwin

package items

import (
	"gabe565.com/nightscout-menu-bar/internal/assets"
	"gabe565.com/nightscout-menu-bar/internal/systheme"
)

// UpdateIcons matches the menu icons to the theme.
func (i Items) UpdateIcons() {
	pick := func(black, white []byte) []byte {
		if systheme.MenuLight() {
			return black
		}
		return white
	}
	i.LastReading.SetIcon(pick(assets.Reading, assets.ReadingWhite))
	i.Error.SetIcon(pick(assets.Error, assets.ErrorWhite))
	if !i.Error.visible {
		i.Error.MenuItem.Hide()
	}
	i.History.SetIcon(pick(assets.History, assets.HistoryWhite))
	i.OpenNightscout.SetIcon(pick(assets.Open, assets.OpenWhite))
	i.Settings.SetIcon(pick(assets.Preferences, assets.PreferencesWhite))
	i.Update.SetIcon(pick(assets.Update, assets.UpdateWhite))
	if !i.Update.visible {
		i.Update.MenuItem.Hide()
	}
	i.About.SetIcon(pick(assets.About, assets.AboutWhite))
	i.Quit.SetIcon(pick(assets.Quit, assets.QuitWhite))
}
