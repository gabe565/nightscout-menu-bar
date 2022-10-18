package nightscout

import (
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/spf13/viper"
	"testing"
	"time"
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
		viper.Set(config.UnitsKey, config.UnitsMgdl)
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
		{"95", args{config.UnitsMgdl}, fields{Last: 95}, "95"},
		{"LOW", args{config.UnitsMgdl}, fields{Last: 39}, "LOW"},
		{"HIGH", args{config.UnitsMgdl}, fields{Last: 401}, "HIGH"},
		{"mmol", args{config.UnitsMmol}, fields{Last: 100}, "5.6"},
	}
	for _, tt := range tests {
		switch tt.args.units {
		case config.UnitsMgdl:
			viper.Set(config.UnitsKey, config.UnitsMgdl)
		case config.UnitsMmol:
			viper.Set(config.UnitsKey, config.UnitsMmol)
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
