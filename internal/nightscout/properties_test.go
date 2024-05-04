package nightscout

import (
	"testing"
	"time"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestProperties_String(t *testing.T) {
	t.Parallel()
	type fields struct {
		Bgnow     Reading
		Buckets   []Reading
		Delta     Delta
		Direction Direction
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
			"mgdl",
			fields{
				Bgnow: Reading{
					Last:  100,
					Mills: Mills{time.Now()},
					Sgvs:  []SGV{{Direction: "Flat"}},
				},
				Delta: Delta{DisplayVal: "+1"},
			},
			args{config.NewDefault()},
			"100 → +1 [0m]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := Properties{
				Bgnow:     tt.fields.Bgnow,
				Buckets:   tt.fields.Buckets,
				Delta:     tt.fields.Delta,
				Direction: tt.fields.Direction,
			}
			assert.Equal(t, tt.want, p.String(tt.args.conf))
		})
	}
}
