package nightscout

import "github.com/spf13/viper"

const (
	UnitsKey  = "units"
	UnitsMgdl = "mg/dL"
	UnitsMmol = "mmol/L"
)

func init() {
	viper.SetDefault(UnitsKey, UnitsMgdl)
}
