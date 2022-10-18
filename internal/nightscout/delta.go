package nightscout

import (
	"fmt"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
	"github.com/spf13/viper"
	"math"
)

type Times struct {
	Previous Mills `json:"previous"`
	Recent   Mills `json:"recent"`
}

type Delta struct {
	Absolute     int     `json:"absolute"`
	DisplayVal   string  `json:"display"`
	ElapsedMins  float64 `json:"elapsedMins"`
	Interpolated bool    `json:"interpolated"`
	Mean5MinsAgo float64 `json:"mean5MinsAgo"`
	Mgdl         int     `json:"mgdl"`
	Previous     Reading `json:"previous"`
	Scaled       int     `json:"scaled"`
	Times        Times   `json:"times"`
}

func (d Delta) Display() string {
	if u := viper.GetString(config.UnitsKey); u == config.UnitsMmol {
		mmol := util.ToMmol(d.Scaled)
		mmol = math.Round(mmol*10) / 10
		return fmt.Sprintf("%+.1f", mmol)
	}

	return d.DisplayVal
}
