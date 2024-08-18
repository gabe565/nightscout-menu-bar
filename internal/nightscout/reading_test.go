package nightscout

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReading_Arrow(t *testing.T) {
	t.Parallel()
	defaultArrows := config.New().Arrows

	type fields struct {
		Mean      json.Number
		Last      Mgdl
		Mills     Mills
		Index     json.Number
		FromMills Mills
		ToMills   Mills
		Sgvs      []SGV
	}
	type args struct {
		arrows config.Arrows
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"TripleUp", fields{Sgvs: []SGV{{Direction: "TripleUp"}}}, args{defaultArrows}, "⇈"},
		{"DoubleUp", fields{Sgvs: []SGV{{Direction: "DoubleUp"}}}, args{defaultArrows}, "⇈"},
		{"SingleUp", fields{Sgvs: []SGV{{Direction: "SingleUp"}}}, args{defaultArrows}, "↑"},
		{"FortyFiveUp", fields{Sgvs: []SGV{{Direction: "FortyFiveUp"}}}, args{defaultArrows}, "↗"},
		{"Flat", fields{Sgvs: []SGV{{Direction: "Flat"}}}, args{defaultArrows}, "→"},
		{"FortyFiveDown", fields{Sgvs: []SGV{{Direction: "FortyFiveDown"}}}, args{defaultArrows}, "↘"},
		{"SingleDown", fields{Sgvs: []SGV{{Direction: "SingleDown"}}}, args{defaultArrows}, "↓"},
		{"DoubleDown", fields{Sgvs: []SGV{{Direction: "DoubleDown"}}}, args{defaultArrows}, "⇊"},
		{"TripleDown", fields{Sgvs: []SGV{{Direction: "TripleDown"}}}, args{defaultArrows}, "⇊"},
		{"unknown", fields{}, args{defaultArrows}, "-"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := Reading{
				Mean:      tt.fields.Mean,
				Last:      tt.fields.Last,
				Mills:     tt.fields.Mills,
				Index:     tt.fields.Index,
				FromMills: tt.fields.FromMills,
				ToMills:   tt.fields.ToMills,
				Sgvs:      tt.fields.Sgvs,
			}
			assert.Equal(t, tt.want, r.Arrow(tt.args.arrows))
		})
	}
}

func TestReading_String(t *testing.T) {
	t.Parallel()
	type fields struct {
		Mean      json.Number
		Last      Mgdl
		Mills     Mills
		Index     json.Number
		FromMills Mills
		ToMills   Mills
		Sgvs      []SGV
	}
	type args struct {
		conf *config.Config
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"simple",
			fields{
				Last:  100,
				Mills: Mills{time.Now()},
				Sgvs:  []SGV{{Direction: "Flat"}},
			},
			args{config.New()},
			"100 → [0m]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := Reading{
				Mean:      tt.fields.Mean,
				Last:      tt.fields.Last,
				Mills:     tt.fields.Mills,
				Index:     tt.fields.Index,
				FromMills: tt.fields.FromMills,
				ToMills:   tt.fields.ToMills,
				Sgvs:      tt.fields.Sgvs,
			}
			assert.Equal(t, tt.want, r.String(tt.args.conf))
		})
	}
}

func TestReading_DisplayBg(t *testing.T) {
	t.Parallel()
	type fields struct {
		Mean      json.Number
		Last      Mgdl
		Mills     Mills
		Index     json.Number
		FromMills Mills
		ToMills   Mills
		Sgvs      []SGV
	}
	type args struct {
		units config.Unit
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   string
	}{
		{"95", args{config.UnitMgdl}, fields{Last: 95}, "95"},
		{"LOW", args{config.UnitMgdl}, fields{Last: 39}, "LOW"},
		{"HIGH", args{config.UnitMgdl}, fields{Last: 401}, "HIGH"},
		{"mmol", args{config.UnitMmol}, fields{Last: 100}, "5.6"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &Reading{
				Mean:      tt.fields.Mean,
				Last:      tt.fields.Last,
				Mills:     tt.fields.Mills,
				Index:     tt.fields.Index,
				FromMills: tt.fields.FromMills,
				ToMills:   tt.fields.ToMills,
				Sgvs:      tt.fields.Sgvs,
			}
			assert.Equal(t, tt.want, r.DisplayBg(tt.args.units))
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
	t.Parallel()
	now := time.Now()

	type fields struct {
		Mean      json.Number
		Last      Mgdl
		Mills     Mills
		Index     json.Number
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
		wantErr require.ErrorAssertionFunc
	}{
		{
			"simple",
			fields{},
			args{[]byte(fmt.Sprintf(normalReading, now.UnixMilli(), now.UnixMilli()))},
			require.NoError,
		},
		{
			"low",
			fields{},
			args{[]byte(fmt.Sprintf(lowReading, now.UnixMilli()))},
			require.NoError,
		},
		{
			"error",
			fields{},
			args{[]byte("{")},
			require.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &Reading{
				Mean:      tt.fields.Mean,
				Last:      tt.fields.Last,
				Mills:     tt.fields.Mills,
				Index:     tt.fields.Index,
				FromMills: tt.fields.FromMills,
				ToMills:   tt.fields.ToMills,
				Sgvs:      tt.fields.Sgvs,
			}
			tt.wantErr(t, r.UnmarshalJSON(tt.args.bytes))
		})
	}
}
