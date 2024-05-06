package nightscout

import (
	"encoding/json"
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

func (r *Reading) Arrow(conf config.Arrows) string {
	var direction string
	if len(r.Sgvs) > 0 {
		direction = r.Sgvs[0].Direction
	}
	switch direction {
	case "DoubleUp", "TripleUp":
		direction = conf.DoubleUp
	case "SingleUp":
		direction = conf.SingleUp
	case "FortyFiveUp":
		direction = conf.FortyFiveUp
	case "Flat":
		direction = conf.Flat
	case "FortyFiveDown":
		direction = conf.FortyFiveDown
	case "SingleDown":
		direction = conf.SingleDown
	case "DoubleDown", "TripleDown":
		direction = conf.DoubleDown
	default:
		direction = conf.Unknown
	}
	return direction
}

func (r *Reading) String(conf *config.Config) string {
	if r.Last == 0 {
		return ""
	}

	return r.DisplayBg(conf.Units) +
		" " + r.Arrow(conf.Arrows) +
		" [" + r.Mills.Relative(conf.Advanced.RoundAge) + "]"
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

func (r *Reading) DisplayBg(units string) string {
	switch r.Last {
	case LowReading:
		return "LOW"
	case HighReading:
		return "HIGH"
	}

	if units == config.UnitsMmol {
		mmol := util.ToMmol(r.Last)
		mmol = math.Round(mmol*10) / 10
		return strconv.FormatFloat(mmol, 'f', 1, 64)
	}

	return strconv.Itoa(r.Last)
}
