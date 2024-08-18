// Code generated by "enumer -type Unit -linecomment -text"; DO NOT EDIT.

package config

import (
	"fmt"
	"strings"
)

const _UnitName = "mg/dLmmol/L"

var _UnitIndex = [...]uint8{0, 5, 11}

const _UnitLowerName = "mg/dlmmol/l"

func (i Unit) String() string {
	if i >= Unit(len(_UnitIndex)-1) {
		return fmt.Sprintf("Unit(%d)", i)
	}
	return _UnitName[_UnitIndex[i]:_UnitIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _UnitNoOp() {
	var x [1]struct{}
	_ = x[UnitMgdl-(0)]
	_ = x[UnitMmol-(1)]
}

var _UnitValues = []Unit{UnitMgdl, UnitMmol}

var _UnitNameToValueMap = map[string]Unit{
	_UnitName[0:5]:       UnitMgdl,
	_UnitLowerName[0:5]:  UnitMgdl,
	_UnitName[5:11]:      UnitMmol,
	_UnitLowerName[5:11]: UnitMmol,
}

var _UnitNames = []string{
	_UnitName[0:5],
	_UnitName[5:11],
}

// UnitString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func UnitString(s string) (Unit, error) {
	if val, ok := _UnitNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _UnitNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Unit values", s)
}

// UnitValues returns all values of the enum
func UnitValues() []Unit {
	return _UnitValues
}

// UnitStrings returns a slice of all String values of the enum
func UnitStrings() []string {
	strs := make([]string, len(_UnitNames))
	copy(strs, _UnitNames)
	return strs
}

// IsAUnit returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Unit) IsAUnit() bool {
	for _, v := range _UnitValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalText implements the encoding.TextMarshaler interface for Unit
func (i Unit) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Unit
func (i *Unit) UnmarshalText(text []byte) error {
	var err error
	*i, err = UnitString(string(text))
	return err
}
