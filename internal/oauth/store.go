package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/util"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/oauth2"
)

const (
	credentialsFile = "oauth.json" //nolint:gosec
	refreshBefore   = time.Minute
)

var (
	ErrNotConnected = errors.New("not connected to Nocturne")
	//nolint:staticcheck
	ErrExpired = errors.New("Nocturne connection expired or was revoked; reconnect in Preferences")
)

// Credentials are stored separately from the config file since the refresh token
// rotates on every refresh, and writing it to the config would trigger a reload.
type Credentials struct {
	URL          string    `json:"url"`
	ClientID     string    `json:"client_id"`
	TokenURL     string    `json:"token_url"`
	RevokeURL    string    `json:"revoke_url,omitempty"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

func NewStore(conf *config.Config) *Store {
	return &Store{
		config: conf,
		client: &http.Client{
			Transport: util.NewUserAgentTransport("nightscout-menu-bar", conf.Version),
			Timeout:   time.Minute,
		},
	}
}

type Store struct {
	mu     sync.Mutex
	config *config.Config
	client *http.Client
	creds  *Credentials
}

func (s *Store) path() (string, error) {
	if s.config.File != "" {
		return filepath.Join(filepath.Dir(s.config.File), credentialsFile), nil
	}
	dir, err := config.GetDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, credentialsFile), nil
}

func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.loadLocked()
	return err
}

func (s *Store) loadLocked() (bool, error) {
	p, err := s.path()
	if err != nil {
		return false, err
	}

	b, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			changed := s.creds != nil
			s.creds = nil
			return changed, nil
		}
		return false, err
	}

	var creds Credentials
	if err := json.Unmarshal(b, &creds); err != nil {
		return false, err
	}
	if s.creds != nil && *s.creds == creds {
		return false, nil
	}
	s.creds = &creds
	slog.Info("Loaded Nocturne credentials", "file", p, "url", creds.URL)
	return true, nil
}

// Watch reloads credentials changed by another process, like the settings window.
func (s *Store) Watch(ctx context.Context, onChange func()) error {
	p, err := s.path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o777); err != nil {
		return err
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	// The file is replaced atomically, so the directory must be watched
	if err := w.Add(filepath.Dir(p)); err != nil {
		_ = w.Close()
		return err
	}

	go func() {
		defer func() {
			_ = w.Close()
		}()

		var debounce <-chan time.Time
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-w.Events:
				if !ok {
					return
				}
				if filepath.Base(event.Name) == credentialsFile {
					debounce = time.After(100 * time.Millisecond)
				}
			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				slog.Error("Nocturne credentials watcher failed", "error", err)
			case <-debounce:
				s.mu.Lock()
				changed, err := s.loadLocked()
				s.mu.Unlock()
				if err != nil {
					slog.Error("Failed to reload Nocturne credentials", "error", err)
				} else if changed {
					onChange()
				}
			}
		}
	}()
	return nil
}

func (s *Store) saveLocked(creds *Credentials) error {
	p, err := s.path()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(creds, "", "  ") //nolint:gosec // Stored with 0600 permissions
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(p), 0o777); err != nil {
		return err
	}

	// Write atomically so a crash can't lose a rotated refresh token
	tmp, err := os.CreateTemp(filepath.Dir(p), credentialsFile+".*")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(tmp.Name())
	}()
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(b); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmp.Name(), p); err != nil {
		return err
	}

	s.creds = creds
	return nil
}

func (s *Store) clearLocked() error {
	s.creds = nil
	p, err := s.path()
	if err != nil {
		return err
	}
	if err := os.Remove(p); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func normalizeURL(u string) string {
	return strings.TrimRight(u, "/")
}

// Supported reports whether the server supports connected apps.
func (s *Store) Supported(ctx context.Context, serverURL string) bool {
	_, err := Discover(ctx, s.client, serverURL)
	if err != nil {
		slog.Debug("Server does not support connected apps", "url", serverURL, "error", err)
	}
	return err == nil
}

func (s *Store) Connected(serverURL string) bool {
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.creds != nil && serverURL != "" && normalizeURL(s.creds.URL) == normalizeURL(serverURL)
}

func (s *Store) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.creds == nil {
		return "", ErrNotConnected
	}

	if s.creds.AccessToken != "" && time.Until(s.creds.Expiry) > refreshBefore {
		return s.creds.AccessToken, nil
	}

	slog.Debug("Refreshing Nocturne access token")
	conf := &oauth2.Config{
		ClientID: s.creds.ClientID,
		Endpoint: oauth2.Endpoint{
			TokenURL:  s.creds.TokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, s.client)
	tok, err := conf.TokenSource(ctx, &oauth2.Token{
		RefreshToken: s.creds.RefreshToken,
		Expiry:       time.Unix(1, 0),
	}).Token()
	if err != nil {
		if retrieveErr := (*oauth2.RetrieveError)(
			nil,
		); errors.As(err, &retrieveErr) &&
			retrieveErr.ErrorCode == "invalid_grant" {
			slog.Warn("Nocturne refresh token was rejected; clearing credentials", "error", err)
			if err := s.clearLocked(); err != nil {
				slog.Error("Failed to remove Nocturne credentials", "error", err)
			}
			return "", ErrExpired
		}
		return "", err
	}

	creds := *s.creds
	creds.AccessToken = tok.AccessToken
	creds.Expiry = tok.Expiry
	if tok.RefreshToken != "" {
		creds.RefreshToken = tok.RefreshToken
	}
	if err := s.saveLocked(&creds); err != nil {
		// Keep the rotated token in memory even if it couldn't be persisted
		slog.Error("Failed to save Nocturne credentials", "error", err)
		s.creds = &creds
	}
	return creds.AccessToken, nil
}

// Connect blocks until the user approves or denies the device code passed to onCode.
func (s *Store) Connect(ctx context.Context, serverURL string, onCode func(*oauth2.DeviceAuthResponse)) error {
	meta, err := Discover(ctx, s.client, serverURL)
	if err != nil {
		return err
	}

	clientID, err := Register(ctx, s.client, meta.RegistrationEndpoint, s.config.Version)
	if err != nil {
		return err
	}
	slog.Debug("Registered Nocturne client", "client_id", clientID)

	conf := oauthConfig(clientID, meta)
	ctx = context.WithValue(ctx, oauth2.HTTPClient, s.client)
	auth, err := conf.DeviceAuth(ctx)
	if err != nil {
		return err
	}
	onCode(auth)

	// Nocturne answers polls that arrive even slightly early with slow_down, which adds 5s
	// to the interval each time. Poll a second slower than asked to absorb network jitter.
	if auth.Interval == 0 {
		auth.Interval = 5
	}
	auth.Interval++

	tok, err := conf.DeviceAccessToken(ctx, auth)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.saveLocked(&Credentials{
		URL:          normalizeURL(serverURL),
		ClientID:     clientID,
		TokenURL:     meta.TokenEndpoint,
		RevokeURL:    meta.RevocationEndpoint,
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		Expiry:       tok.Expiry,
	}); err != nil {
		return err
	}
	slog.Info("Connected to Nocturne", "url", serverURL)
	return nil
}

func (s *Store) Disconnect(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Another process may have rotated the refresh token since it was loaded
	if _, err := s.loadLocked(); err != nil {
		slog.Warn("Failed to reload Nocturne credentials", "error", err)
	}

	if s.creds == nil {
		return nil
	}

	if s.creds.RevokeURL != "" && s.creds.RefreshToken != "" {
		if err := s.revokeLocked(ctx); err != nil {
			slog.Warn("Failed to revoke Nocturne token", "error", err)
		}
	}

	if err := s.clearLocked(); err != nil {
		return err
	}
	slog.Info("Disconnected from Nocturne")
	return nil
}

func (s *Store) revokeLocked(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	form := url.Values{
		"token":           {s.creds.RefreshToken},
		"token_type_hint": {"refresh_token"},
		"client_id":       {s.creds.ClientID},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.creds.RevokeURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return nil
}
