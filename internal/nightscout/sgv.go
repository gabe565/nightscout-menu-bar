package nightscout

type SGV struct {
	ID         string `json:"_id"`
	Device     string `json:"device"`
	Direction  string `json:"direction"`
	Filtered   int    `json:"filtered"`
	Mgdl       int    `json:"mgdl"`
	Mills      Mills  `json:"mills"`
	Noise      int    `json:"noise"`
	Rssi       int    `json:"rssi"`
	Scaled     int    `json:"scaled"`
	Type       string `json:"type"`
	Unfiltered int    `json:"unfiltered"`
}
