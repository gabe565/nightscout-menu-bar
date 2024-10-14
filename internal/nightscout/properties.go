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

func (p Properties) String(conf *config.Config) string {
	result := p.Bgnow.DisplayBg(conf.Units) +
		" " + p.Bgnow.Arrow(conf.Arrows)
	if delta := p.Delta.Display(conf.Units); delta != "" {
		result += " " + p.Delta.Display(conf.Units)
	}
	if rel := p.Bgnow.Mills.Relative(conf.Advanced.RoundAge); rel != "" {
		result += " [" + p.Bgnow.Mills.Relative(conf.Advanced.RoundAge) + "]"
	}
	return result
}
