package nightscout

import (
	"encoding/json"
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"github.com/gabe565/nightscout-menu-bar/internal/app/settings"
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
	FromMills Mills       `json:"fromMills,omitempty"`
	ToMills   Mills       `json:"toMills,omitempty"`
	Sgvs      []SGV       `json:"sgvs"`
}

func (r *Reading) Arrow() string {
	var direction string
	if len(r.Sgvs) > 0 {
		direction = r.Sgvs[0].Direction
	}
	switch direction {
	case "DoubleUp", "TripleUp":
		direction = "⇈"
	case "SingleUp":
		direction = "↑"
	case "FortyFiveUp":
		direction = "↗"
	case "Flat":
		direction = "→"
	case "FortyFiveDown":
		direction = "↘"
	case "SingleDown":
		direction = "↓"
	case "DoubleDown", "TripleDown":
		direction = "⇊"
	default:
		direction = "-"
	}
	return direction
}

func (r *Reading) String(prefs fyne.Preferences) string {
	if r.Last == 0 {
		return ""
	}

	var units settings.Unit
	if err := units.UnmarshalText([]byte(prefs.String(settings.UnitsKey))); err != nil {
		units = settings.UnitMgdl
	}

	result := r.DisplayBg(units) +
		" " + r.Arrow()
	if rel := r.Mills.Relative(); rel != "" {
		result += " [" + rel + "]"
	}
	return result
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

func (r *Reading) DisplayBg(units settings.Unit) string {
	switch r.Last {
	case LowReading:
		return "LOW"
	case HighReading:
		return "HIGH"
	}

	if units == settings.UnitMmol {
		mmol := r.Last.Mmol()
		mmol = math.Round(mmol*10) / 10
		return strconv.FormatFloat(mmol, 'f', 1, 64)
	}

	return strconv.Itoa(r.Last.Mgdl())
}
