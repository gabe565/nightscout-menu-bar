package oauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

type fakeServer struct {
	*httptest.Server
	polls        atomic.Int32
	refreshToken atomic.Value
	revoked      atomic.Bool
}

func newFakeServer(t *testing.T) *fakeServer {
	t.Helper()
	f := &fakeServer{}
	f.refreshToken.Store("refresh-1")

	writeJSON := func(w http.ResponseWriter, status int, v any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(v)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /.well-known/oauth-authorization-server", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"token_endpoint":                f.URL + "/api/oauth/token",
			"device_authorization_endpoint": "/api/oauth/device",
			"registration_endpoint":         "/api/oauth/register",
			"revocation_endpoint":           "/api/oauth/revoke",
		})
	})
	mux.HandleFunc("POST /api/oauth/register", func(w http.ResponseWriter, r *http.Request) {
		var req registrationRequest
		if !assert.NoError(t, json.NewDecoder(r.Body).Decode(&req)) {
			return
		}
		assert.Equal(t, SoftwareID, req.SoftwareID)
		assert.Equal(t, Scope, req.Scope)
		assert.NotNil(t, req.RedirectURIs)
		writeJSON(w, http.StatusCreated, map[string]string{"client_id": "client-1"})
	})
	mux.HandleFunc("POST /api/oauth/device", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "client-1", r.FormValue("client_id"))
		assert.Equal(t, Scope, r.FormValue("scope"))
		writeJSON(w, http.StatusOK, map[string]any{
			"device_code":               "device-1",
			"user_code":                 "BCTV-3KMN",
			"verification_uri":          f.URL + "/oauth/device",
			"verification_uri_complete": f.URL + "/oauth/device?user_code=BCTV-3KMN",
			"expires_in":                1800,
			"interval":                  1,
		})
	})
	mux.HandleFunc("POST /api/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		switch r.FormValue("grant_type") {
		case "urn:ietf:params:oauth:grant-type:device_code":
			assert.Equal(t, "device-1", r.FormValue("device_code"))
			if f.polls.Add(1) == 1 {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "authorization_pending"})
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{
				"access_token":  "access-1",
				"token_type":    "Bearer",
				"expires_in":    900,
				"refresh_token": "refresh-1",
				"scope":         Scope,
			})
		case "refresh_token":
			if r.FormValue("refresh_token") != f.refreshToken.Load() {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_grant"})
				return
			}
			f.refreshToken.Store("refresh-2")
			writeJSON(w, http.StatusOK, map[string]any{
				"access_token":  "access-2",
				"token_type":    "Bearer",
				"expires_in":    900,
				"refresh_token": "refresh-2",
			})
		default:
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unsupported_grant_type"})
		}
	})
	mux.HandleFunc("POST /api/oauth/revoke", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, f.refreshToken.Load(), r.FormValue("token"))
		f.revoked.Store(true)
		w.WriteHeader(http.StatusOK)
	})

	f.Server = httptest.NewServer(mux)
	t.Cleanup(f.Close)
	return f
}

func newTestStore(t *testing.T) *Store {
	t.Helper()
	conf := config.New()
	conf.File = filepath.Join(t.TempDir(), "config.toml")
	return NewStore(conf)
}

func TestStore_Connect(t *testing.T) {
	t.Parallel()
	server := newFakeServer(t)
	store := newTestStore(t)

	var code *oauth2.DeviceAuthResponse
	require.NoError(t, store.Connect(t.Context(), server.URL+"/", func(auth *oauth2.DeviceAuthResponse) {
		code = auth
	}))
	require.NotNil(t, code)
	assert.Equal(t, "BCTV-3KMN", code.UserCode)
	assert.EqualValues(t, 2, server.polls.Load())

	assert.True(t, store.Connected(server.URL))
	assert.True(t, store.Connected(server.URL+"/"))
	assert.False(t, store.Connected("https://example.com"))

	token, err := store.Token(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "access-1", token)

	// Credentials persist with restricted permissions
	p, err := store.path()
	require.NoError(t, err)
	stat, err := os.Stat(p)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), stat.Mode().Perm())

	reloaded := NewStore(store.config)
	require.NoError(t, reloaded.Load())
	assert.True(t, reloaded.Connected(server.URL))
	assert.Equal(t, server.URL+"/api/oauth/revoke", reloaded.creds.RevokeURL)
}

func TestStore_Token(t *testing.T) {
	t.Parallel()

	t.Run("not connected", func(t *testing.T) {
		t.Parallel()
		_, err := newTestStore(t).Token(t.Context())
		require.ErrorIs(t, err, ErrNotConnected)
	})

	t.Run("refreshes and rotates", func(t *testing.T) {
		t.Parallel()
		server := newFakeServer(t)
		store := newTestStore(t)
		require.NoError(t, store.saveLocked(&Credentials{
			URL:          server.URL,
			ClientID:     "client-1",
			TokenURL:     server.URL + "/api/oauth/token",
			AccessToken:  "access-1",
			RefreshToken: "refresh-1",
			Expiry:       time.Now().Add(30 * time.Second),
		}))

		token, err := store.Token(t.Context())
		require.NoError(t, err)
		assert.Equal(t, "access-2", token)

		reloaded := NewStore(store.config)
		require.NoError(t, reloaded.Load())
		assert.Equal(t, "refresh-2", reloaded.creds.RefreshToken)
		assert.Equal(t, "access-2", reloaded.creds.AccessToken)
	})

	t.Run("invalid grant clears credentials", func(t *testing.T) {
		t.Parallel()
		server := newFakeServer(t)
		store := newTestStore(t)
		require.NoError(t, store.saveLocked(&Credentials{
			URL:          server.URL,
			ClientID:     "client-1",
			TokenURL:     server.URL + "/api/oauth/token",
			RefreshToken: "revoked",
		}))

		_, err := store.Token(t.Context())
		require.ErrorIs(t, err, ErrExpired)
		assert.False(t, store.Connected(server.URL))

		p, err := store.path()
		require.NoError(t, err)
		assert.NoFileExists(t, p)
	})
}

func TestStore_Disconnect(t *testing.T) {
	t.Parallel()
	server := newFakeServer(t)
	store := newTestStore(t)
	require.NoError(t, store.saveLocked(&Credentials{
		URL:          server.URL,
		ClientID:     "client-1",
		RevokeURL:    server.URL + "/api/oauth/revoke",
		RefreshToken: "refresh-1",
	}))

	require.NoError(t, store.Disconnect(t.Context()))
	assert.True(t, server.revoked.Load())
	assert.False(t, store.Connected(server.URL))
}

func TestDiscover_Unsupported(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.NotFoundHandler())
	t.Cleanup(server.Close)

	_, err := Discover(t.Context(), server.Client(), server.URL)
	require.ErrorIs(t, err, ErrUnsupported)
}

// Nocturne answers polls that arrive before the interval has elapsed with slow_down,
// with no tolerance for network jitter.
func TestStore_ConnectPollsWithMargin(t *testing.T) {
	t.Parallel()
	var (
		mu       sync.Mutex
		lastPoll time.Time
		polls    int
		slowDown int
	)
	mux := http.NewServeMux()
	var server *httptest.Server
	writeJSON := func(w http.ResponseWriter, status int, v any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(v)
	}
	mux.HandleFunc("GET /.well-known/oauth-authorization-server", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"token_endpoint":                "/token",
			"device_authorization_endpoint": "/device",
			"registration_endpoint":         "/register",
		})
	})
	mux.HandleFunc("POST /register", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusCreated, map[string]string{"client_id": "client-1"})
	})
	mux.HandleFunc("POST /device", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"device_code": "device-1", "user_code": "CODE", "verification_uri": server.URL, "interval": 1,
		})
	})
	mux.HandleFunc("POST /token", func(w http.ResponseWriter, _ *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		now := time.Now()
		early := !lastPoll.IsZero() && now.Sub(lastPoll) < time.Second
		lastPoll = now
		if early {
			slowDown++
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "slow_down"})
			return
		}
		if polls++; polls < 3 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "authorization_pending"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"access_token": "access-1", "token_type": "Bearer", "expires_in": 900, "refresh_token": "refresh-1",
		})
	})
	server = httptest.NewServer(mux)
	t.Cleanup(server.Close)

	store := newTestStore(t)
	require.NoError(t, store.Connect(t.Context(), server.URL, func(*oauth2.DeviceAuthResponse) {}))
	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, 3, polls)
	assert.Zero(t, slowDown)
}
