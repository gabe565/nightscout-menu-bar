package nightscout

import "github.com/gabe565/nightscout-menu-bar/internal/app/settings"

type Mgdl int

func (m Mgdl) Mgdl() int { return int(m) }

func (m Mgdl) Mmol() float64 { return float64(m) * settings.MmolConversionFactor }
