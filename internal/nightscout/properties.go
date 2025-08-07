package nightscout

import (
	"gabe565.com/nightscout-menu-bar/internal/config"
)

type Properties struct {
	Bgnow     Reading   `json:"bgnow"`
	Buckets   []Reading `json:"buckets"`
	Delta     Delta     `json:"delta"`
	Direction Direction `json:"direction"`
}

func (p Properties) String(data config.Data) string {
	result := p.Bgnow.DisplayBg(data.Units) +
		" " + p.Bgnow.Arrow(data.Arrows)
	if delta := p.Delta.Display(data.Units); delta != "" {
		result += " " + p.Delta.Display(data.Units)
	}
	if rel := p.Bgnow.Mills.Relative(data.Advanced.RoundAge); rel != "" {
		result += " [" + p.Bgnow.Mills.Relative(data.Advanced.RoundAge) + "]"
	}
	return result
}
