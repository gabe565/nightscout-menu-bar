package nightscout

import (
	"encoding/json"
	"fmt"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
	"github.com/spf13/viper"
	"strconv"
)

const (
	LowReading  = 39
	HighReading = 401
)

func init() {
	viper.SetDefault("arrows.double-up", "⇈")
	viper.SetDefault("arrows.single-up", "↑")
	viper.SetDefault("arrows.forty-five-up", "↗")
	viper.SetDefault("arrows.flat", "→")
	viper.SetDefault("arrows.forty-five-down", "↘")
	viper.SetDefault("arrows.single-down", "↓")
	viper.SetDefault("arrows.double-down", "⇊")
	viper.SetDefault("arrows.unknown", "-")
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
		direction = viper.GetString("arrows.double-up")
	case "SingleUp":
		direction = viper.GetString("arrows.single-up")
	case "FortyFiveUp":
		direction = viper.GetString("arrows.forty-five-up")
	case "Flat":
		direction = viper.GetString("arrows.flat")
	case "FortyFiveDown":
		direction = viper.GetString("arrows.forty-five-down")
	case "SingleDown":
		direction = viper.GetString("arrows.single-down")
	case "DoubleDown", "TripleDown":
		direction = viper.GetString("arrows.double-down")
	default:
		direction = viper.GetString("arrows.unknown")
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
	default:
		return strconv.Itoa(r.Last)
	}
}
