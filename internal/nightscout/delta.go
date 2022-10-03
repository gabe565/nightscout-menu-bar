package nightscout

type Times struct {
	Previous Mills `json:"previous"`
	Recent   Mills `json:"recent"`
}

type Delta struct {
	Absolute     int     `json:"absolute"`
	Display      string  `json:"display"`
	ElapsedMins  float64 `json:"elapsedMins"`
	Interpolated bool    `json:"interpolated"`
	Mean5MinsAgo float64 `json:"mean5MinsAgo"`
	Mgdl         int     `json:"mgdl"`
	Previous     Reading `json:"previous"`
	Scaled       int     `json:"scaled"`
	Times        Times   `json:"times"`
}
