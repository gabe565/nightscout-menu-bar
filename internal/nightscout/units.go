package nightscout

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	UnitsKey  = "units"
	UnitsMgdl = "mg/dL"
	UnitsMmol = "mmol/L"
)

func init() {
	flag.String("units", UnitsMgdl, "Units to use (one of: mg/dL, mmol/L)")
	if err := viper.BindPFlag(UnitsKey, flag.Lookup(UnitsKey)); err != nil {
		panic(err)
	}
}
