package settings

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/dynamicicon"
	"gabe565.com/nightscout-menu-bar/internal/nightscout/testproperties"
	"gabe565.com/nightscout-menu-bar/internal/oauth"
	"gabe565.com/nightscout-menu-bar/internal/secret"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/font/sfnt"
)

func newTestService(t *testing.T) (*Service, *config.Config) {
	t.Helper()
	data := config.New().Data()
	data.URL = "https://example.com"
	dir := t.TempDir()
	conf := config.New(config.WithData(data))
	conf.File = filepath.Join(dir, "config.toml")
	conf.Secrets = secret.NewWithBackend(dir, &secret.Memory{})
	return NewService(conf, oauth.NewStore(conf)), conf
}

func TestService_Save(t *testing.T) {
	t.Parallel()
	s, conf := newTestService(t)

	// Load keeps StartOnLogin at its current value, so the real login item isn't changed
	settings := s.Load()
	settings.URL = " https://example.com/nightscout "
	settings.Units = config.UnitMmol.String()
	settings.ShowDelta = false
	settings.FetchDelay = "1m"
	settings.FontColor = "#ff0000"
	require.NoError(t, s.Save(settings))

	data := conf.Data()
	assert.Equal(t, "https://example.com/nightscout", data.URL)
	assert.Equal(t, config.UnitMmol, data.Units)
	assert.True(t, data.LastReading.HideDelta)
	assert.False(t, data.LastReading.HideArrow)
	assert.Equal(t, time.Minute, data.Advanced.FetchDelay.Duration)

	contents, err := os.ReadFile(conf.File)
	require.NoError(t, err)
	assert.Contains(t, string(contents), `url = 'https://example.com/nightscout'`)

	assert.Equal(t, "#ff0000", s.Load().FontColor)
}

func TestService_SaveInvalid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		modify func(*Settings)
		want   error
	}{
		{"duration", func(s *Settings) { s.FallbackInterval = "soon" }, ErrDuration},
		{"negative duration", func(s *Settings) { s.StaleThreshold = "-5m" }, ErrDuration},
		{"font size", func(s *Settings) { s.MaxFontSize = 0 }, ErrFontSize},
		{"units", func(s *Settings) { s.Units = "furlongs" }, ErrUnknownUnit},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, conf := newTestService(t)
			settings := s.Load()
			tt.modify(&settings)
			require.ErrorIs(t, s.Save(settings), tt.want)
			assert.NoFileExists(t, conf.File)
		})
	}
}

func TestService_Preview(t *testing.T) {
	t.Parallel()
	s, _ := newTestService(t)
	settings := s.Load()
	settings.DynamicIcon = false

	p := s.Preview(settings)
	assert.Equal(t, "123 → +2 [3m]", p.Text)
	assert.False(t, p.Stale)
	assert.False(t, p.Live)
	assert.Empty(t, p.Icon)

	settings.ShowDelta = false
	settings.Units = config.UnitMmol.String()
	settings.StaleThreshold = "1m"
	p = s.Preview(settings)
	assert.Equal(t, "6.8 → [3m]", p.Text)
	assert.True(t, p.Stale)

	settings.DynamicIcon = true
	p = s.Preview(settings)
	assert.True(t, strings.HasPrefix(p.Icon, "data:image/png;base64,"))
	assert.Empty(t, p.Error)

	settings.FontFile = "does-not-exist.ttf"
	p = s.Preview(settings)
	assert.Contains(t, p.Error, "Dynamic icon")

	settings.FetchDelay = "soon"
	assert.Contains(t, s.Preview(settings).Error, ErrDuration.Error())
}

func TestService_Test(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/v2/properties") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write(testproperties.JSON)
	}))
	t.Cleanup(server.Close)

	s, _ := newTestService(t)
	settings := s.Load()

	settings.URL = server.URL + "/missing"
	_, err := s.Test(t.Context(), settings)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
	assert.False(t, s.Preview(settings).Live)

	settings.URL = server.URL
	reading, err := s.Test(t.Context(), settings)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(reading, "123 "))

	// The fetched reading replaces the sample in previews
	assert.True(t, s.Preview(settings).Live)
}

func TestService_TestNoURL(t *testing.T) {
	t.Parallel()
	s, _ := newTestService(t)
	settings := s.Load()
	settings.URL = "  "
	_, err := s.Test(t.Context(), settings)
	require.Error(t, err)
}

func TestService_CloseGuard(t *testing.T) {
	t.Parallel()
	s, _ := newTestService(t)

	// The window can close normally until the page has loaded
	assert.True(t, s.canClose())

	s.GuardClose()
	assert.False(t, s.canClose())

	// Close always closes, even without a window
	s.Close()
	assert.True(t, s.canClose())
}

// newFakeNocturne serves just enough of Nocturne's OAuth API to approve a device code.
func newFakeNocturne(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	var server *httptest.Server
	writeJSON := func(w http.ResponseWriter, v string) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(v))
	}
	mux.HandleFunc("GET /.well-known/oauth-authorization-server", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, `{"token_endpoint":"/api/oauth/token","device_authorization_endpoint":"/api/oauth/device",`+
			`"registration_endpoint":"/api/oauth/register"}`)
	})
	mux.HandleFunc("POST /api/oauth/register", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
		writeJSON(w, `{"client_id":"client-1"}`)
	})
	mux.HandleFunc("POST /api/oauth/device", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, `{"device_code":"device-1","user_code":"BCTV-3KMN","verification_uri":"`+server.URL+
			`/oauth/device","expires_in":1800,"interval":1}`)
	})
	mux.HandleFunc("POST /api/oauth/token", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, `{"access_token":"access-1","token_type":"Bearer","expires_in":900,"refresh_token":"refresh-1"}`)
	})
	server = httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server
}

func TestService_ConnectNocturne(t *testing.T) {
	t.Parallel()
	server := newFakeNocturne(t)
	s, conf := newTestService(t)

	var opened string
	s.openURL = func(url string) error {
		opened = url
		return nil
	}

	// An unrelated unsaved change shouldn't be saved by connecting
	before := conf.Data()

	require.NoError(t, s.ConnectNocturne(t.Context(), " "+server.URL+" "))
	assert.Equal(t, server.URL+"/oauth/device", opened)
	assert.True(t, s.NocturneConnected(server.URL))

	// The URL is saved so the tray uses the new link
	data := conf.Data()
	assert.Equal(t, server.URL, data.URL)
	data.URL = before.URL
	assert.Equal(t, before, data)

	contents, err := os.ReadFile(conf.File)
	require.NoError(t, err)
	assert.Contains(t, string(contents), "url = '"+server.URL+"'")
}

func TestService_NocturneSupported(t *testing.T) {
	t.Parallel()
	s, _ := newTestService(t)
	assert.True(t, s.NocturneSupported(t.Context(), newFakeNocturne(t).URL))
	assert.False(t, s.NocturneSupported(t.Context(), " "))

	nightscout := httptest.NewServer(http.NotFoundHandler())
	t.Cleanup(nightscout.Close)
	assert.False(t, s.NocturneSupported(t.Context(), nightscout.URL))
}

func TestFontName(t *testing.T) {
	t.Parallel()
	b, err := dynamicicon.LoadFont("")
	require.NoError(t, err)
	path := filepath.Join(t.TempDir(), "embedded.ttf")
	require.NoError(t, os.WriteFile(path, b, 0o600))

	var buf sfnt.Buffer
	assert.Equal(t, "Roboto Condensed SemiBold", fontName(&buf, path))

	// Falls back to the file name
	notFont := filepath.Join(t.TempDir(), "Not A Font.ttf")
	require.NoError(t, os.WriteFile(notFont, []byte("nope"), 0o600))
	assert.Equal(t, "Not A Font", fontName(&buf, notFont))
}

func TestService_ServeFont(t *testing.T) {
	t.Parallel()
	s, _ := newTestService(t)
	serve := func(file string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		s.ServeFont(w, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/font?file="+file, nil))
		return w
	}

	w := serve("")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "font/ttf", w.Header().Get("Content-Type"))

	assert.Equal(t, http.StatusBadRequest, serve("/etc/passwd.txt").Code)
	assert.Equal(t, http.StatusNotFound, serve("does-not-exist.ttf").Code)
}
