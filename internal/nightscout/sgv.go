package nightscout

import "encoding/json"

type SGV struct {
	ID         string      `json:"_id"`
	Device     string      `json:"device"`
	Direction  string      `json:"direction"`
	Filtered   int         `json:"filtered"`
	Mgdl       int         `json:"mgdl"`
	Mills      Mills       `json:"mills"`
	Noise      int         `json:"noise"`
	Rssi       int         `json:"rssi"`
	Scaled     json.Number `json:"scaled"`
	Type       string      `json:"type"`
	Unfiltered int         `json:"unfiltered"`
}
