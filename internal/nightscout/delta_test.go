package nightscout

import (
	"encoding/json"
	"testing"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestDelta_Display(t *testing.T) {
	t.Parallel()
	type fields struct {
		Absolute     json.Number
		DisplayVal   string
		ElapsedMins  json.Number
		Interpolated bool
		Mean5MinsAgo json.Number
		Mgdl         json.Number
		Previous     Reading
		Scaled       int
		Times        Times
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
		{
			"mgdl",
			args{config.UnitsMgdl},
			fields{DisplayVal: "+1"},
			"+1",
		},
		{
			"mmol",
			args{config.UnitsMmol},
			fields{Scaled: 9},
			"+0.5",
		},
		{
			"mmol no decimal",
			args{config.UnitsMmol},
			fields{Scaled: 0},
			"+0",
		},
		{
			"mmol negative",
			args{config.UnitsMmol},
			fields{Scaled: -9},
			"-0.5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			d := Delta{
				Absolute:     tt.fields.Absolute,
				DisplayVal:   tt.fields.DisplayVal,
				ElapsedMins:  tt.fields.ElapsedMins,
				Interpolated: tt.fields.Interpolated,
				Mean5MinsAgo: tt.fields.Mean5MinsAgo,
				Mgdl:         tt.fields.Mgdl,
				Previous:     tt.fields.Previous,
				Scaled:       tt.fields.Scaled,
				Times:        tt.fields.Times,
			}
			assert.Equal(t, tt.want, d.Display(tt.args.units))
		})
	}
}
