package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMinAgo(t *testing.T) {
	t.Parallel()
	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"0m", args{time.Now()}, "0m"},
		{"59s", args{time.Now().Add(-59 * time.Second)}, "1m"},
		{"1m", args{time.Now().Add(-time.Minute)}, "1m"},
		{"1m30s", args{time.Now().Add(-time.Minute - 30*time.Second)}, "2m"},
		{"2m35s", args{time.Now().Add(-2*time.Minute - 35*time.Second)}, "3m"},
		{"4m15s", args{time.Now().Add(-4*time.Minute - 15*time.Second)}, "4m"},
		{"5m1s", args{time.Now().Add(-5*time.Minute - time.Second)}, "5m"},
		{"now", args{time.Now()}, "0m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, MinAgo(tt.args.date))
		})
	}
}

func TestGetNextMinChange(t *testing.T) {
	t.Parallel()
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{"now", args{time.Now()}, 30 * time.Second},
		{"1s ago", args{time.Now().Add(-time.Second)}, 29 * time.Second},
		{"15s ago", args{time.Now().Add(-15 * time.Second)}, 15 * time.Second},
		{"29s ago", args{time.Now().Add(-29 * time.Second)}, time.Second},
		{"30s ago", args{time.Now().Add(-30 * time.Second)}, time.Minute},
		{"31s ago", args{time.Now().Add(-31 * time.Second)}, 59 * time.Second},
		{"45s ago", args{time.Now().Add(-45 * time.Second)}, 45 * time.Second},
		{"59s ago", args{time.Now().Add(-59 * time.Second)}, 31 * time.Second},
		{"1m ago", args{time.Now().Add(-time.Minute)}, 30 * time.Second},
		{"4m40s ago", args{time.Now().Add(-4*time.Minute - 40*time.Second)}, 50 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.want = tt.want.Round(time.Second)
			got := GetNextMinChange(tt.args.t).Round(time.Second)
			assert.Equal(t, tt.want, got)
		})
	}
}
