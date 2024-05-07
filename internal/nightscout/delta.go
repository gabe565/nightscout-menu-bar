package nightscout

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
)

type Times struct {
	Previous Mills `json:"previous"`
	Recent   Mills `json:"recent"`
}

type Delta struct {
	Absolute     json.Number `json:"absolute"`
	DisplayVal   string      `json:"display"`
	ElapsedMins  json.Number `json:"elapsedMins"`
	Interpolated bool        `json:"interpolated"`
	Mean5MinsAgo json.Number `json:"mean5MinsAgo"`
	Mgdl         json.Number `json:"mgdl"`
	Previous     Reading     `json:"previous"`
	Scaled       int         `json:"scaled"`
	Times        Times       `json:"times"`
}

func (d Delta) Display(units string) string {
	if units == config.UnitsMmol {
		mmol := util.ToMmol(d.Scaled)
		mmol = math.Round(mmol*10) / 10
		f := strconv.FormatFloat(mmol, 'f', -1, 64)
		if mmol >= 0 {
			return "+" + f
		}
		return f
	}

	return d.DisplayVal
}
