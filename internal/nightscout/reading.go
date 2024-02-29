package nightscout

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
)

const (
	LowReading  = 39
	HighReading = 401
)

type Reading struct {
	Mean      int   `json:"mean"`
	Last      int   `json:"last"`
	Mills     Mills `json:"mills"`
	Index     int   `json:"index,omitempty"`
	FromMills Mills `json:"fromMills,omitempty"`
	ToMills   Mills `json:"toMills,omitempty"`
	Sgvs      []SGV `json:"sgvs"`
}

func (r *Reading) Arrow() string {
	var direction string
	if len(r.Sgvs) > 0 {
		direction = r.Sgvs[0].Direction
	}
	switch direction {
	case "DoubleUp", "TripleUp":
		direction = config.Default.Arrows.DoubleUp
	case "SingleUp":
		direction = config.Default.Arrows.SingleUp
	case "FortyFiveUp":
		direction = config.Default.Arrows.FortyFiveUp
	case "Flat":
		direction = config.Default.Arrows.Flat
	case "FortyFiveDown":
		direction = config.Default.Arrows.FortyFiveDown
	case "SingleDown":
		direction = config.Default.Arrows.SingleDown
	case "DoubleDown", "TripleDown":
		direction = config.Default.Arrows.DoubleDown
	default:
		direction = config.Default.Arrows.Unknown
	}
	return direction
}

func (r *Reading) String() string {
	return fmt.Sprintf(
		"%s %s [%s]",
		r.DisplayBg(),
		r.Arrow(),
		util.MinAgo(r.Mills.Time),
	)
}

func (r *Reading) UnmarshalJSON(bytes []byte) error {
	type rawReading Reading
	if err := json.Unmarshal(bytes, (*rawReading)(r)); err != nil {
		return err
	}

	// Last is unset if reading is out of range.
	// Will be set from sgvs.
	if r.Last == 0 && len(r.Sgvs) > 0 {
		r.Last = r.Sgvs[0].Mgdl
		r.Mills = r.Sgvs[0].Mills
	}

	return nil
}

func (r *Reading) DisplayBg() string {
	switch r.Last {
	case LowReading:
		return "LOW"
	case HighReading:
		return "HIGH"
	}

	if config.Default.Units == config.UnitsMmol {
		mmol := util.ToMmol(r.Last)
		mmol = math.Round(mmol*10) / 10
		return strconv.FormatFloat(mmol, 'f', 1, 64)
	}

	return strconv.Itoa(r.Last)
}
