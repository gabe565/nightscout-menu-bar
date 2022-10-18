package util

const ConversionFactor = 0.0555

func ToMmol(mgdl int) float64 {
	return float64(mgdl) * ConversionFactor
}
