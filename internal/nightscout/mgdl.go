package nightscout

import "gabe565.com/nightscout-menu-bar/internal/config"

type Mgdl int

func (m Mgdl) Mgdl() int { return int(m) }

func (m Mgdl) Mmol() float64 { return float64(m) * config.MmolConversionFactor }
