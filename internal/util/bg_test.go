package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToMmol(t *testing.T) {
	type args struct {
		mgdl int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"100", args{100}, 5.55},
		{"50", args{50}, 2.775},
		{"300", args{300}, 16.65},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ToMmol(tt.args.mgdl))
		})
	}
}
