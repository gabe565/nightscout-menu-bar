package util

import (
	"errors"
	"testing"
)

func TestSoftError_Error(t1 *testing.T) {
	type fields struct {
		Err error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"simple", fields{errors.New("a")}, "a"},
		{"nil error", fields{}, ""},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := SoftError{
				Err: tt.fields.Err,
			}
			if got := t.Error(); got != tt.want {
				t1.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSoftError_Unwrap(t1 *testing.T) {
	err := errors.New("a")

	type fields struct {
		Err error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"simple", fields{err}, true},
		{"nil error", fields{}, false},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := SoftError{
				Err: tt.fields.Err,
			}
			if err := t.Unwrap(); (err != nil) != tt.wantErr {
				t1.Errorf("Unwrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
