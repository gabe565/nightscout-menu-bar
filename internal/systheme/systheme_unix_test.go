//go:build !darwin && !windows

package systheme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMenuLight(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		scheme  uint32
		desktop string
		want    bool
	}{
		{"dark", preferDark, "KDE", false},
		{"light", preferLight, "KDE", true},
		{"stock GNOME light style", noPreference, "GNOME", false},
		{"Ubuntu light style", noPreference, "ubuntu:GNOME", true},
		{"unknown desktop", noPreference, "XFCE", true},
		{"no portal", noPreference, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, menuLight(tt.scheme, tt.desktop))
		})
	}
}
