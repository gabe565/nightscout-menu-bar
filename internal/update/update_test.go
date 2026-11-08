package update

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecker_Check(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"tag_name":"v1.2.0","html_url":"https://example.com/v1.2.0"}`))
	}))
	t.Cleanup(server.Close)

	tests := []struct {
		version string
		want    *Release
	}{
		{"v1.1.9", &Release{Version: "v1.2.0", URL: "https://example.com/v1.2.0"}},
		{"v1.2.0", nil},
		{"v1.3.0", nil},
	}
	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			t.Parallel()
			c := NewChecker(tt.version)
			c.url = server.URL
			got, err := c.Check(t.Context())
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestChecker_CheckError(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.NotFoundHandler())
	t.Cleanup(server.Close)

	c := NewChecker("v1.0.0")
	c.url = server.URL
	_, err := c.Check(t.Context())
	require.ErrorIs(t, err, ErrUnexpectedStatus)
}

func TestChecker_Enabled(t *testing.T) {
	t.Parallel()
	assert.True(t, NewChecker("v1.0.0").Enabled())
	assert.False(t, NewChecker("beta").Enabled())
}

func TestInstallMethod(t *testing.T) {
	t.Parallel()
	none := func(string) bool { return false }
	cask := func(path string) bool { return path == "/opt/homebrew/Caskroom/nightscout-menu-bar" }

	tests := []struct {
		name   string
		goos   string
		exe    string
		exists func(string) bool
		want   Method
	}{
		{
			"linux brew",
			"linux",
			"/home/linuxbrew/.linuxbrew/Caskroom/nightscout-menu-bar/1.0.0/nightscout-menu-bar",
			none,
			MethodHomebrew,
		},
		{"linux package", "linux", "/usr/bin/nightscout-menu-bar", none, MethodPackage},
		{"linux tarball", "linux", "/home/user/bin/nightscout-menu-bar", none, MethodManual},
		{
			"darwin brew",
			"darwin",
			"/Applications/Nightscout Menu Bar.app/Contents/MacOS/nightscout-menu-bar",
			cask,
			MethodHomebrew,
		},
		{
			"darwin manual",
			"darwin",
			"/Applications/Nightscout Menu Bar.app/Contents/MacOS/nightscout-menu-bar",
			none,
			MethodManual,
		},
		{"windows", "windows", `C:\Program Files\Nightscout Menu Bar.exe`, none, MethodManual},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, installMethod(tt.goos, tt.exe, tt.exists))
		})
	}
}
