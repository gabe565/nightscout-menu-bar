package nightscout

import (
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
