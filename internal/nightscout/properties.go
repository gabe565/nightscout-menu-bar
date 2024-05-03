package nightscout

import (
	"fmt"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
)

type Properties struct {
	Bgnow     Reading   `json:"bgnow"`
	Buckets   []Reading `json:"buckets"`
	Delta     Delta     `json:"delta"`
	Direction Direction `json:"direction"`
}

func (p Properties) String(units string, arrows config.Arrows) string {
	return fmt.Sprintf(
		"%s %s %s [%s]",
		p.Bgnow.DisplayBg(units),
		p.Bgnow.Arrow(arrows),
		p.Delta.Display(units),
		util.MinAgo(p.Bgnow.Mills.Time),
	)
}
