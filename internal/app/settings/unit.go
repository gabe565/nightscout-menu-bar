package settings

const MmolConversionFactor = 0.0555

//go:generate enumer -type Unit -linecomment -text

type Unit uint8

const (
	UnitMgdl Unit = iota // mg/dL
	UnitMmol             // mmol/L
)
