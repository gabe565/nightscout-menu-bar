package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHexColor_MarshalText(t *testing.T) {
	t.Parallel()
	type fields struct {
		HexColor HexColor
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr require.ErrorAssertionFunc
	}{
		{"white", fields{HexColor{R: 0xFF, G: 0xFF, B: 0xFF}}, []byte("#fff"), require.NoError},
		{"black", fields{HexColor{}}, []byte("#000"), require.NoError},
		{"red", fields{HexColor{R: 0xFF}}, []byte("#f00"), require.NoError},
		{"green", fields{HexColor{G: 0xFF}}, []byte("#0f0"), require.NoError},
		{"blue", fields{HexColor{B: 0xFF}}, []byte("#00f"), require.NoError},
		{"blue-gray", fields{HexColor{R: 0x60, G: 0x7D, B: 0x8B}}, []byte("#607d8b"), require.NoError},
		{"increment", fields{HexColor{R: 1, G: 2, B: 3}}, []byte("#010203"), require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := tt.fields.HexColor
			got, err := h.MarshalText()
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHexColor_UnmarshalText(t *testing.T) {
	t.Parallel()
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		args    args
		want    HexColor
		wantErr require.ErrorAssertionFunc
	}{
		{"white", args{[]byte("#fff")}, HexColor{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}, require.NoError},
		{"black", args{[]byte("#000")}, HexColor{A: 0xFF}, require.NoError},
		{"red", args{[]byte("#f00")}, HexColor{R: 0xFF, A: 0xFF}, require.NoError},
		{"green", args{[]byte("#0f0")}, HexColor{G: 0xFF, A: 0xFF}, require.NoError},
		{"blue", args{[]byte("#00f")}, HexColor{B: 0xFF, A: 0xFF}, require.NoError},
		{"blue-gray", args{[]byte("#607d8b")}, HexColor{R: 0x60, G: 0x7D, B: 0x8B, A: 0xFF}, require.NoError},
		{"missing-prefix", args{[]byte("fff")}, HexColor{}, require.Error},
		{"too-long", args{[]byte("#fffffff")}, HexColor{}, require.Error},
		{"too-short", args{[]byte("#fffff")}, HexColor{}, require.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := HexColor{}
			tt.wantErr(t, h.UnmarshalText(tt.args.text))
			assert.Equal(t, tt.want, h)
		})
	}
}
