package nightscout

import (
	"fmt"
	"github.com/gabe565/nightscout-systray/internal/util"
)

type Properties struct {
	Bgnow     Reading   `json:"bgnow"`
	Buckets   []Reading `json:"buckets"`
	Delta     Delta     `json:"delta"`
	Direction Direction `json:"direction"`
}

func (p Properties) String() string {
	return fmt.Sprintf(
		"%d %s %s [%s]",
		p.Bgnow.Last,
		p.Bgnow.Arrow(),
		p.Delta.Display,
		util.MinAgo(p.Bgnow.Time()),
	)
}
