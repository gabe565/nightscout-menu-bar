package nightscout

import (
	"fmt"
	"github.com/gabe565/nightscout-systray/internal/util"
	"time"
)

type Reading struct {
	Mean      int   `json:"mean"`
	Last      int   `json:"last"`
	Mills     int   `json:"mills"`
	Index     int   `json:"index"`
	FromMills int   `json:"fromMills"`
	ToMills   int   `json:"toMills"`
	Sgvs      []SGV `json:"sgvs"`
}

func (r Reading) Time() time.Time {
	return time.Unix(int64(r.Mills/1000), 0)
}

func (r Reading) Arrow() string {
	direction := r.Sgvs[0].Direction
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
		util.MinAgo(r.Time()),
	)
}
