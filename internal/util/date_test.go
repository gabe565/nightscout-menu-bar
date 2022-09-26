package util

import (
	"testing"
	"time"
)

func TestMinAgo(t *testing.T) {
	type args struct {
		date any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"simple", args{time.Now().Add(-time.Minute)}, "1m"},
		{"now", args{time.Now()}, "0m"},
		{"nil", args{nil}, "0m"},
		{"unix", args{time.Now().Add(-time.Minute).Unix()}, "1m"},
		{"unix64", args{int(time.Now().Add(-time.Minute).Unix())}, "1m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinAgo(tt.args.date); got != tt.want {
				t.Errorf("MinAgo() = %v, want %v", got, tt.want)
			}
		})
	}
}
