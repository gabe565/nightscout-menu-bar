package nightscout

import (
	"encoding/json"
	"math"
	"strconv"

	"gabe565.com/nightscout-menu-bar/internal/config"
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
	Mgdl         Mgdl        `json:"mgdl"`
	Previous     Reading     `json:"previous"`
	Scaled       json.Number `json:"scaled"`
	Times        Times       `json:"times"`
}

func (d Delta) Display(units config.Unit) string {
	if units == config.UnitMmol {
		mmol := d.Mgdl.Mmol()
		mmol = math.Round(mmol*10) / 10
		f := strconv.FormatFloat(mmol, 'f', -1, 64)
		if mmol >= 0 {
			return "+" + f
		}
		return f
	}

	mgdl := d.Mgdl.Mgdl()
	val := strconv.Itoa(mgdl)
	if mgdl >= 0 {
		return "+" + val
	}
	return val
}
