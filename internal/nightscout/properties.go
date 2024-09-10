package nightscout

import (
	"fyne.io/fyne/v2"
	"github.com/gabe565/nightscout-menu-bar/internal/app/settings"
)

type Properties struct {
	Bgnow     Reading   `json:"bgnow"`
	Buckets   []Reading `json:"buckets"`
	Delta     Delta     `json:"delta"`
	Direction Direction `json:"direction"`
}

func (p Properties) String(prefs fyne.Preferences) string {
	var units settings.Unit
	if err := units.UnmarshalText([]byte(prefs.String(settings.UnitsKey))); err != nil {
		units = settings.UnitMgdl
	}

	result := p.Bgnow.DisplayBg(units) +
		" " + p.Bgnow.Arrow()
	if delta := p.Delta.Display(units); delta != "" {
		result += " " + p.Delta.Display(units)
	}

	if rel := p.Bgnow.Mills.Relative(); rel != "" {
		result += " [" + rel + "]"
	}
	return result
}
