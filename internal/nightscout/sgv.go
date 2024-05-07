package nightscout

import "encoding/json"

type SGV struct {
	ID         string      `json:"_id"`
	Device     string      `json:"device"`
	Direction  string      `json:"direction"`
	Filtered   json.Number `json:"filtered"`
	Mgdl       int         `json:"mgdl"`
	Mills      Mills       `json:"mills"`
	Noise      json.Number `json:"noise"`
	Rssi       json.Number `json:"rssi"`
	Scaled     json.Number `json:"scaled"`
	Type       string      `json:"type"`
	Unfiltered json.Number `json:"unfiltered"`
}
