package nightscout

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"

	"gabe565.com/nightscout-menu-bar/internal/config"
)

const (
	LowReading  = 39
	HighReading = 401
)

type Reading struct {
	Mean      json.Number `json:"mean"`
	Last      Mgdl        `json:"last"`
	Mills     Mills       `json:"mills"`
	Index     json.Number `json:"index,omitempty"`
	FromMills Mills       `json:"fromMills"`
	ToMills   Mills       `json:"toMills"`
	Sgvs      []SGV       `json:"sgvs"`
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

func (r *Reading) String(data config.Data) string {
	if r.Last == 0 {
		return ""
	}

	var result strings.Builder

	result.WriteString(r.DisplayBg(data.Units))
	result.WriteRune(' ')
	result.WriteString(r.Arrow(data.Arrows))

	if rel := r.Mills.Relative(data.Advanced.RoundAge); rel != "" {
		result.WriteString(" [")
		result.WriteString(rel)
		result.WriteRune(']')
	}

	return result.String()
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

func (r *Reading) DisplayBg(units config.Unit) string {
	switch r.Last {
	case LowReading:
		return "LOW"
	case HighReading:
		return "HIGH"
	}

	if units == config.UnitMmol {
		mmol := r.Last.Mmol()
		mmol = math.Round(mmol*10) / 10
		return strconv.FormatFloat(mmol, 'f', 1, 64)
	}

	return strconv.Itoa(r.Last.Mgdl())
}
