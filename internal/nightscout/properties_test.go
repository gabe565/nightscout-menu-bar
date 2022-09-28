package nightscout

import (
	"testing"
	"time"
)

func TestProperties_String(t *testing.T) {
	type fields struct {
		Bgnow     Reading
		Buckets   []Reading
		Delta     Delta
		Direction Direction
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"simple",
			fields{
				Bgnow: Reading{
					Last:  100,
					Mills: Mills{time.Now()},
					Sgvs:  []SGV{{Direction: "Flat"}},
				},
				Delta: Delta{Display: "+1"},
			},
			"100 → +1 [0m]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Properties{
				Bgnow:     tt.fields.Bgnow,
				Buckets:   tt.fields.Buckets,
				Delta:     tt.fields.Delta,
				Direction: tt.fields.Direction,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
