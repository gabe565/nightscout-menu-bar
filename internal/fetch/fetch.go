package fetch

import (
	"context"
	"crypto/sha1" //nolint:gosec
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/nightscout"
	"gabe565.com/nightscout-menu-bar/internal/util"
)

var (
	ErrHTTP        = errors.New("unexpected HTTP error")
	ErrNotModified = errors.New("not modified")
	ErrNoURL       = errors.New("please configure your Nightscout URL")
)

func NewFetch(conf *config.Config) *Fetch {
	return &Fetch{
		config: conf,
		client: &http.Client{
			Transport: util.NewUserAgentTransport("nightscout-menu-bar", conf.Version),
			Timeout:   time.Minute,
		},
	}
}

type Fetch struct {
	mu            sync.Mutex
	config        *config.Config
	client        *http.Client
	url           string
	tokenChecksum string
	etag          string
}

func (f *Fetch) Do(ctx context.Context) (*nightscout.Properties, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	start := time.Now()

	if f.url == "" {
		if err := f.updateURLLocked(); err != nil {
			return nil, err
		}
	}

	// Fetch JSON
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.url, nil)
	if err != nil {
		return nil, err
	}
	if f.etag != "" {
		req.Header.Set("If-None-Match", f.etag)
	}

	if f.tokenChecksum != "" {
		req.Header.Set("Api-Secret", f.tokenChecksum)
	}

	slog.Debug("Fetching data",
		"etag", f.etag != "",
		"secret", f.tokenChecksum != "",
	)

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusNotModified:
		slog.Debug("Data was not modified", "took", time.Since(start))
		return nil, ErrNotModified
	case http.StatusOK:
		// Decode JSON
		var properties nightscout.Properties
		if err := json.NewDecoder(resp.Body).Decode(&properties); err != nil {
			return nil, err
		}

		slog.Debug("Parsed response", "took", time.Since(start), "data", properties)

		f.etag = resp.Header.Get("etag")
		return &properties, nil
	default:
		f.etag = ""
		return nil, fmt.Errorf("%w: %d", ErrHTTP, resp.StatusCode)
	}
}

func (f *Fetch) UpdateURL() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.updateURLLocked()
}

func (f *Fetch) updateURLLocked() error {
	data := f.config.Data()

	u, err := BuildURL(data)
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, "api", "v2", "properties", "bgnow,buckets,delta,direction")
	f.url = u.String()
	slog.Debug("Generated URL", "value", f.url)

	if token := data.Token; token != "" {
		rawChecksum := sha1.Sum([]byte(token)) //nolint:gosec
		f.tokenChecksum = hex.EncodeToString(rawChecksum[:])
		slog.Debug("Generated token checksum", "value", f.tokenChecksum)
	} else {
		f.tokenChecksum = ""
	}

	return nil
}

func (f *Fetch) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()

	slog.Debug("Resetting fetch cache")
	f.url = ""
	f.tokenChecksum = ""
	f.etag = ""
}

func BuildURL(conf config.Data) (*url.URL, error) {
	if conf.URL == "" {
		return nil, ErrNoURL
	}

	return url.Parse(conf.URL)
}

func BuildURLWithToken(conf config.Data) (*url.URL, error) {
	u, err := BuildURL(conf)
	if err != nil {
		return u, err
	}

	if token := conf.Token; token != "" {
		query := u.Query()
		query.Set("token", conf.Token)
		u.RawQuery = query.Encode()
	}

	return u, nil
}
