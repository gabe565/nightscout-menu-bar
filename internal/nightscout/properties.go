package nightscout

import (
	"strings"

	"gabe565.com/nightscout-menu-bar/internal/config"
)

type Properties struct {
	Bgnow     Reading   `json:"bgnow"`
	Buckets   []Reading `json:"buckets"`
	Delta     Delta     `json:"delta"`
	Direction Direction `json:"direction"`
}

func (p Properties) String(data config.Data) string {
	var result strings.Builder

	result.WriteString(p.Bgnow.DisplayBg(data.Units))
	if !data.LastReading.HideArrow {
		result.WriteRune(' ')
		result.WriteString(p.Bgnow.Arrow(data.Arrows))
	}

	if !data.LastReading.HideDelta {
		if delta := p.Delta.Display(data.Units); delta != "" {
			result.WriteRune(' ')
			result.WriteString(delta)
		}
	}

	if !data.LastReading.HideTimeAgo {
		if rel := p.Bgnow.Mills.Relative(data.Advanced.RoundAge); rel != "" {
			result.WriteString(" [")
			result.WriteString(rel)
			result.WriteRune(']')
		}
	}

	return result.String()
}
