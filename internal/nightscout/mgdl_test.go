package nightscout

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMgdl_ToMmol(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		m    Mgdl
		want float64
	}{
		{"100", Mgdl(100), 5.55},
		{"50", Mgdl(50), 2.775},
		{"300", Mgdl(300), 16.65},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.InDelta(t, tt.want, tt.m.Mmol(), 0.001)
		})
	}
}
