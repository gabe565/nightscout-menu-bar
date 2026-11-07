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

	p, err := s.path()
	if err != nil {
		return err
	}

	b, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			s.creds = nil
			return nil
		}
		return err
	}

	var creds Credentials
	if err := json.Unmarshal(b, &creds); err != nil {
		return err
	}
	s.creds = &creds
	slog.Info("Loaded Nocturne credentials", "file", p, "url", creds.URL)
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
