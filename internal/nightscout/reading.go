package nightscout

import (
	"fmt"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("arrows.forty-five-up", "↗")
	viper.SetDefault("arrows.forty-five-down", "↘")
	viper.SetDefault("arrows.single-up", "↑")
	viper.SetDefault("arrows.single-down", "↓")
	viper.SetDefault("arrows.flat", "→")
	viper.SetDefault("arrows.double-up", "⇈")
	viper.SetDefault("arrows.double-down", "⇊")
	viper.SetDefault("arrows.unknown", "-")
}

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
	direction := ""
	if len(r.Sgvs) > 0 {
		direction = r.Sgvs[0].Direction
	}
	switch direction {
	case "FortyFiveUp":
		direction = viper.GetString("arrows.forty-five-up")
	case "FortyFiveDown":
		direction = viper.GetString("arrows.forty-five-down")
	case "SingleUp":
		direction = viper.GetString("arrows.single-up")
	case "SingleDown":
		direction = viper.GetString("arrows.single-down")
	case "Flat":
		direction = viper.GetString("arrows.flat")
	case "DoubleUp":
		direction = viper.GetString("arrows.double-up")
	case "DoubleDown":
		direction = viper.GetString("arrows.double-down")
	default:
		direction = viper.GetString("arrows.unknown")
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
