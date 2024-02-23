package nightscout

import (
	"fmt"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestReading_Arrow(t *testing.T) {
	type fields struct {
		Mean      int
		Last      int
		Mills     Mills
		Index     int
		FromMills Mills
		ToMills   Mills
		Sgvs      []SGV
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"TripleUp", fields{Sgvs: []SGV{{Direction: "TripleUp"}}}, "⇈"},
		{"DoubleUp", fields{Sgvs: []SGV{{Direction: "DoubleUp"}}}, "⇈"},
		{"SingleUp", fields{Sgvs: []SGV{{Direction: "SingleUp"}}}, "↑"},
		{"FortyFiveUp", fields{Sgvs: []SGV{{Direction: "FortyFiveUp"}}}, "↗"},
		{"Flat", fields{Sgvs: []SGV{{Direction: "Flat"}}}, "→"},
		{"FortyFiveDown", fields{Sgvs: []SGV{{Direction: "FortyFiveDown"}}}, "↘"},
		{"SingleDown", fields{Sgvs: []SGV{{Direction: "SingleDown"}}}, "↓"},
		{"DoubleDown", fields{Sgvs: []SGV{{Direction: "DoubleDown"}}}, "⇊"},
		{"TripleDown", fields{Sgvs: []SGV{{Direction: "TripleDown"}}}, "⇊"},
		{"unknown", fields{}, "-"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Reading{
				Mean:      tt.fields.Mean,
				Last:      tt.fields.Last,
				Mills:     tt.fields.Mills,
				Index:     tt.fields.Index,
				FromMills: tt.fields.FromMills,
				ToMills:   tt.fields.ToMills,
				Sgvs:      tt.fields.Sgvs,
			}
			if got := r.Arrow(); got != tt.want {
				t.Errorf("Arrow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReading_String(t *testing.T) {
	type fields struct {
		Mean      int
		Last      int
		Mills     Mills
		Index     int
		FromMills Mills
		ToMills   Mills
		Sgvs      []SGV
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"simple",
			fields{
				Last:  100,
				Mills: Mills{time.Now()},
				Sgvs:  []SGV{{Direction: "Flat"}},
			},
			"100 → [0m]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Reading{
				Mean:      tt.fields.Mean,
				Last:      tt.fields.Last,
				Mills:     tt.fields.Mills,
				Index:     tt.fields.Index,
				FromMills: tt.fields.FromMills,
				ToMills:   tt.fields.ToMills,
				Sgvs:      tt.fields.Sgvs,
			}
			if got := r.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReading_DisplayBg(t *testing.T) {
	defer func() {
		viper.Set(UnitsKey, UnitsMgdl)
	}()

	type fields struct {
		Mean      int
		Last      int
		Mills     Mills
		Index     int
		FromMills Mills
		ToMills   Mills
		Sgvs      []SGV
	}
	type args struct {
		units string
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   string
	}{
		{"95", args{UnitsMgdl}, fields{Last: 95}, "95"},
		{"LOW", args{UnitsMgdl}, fields{Last: 39}, "LOW"},
		{"HIGH", args{UnitsMgdl}, fields{Last: 401}, "HIGH"},
		{"mmol", args{UnitsMmol}, fields{Last: 100}, "5.6"},
	}
	for _, tt := range tests {
		switch tt.args.units {
		case UnitsMgdl:
			viper.Set(UnitsKey, UnitsMgdl)
		case UnitsMmol:
			viper.Set(UnitsKey, UnitsMmol)
		}
		t.Run(tt.name, func(t *testing.T) {
			r := &Reading{
				Mean:      tt.fields.Mean,
				Last:      tt.fields.Last,
				Mills:     tt.fields.Mills,
				Index:     tt.fields.Index,
				FromMills: tt.fields.FromMills,
				ToMills:   tt.fields.ToMills,
				Sgvs:      tt.fields.Sgvs,
			}
			if got := r.DisplayBg(); got != tt.want {
				t.Errorf("DisplayBg() = %v, want %v", got, tt.want)
			}
		})
	}
}

var normalReading = `{
  "mean": 100,
  "last": 100,
  "mills": %d,
  "sgvs": [
    {
      "_id": "a",
      "mgdl": 100,
      "mills": %d,
      "device": "xDrip-DexcomG5",
      "direction": "Flat",
      "filtered": 0,
      "unfiltered": 0,
      "noise": 1,
      "rssi": 100,
      "type": "sgv",
      "scaled": 100
    }
  ]
}`

var lowReading = `{
  "sgvs": [
    {
      "_id": "a",
      "mgdl": 39,
      "mills": %d,
      "device": "xDrip-DexcomG5",
      "direction": "Flat",
      "filtered": 0,
      "unfiltered": 0,
      "noise": 1,
      "rssi": 100,
      "type": "sgv",
      "scaled": 39
    }
  ]
}`

func TestReading_UnmarshalJSON(t *testing.T) {
	now := time.Now()

	type fields struct {
		Mean      int
		Last      int
		Mills     Mills
		Index     int
		FromMills Mills
		ToMills   Mills
		Sgvs      []SGV
	}
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"simple",
			fields{},
			args{[]byte(fmt.Sprintf(normalReading, now.UnixMilli(), now.UnixMilli()))},
			false,
		},
		{
			"low",
			fields{},
			args{[]byte(fmt.Sprintf(lowReading, now.UnixMilli()))},
			false,
		},
		{
			"error",
			fields{},
			args{[]byte("{")},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Reading{
				Mean:      tt.fields.Mean,
				Last:      tt.fields.Last,
				Mills:     tt.fields.Mills,
				Index:     tt.fields.Index,
				FromMills: tt.fields.FromMills,
				ToMills:   tt.fields.ToMills,
				Sgvs:      tt.fields.Sgvs,
			}
			if err := r.UnmarshalJSON(tt.args.bytes); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
