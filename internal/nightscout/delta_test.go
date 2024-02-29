package nightscout

import (
	"testing"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
)

func TestDelta_Display(t *testing.T) {
	defer func() {
		config.Default.Units = config.UnitsMgdl
	}()

	type fields struct {
		Absolute     int
		DisplayVal   string
		ElapsedMins  float64
		Interpolated bool
		Mean5MinsAgo float64
		Mgdl         int
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.args.units {
			case config.UnitsMgdl:
				config.Default.Units = config.UnitsMgdl
			case config.UnitsMmol:
				config.Default.Units = config.UnitsMmol
			}
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
			if got := d.Display(); got != tt.want {
				t.Errorf("Display() = %v, want %v", got, tt.want)
			}
		})
	}
}
