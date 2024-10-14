package messages

import "gabe565.com/nightscout-menu-bar/internal/nightscout"

type ReloadConfigMsg struct{}

type RenderType uint8

const (
	RenderTypeFetch RenderType = iota
	RenderTypeTimestamp
)

type RenderMessage struct {
	Type       RenderType
	Properties *nightscout.Properties
}
