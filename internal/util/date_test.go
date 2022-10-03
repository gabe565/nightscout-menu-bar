package util

import (
	"testing"
	"time"
)

func TestMinAgo(t *testing.T) {
	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"simple", args{time.Now().Add(-time.Minute)}, "1m"},
		{"now", args{time.Now()}, "0m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinAgo(tt.args.date); got != tt.want {
				t.Errorf("MinAgo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNextMinChange(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{"now", args{time.Now()}, time.Minute},
		{"1s ago", args{time.Now().Add(-time.Second)}, 59 * time.Second},
		{"15s ago", args{time.Now().Add(-15 * time.Second)}, 45 * time.Second},
		{"30s ago", args{time.Now().Add(-30 * time.Second)}, 30 * time.Second},
		{"45s ago", args{time.Now().Add(-45 * time.Second)}, 15 * time.Second},
		{"59s ago", args{time.Now().Add(-59 * time.Second)}, time.Second},
		{"1m ago", args{time.Now().Add(-time.Minute)}, time.Minute},
		{"4m40s ago", args{time.Now().Add(-4*time.Minute - 40*time.Second)}, 20 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want = tt.want.Round(time.Second)
			got := GetNextMinChange(tt.args.t).Round(time.Second)
			if got != tt.want {
				t.Errorf("GetNextMinChange() = %v, want %v", got, tt.want)
			}
		})
	}
}
