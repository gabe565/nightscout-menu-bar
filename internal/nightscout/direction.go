package nightscout

type Direction struct {
	Display any    `json:"display"`
	Entity  string `json:"entity"`
	Label   string `json:"label"`
	Value   string `json:"value"`
}
