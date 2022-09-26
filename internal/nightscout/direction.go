package nightscout

type Direction struct {
	Display interface{} `json:"display"`
	Entity  string      `json:"entity"`
	Label   string      `json:"label"`
	Value   string      `json:"value"`
}
