package nightscout

import (
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
)

type Properties struct {
	Bgnow     Reading   `json:"bgnow"`
	Buckets   []Reading `json:"buckets"`
	Delta     Delta     `json:"delta"`
	Direction Direction `json:"direction"`
}

func (p Properties) String(conf *config.Config) string {
	return p.Bgnow.DisplayBg(conf.Units) +
		" " + p.Bgnow.Arrow(conf.Arrows) +
		" " + p.Delta.Display(conf.Units) +
		" [" + util.MinAgo(p.Bgnow.Mills.Time, conf.Advanced.RoundAge) + "]"
}
