package messages

import (
	"gabe565.com/nightscout-menu-bar/internal/nightscout"
	"gabe565.com/nightscout-menu-bar/internal/update"
)

type ReloadConfigMsg struct{}

type ThemeChangedMsg struct{}

type OpenSettingsMsg struct{}

type UpdateAvailableMsg struct {
	Release *update.Release
}

type RenderType uint8

const (
	RenderTypeFetch RenderType = iota
	RenderTypeTimestamp
)

type RenderMessage struct {
	Type       RenderType
	Properties *nightscout.Properties
}
