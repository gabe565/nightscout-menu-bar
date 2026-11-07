package oauth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"golang.org/x/oauth2"
)

const (
	ClientName = "Nightscout Menu Bar"
	SoftwareID = "com.gabe565.nightscout-menu-bar"
	ClientURI  = "https://github.com/gabe565/nightscout-menu-bar"
	Scope      = "glucose.read"
)

var (
	ErrUnsupported  = errors.New("server does not support connected apps; is this a Nocturne server?")
	ErrRegistration = errors.New("client registration failed")
)

// Metadata is the subset of RFC 8414 authorization server metadata used by the device flow.
type Metadata struct {
	TokenEndpoint               string `json:"token_endpoint"`
	DeviceAuthorizationEndpoint string `json:"device_authorization_endpoint"`
	RegistrationEndpoint        string `json:"registration_endpoint"`
	RevocationEndpoint          string `json:"revocation_endpoint"`
}

func Discover(ctx context.Context, client *http.Client, baseURL string) (*Metadata, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	u := *base
	u.Path = path.Join(u.Path, ".well-known", "oauth-authorization-server")
	u.RawQuery = ""

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w (discovery returned %d)", ErrUnsupported, resp.StatusCode)
	}

	var meta Metadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnsupported, err)
	}
	if meta.TokenEndpoint == "" || meta.DeviceAuthorizationEndpoint == "" || meta.RegistrationEndpoint == "" {
		return nil, fmt.Errorf("%w (device flow not advertised)", ErrUnsupported)
	}

	// Endpoints may be relative to the server
	for _, endpoint := range []*string{
		&meta.TokenEndpoint,
		&meta.DeviceAuthorizationEndpoint,
		&meta.RegistrationEndpoint,
		&meta.RevocationEndpoint,
	} {
		if *endpoint == "" {
			continue
		}
		ref, err := url.Parse(*endpoint)
		if err != nil {
			return nil, err
		}
		*endpoint = base.ResolveReference(ref).String()
	}

	return &meta, nil
}

type registrationRequest struct {
	ClientName      string   `json:"client_name"`
	ClientURI       string   `json:"client_uri"`
	RedirectURIs    []string `json:"redirect_uris"`
	Scope           string   `json:"scope"`
	SoftwareID      string   `json:"software_id"`
	SoftwareVersion string   `json:"software_version,omitempty"`
}

type registrationResponse struct {
	ClientID string `json:"client_id"`
}

// Register sends no redirect URIs since only the device flow is used.
func Register(ctx context.Context, client *http.Client, endpoint, version string) (string, error) {
	body, err := json.Marshal(registrationRequest{
		ClientName:      ClientName,
		ClientURI:       ClientURI,
		RedirectURIs:    []string{},
		Scope:           Scope,
		SoftwareID:      SoftwareID,
		SoftwareVersion: version,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", fmt.Errorf("%w: %d %s", ErrRegistration, resp.StatusCode, bytes.TrimSpace(msg))
	}

	var reg registrationResponse
	if err := json.NewDecoder(resp.Body).Decode(&reg); err != nil {
		return "", err
	}
	if reg.ClientID == "" {
		return "", fmt.Errorf("%w: no client_id returned", ErrRegistration)
	}
	return reg.ClientID, nil
}

func oauthConfig(clientID string, meta *Metadata) *oauth2.Config {
	return &oauth2.Config{
		ClientID: clientID,
		Scopes:   []string{Scope},
		Endpoint: oauth2.Endpoint{
			TokenURL:      meta.TokenEndpoint,
			DeviceAuthURL: meta.DeviceAuthorizationEndpoint,
			AuthStyle:     oauth2.AuthStyleInParams,
		},
	}
}
