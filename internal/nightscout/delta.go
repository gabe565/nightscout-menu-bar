package nightscout

type Delta struct {
	Absolute     int     `json:"absolute"`
	Display      string  `json:"display"`
	ElapsedMins  float64 `json:"elapsedMins"`
	Interpolated bool    `json:"interpolated"`
	Mean5MinsAgo float64 `json:"mean5MinsAgo"`
	Mgdl         int     `json:"mgdl"`
	Previous     Reading `json:"previous"`
	Scaled       int     `json:"scaled"`
	Times        struct {
		Previous int `json:"previous"`
		Recent   int `json:"recent"`
	} `json:"times"`
}
