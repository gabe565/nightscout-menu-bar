package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMinAgo(t *testing.T) {
	t.Parallel()
	type args struct {
		date  time.Time
		round bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"0m rounded", args{time.Now(), true}, "0m"},
		{"59s rounded", args{time.Now().Add(-59 * time.Second), true}, "1m"},
		{"1m rounded", args{time.Now().Add(-time.Minute), true}, "1m"},
		{"1m30s rounded", args{time.Now().Add(-time.Minute - 30*time.Second), true}, "2m"},
		{"2m35s rounded", args{time.Now().Add(-2*time.Minute - 35*time.Second), true}, "3m"},
		{"4m15s rounded", args{time.Now().Add(-4*time.Minute - 15*time.Second), true}, "4m"},
		{"5m1s rounded", args{time.Now().Add(-5*time.Minute - time.Second), true}, "5m"},
		{"now rounded", args{time.Now(), true}, "0m"},

		{"0m", args{time.Now(), false}, "0m"},
		{"59s", args{time.Now().Add(-59 * time.Second), false}, "0m"},
		{"1m", args{time.Now().Add(-time.Minute), false}, "1m"},
		{"1m30s", args{time.Now().Add(-time.Minute - 30*time.Second), false}, "1m"},
		{"2m35s", args{time.Now().Add(-2*time.Minute - 35*time.Second), false}, "2m"},
		{"4m15s", args{time.Now().Add(-4*time.Minute - 15*time.Second), false}, "4m"},
		{"5m1s", args{time.Now().Add(-5*time.Minute - time.Second), false}, "5m"},
		{"now", args{time.Now(), false}, "0m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, MinAgo(tt.args.date, tt.args.round))
		})
	}
}

func TestGetNextMinChange(t *testing.T) {
	t.Parallel()
	type args struct {
		t     time.Time
		round bool
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{"now rounded", args{time.Now(), true}, 30 * time.Second},
		{"1s ago rounded", args{time.Now().Add(-time.Second), true}, 29 * time.Second},
		{"15s ago rounded", args{time.Now().Add(-15 * time.Second), true}, 15 * time.Second},
		{"29s ago rounded", args{time.Now().Add(-29 * time.Second), true}, time.Second},
		{"30s ago rounded", args{time.Now().Add(-30 * time.Second), true}, time.Minute},
		{"31s ago rounded", args{time.Now().Add(-31 * time.Second), true}, 59 * time.Second},
		{"45s ago rounded", args{time.Now().Add(-45 * time.Second), true}, 45 * time.Second},
		{"59s ago rounded", args{time.Now().Add(-59 * time.Second), true}, 31 * time.Second},
		{"1m ago rounded", args{time.Now().Add(-time.Minute), true}, 30 * time.Second},
		{"4m40s ago rounded", args{time.Now().Add(-4*time.Minute - 40*time.Second), true}, 50 * time.Second},

		{"now", args{time.Now(), false}, time.Minute},
		{"1s ago", args{time.Now().Add(-time.Second), false}, 59 * time.Second},
		{"15s ago", args{time.Now().Add(-15 * time.Second), false}, 45 * time.Second},
		{"29s ago", args{time.Now().Add(-29 * time.Second), false}, 31 * time.Second},
		{"30s ago", args{time.Now().Add(-30 * time.Second), false}, 30 * time.Second},
		{"31s ago", args{time.Now().Add(-31 * time.Second), false}, 29 * time.Second},
		{"45s ago", args{time.Now().Add(-45 * time.Second), false}, 15 * time.Second},
		{"59s ago", args{time.Now().Add(-59 * time.Second), false}, time.Second},
		{"1m ago", args{time.Now().Add(-time.Minute), false}, time.Minute},
		{"4m40s ago", args{time.Now().Add(-4*time.Minute - 40*time.Second), false}, 20 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.want = tt.want.Round(time.Second)
			got := GetNextMinChange(tt.args.t, tt.args.round).Round(time.Second)
			assert.Equal(t, tt.want, got)
		})
	}
}
