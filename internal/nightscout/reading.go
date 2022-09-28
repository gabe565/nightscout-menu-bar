package nightscout

import (
	"fmt"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
)

type Reading struct {
	Mean      int   `json:"mean"`
	Last      int   `json:"last"`
	Mills     Mills `json:"mills"`
	Index     int   `json:"index"`
	FromMills Mills `json:"fromMills"`
	ToMills   Mills `json:"toMills"`
	Sgvs      []SGV `json:"sgvs"`
}

func (r Reading) Arrow() string {
	direction := "-"
	if len(r.Sgvs) > 0 {
		direction = r.Sgvs[0].Direction
	}
	switch direction {
	case "FortyFiveUp":
		direction = "↗"
	case "FortyFiveDown":
		direction = "↘"
	case "SingleUp":
		direction = "↑"
	case "SingleDown":
		direction = "↓"
	case "Flat":
		direction = "→"
	case "DoubleUp":
		direction = "⇈"
	case "DoubleDown":
		direction = "⇊"
	}
	return direction
}

func (r Reading) String() string {
	return fmt.Sprintf(
		"%d %s [%s]",
		r.Last,
		r.Arrow(),
		util.MinAgo(r.Mills.Time),
	)
}
