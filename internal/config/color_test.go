package config

import (
	"reflect"
	"testing"
)

func TestHexColor_MarshalText(t *testing.T) {
	type fields struct {
		HexColor HexColor
	}
	tests := []struct {
		name     string
		fields   fields
		wantText []byte
		wantErr  bool
	}{
		{"white", fields{HexColor{R: 0xFF, G: 0xFF, B: 0xFF}}, []byte("#fff"), false},
		{"black", fields{HexColor{}}, []byte("#000"), false},
		{"red", fields{HexColor{R: 0xFF}}, []byte("#f00"), false},
		{"green", fields{HexColor{G: 0xFF}}, []byte("#0f0"), false},
		{"blue", fields{HexColor{B: 0xFF}}, []byte("#00f"), false},
		{"blue-gray", fields{HexColor{R: 0x60, G: 0x7D, B: 0x8B}}, []byte("#607d8b"), false},
		{"increment", fields{HexColor{R: 1, G: 2, B: 3}}, []byte("#010203"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.fields.HexColor
			gotText, err := h.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotText, tt.wantText) {
				t.Errorf("MarshalText() gotText = %v, want %v", string(gotText), string(tt.wantText))
			}
		})
	}
}

func TestHexColor_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		args    args
		want    HexColor
		wantErr bool
	}{
		{"white", args{[]byte("#fff")}, HexColor{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}, false},
		{"black", args{[]byte("#000")}, HexColor{A: 0xFF}, false},
		{"red", args{[]byte("#f00")}, HexColor{R: 0xFF, A: 0xFF}, false},
		{"green", args{[]byte("#0f0")}, HexColor{G: 0xFF, A: 0xFF}, false},
		{"blue", args{[]byte("#00f")}, HexColor{B: 0xFF, A: 0xFF}, false},
		{"blue-gray", args{[]byte("#607d8b")}, HexColor{R: 0x60, G: 0x7D, B: 0x8B, A: 0xFF}, false},
		{"missing-prefix", args{[]byte("fff")}, HexColor{}, true},
		{"too-long", args{[]byte("#fffffff")}, HexColor{}, true},
		{"too-short", args{[]byte("#fffff")}, HexColor{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HexColor{}
			if err := h.UnmarshalText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("MarshalText() gotText = %v, want %v", h, tt.want)
			}
		})
	}
}
