package nightscout

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/gabe565/nightscout-menu-bar/internal/util"
	"github.com/spf13/viper"
)

const (
	DoubleUpKey      = "arrows.double-up"
	SingleUpKey      = "arrows.single-up"
	FortyFiveUpKey   = "arrows.forty-five-up"
	FlatKey          = "arrows.flat"
	FortyFiveDownKey = "arrows.forty-five-down"
	SingleDownKey    = "arrows.single-down"
	DoubleDownKey    = "arrows.double-down"
	UnknownKey       = "arrows.unknown"
)

const (
	LowReading  = 39
	HighReading = 401
)

func init() {
	viper.SetDefault(DoubleUpKey, "⇈")
	viper.SetDefault(SingleUpKey, "↑")
	viper.SetDefault(FortyFiveUpKey, "↗")
	viper.SetDefault(FlatKey, "→")
	viper.SetDefault(FortyFiveDownKey, "↘")
	viper.SetDefault(SingleDownKey, "↓")
	viper.SetDefault(DoubleDownKey, "⇊")
	viper.SetDefault(UnknownKey, "-")
}

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
		direction = viper.GetString(DoubleUpKey)
	case "SingleUp":
		direction = viper.GetString(SingleUpKey)
	case "FortyFiveUp":
		direction = viper.GetString(FortyFiveUpKey)
	case "Flat":
		direction = viper.GetString(FlatKey)
	case "FortyFiveDown":
		direction = viper.GetString(FortyFiveDownKey)
	case "SingleDown":
		direction = viper.GetString(SingleDownKey)
	case "DoubleDown", "TripleDown":
		direction = viper.GetString(DoubleDownKey)
	default:
		direction = viper.GetString(UnknownKey)
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

	if u := viper.GetString(UnitsKey); u == UnitsMmol {
		mmol := util.ToMmol(r.Last)
		mmol = math.Round(mmol*10) / 10
		return strconv.FormatFloat(mmol, 'f', 1, 64)
	}

	return strconv.Itoa(r.Last)
}
