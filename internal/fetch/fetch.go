package fetch

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/rs/zerolog/log"
)

var (
	ErrHTTP        = errors.New("unexpected HTTP error")
	ErrNotModified = errors.New("not modified")
	ErrNoURL       = errors.New("please configure your Nightscout URL")
)

func NewFetch(conf *config.Config) *Fetch {
	return &Fetch{
		config: conf,
	}
}

type Fetch struct {
	config        *config.Config
	url           string
	tokenChecksum string
	etag          string
}

func (f *Fetch) Do(ctx context.Context) (*nightscout.Properties, error) {
	if f.url == "" {
		if err := f.UpdateURL(); err != nil {
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

	log.Debug().
		Bool("etag", f.etag != "").
		Bool("secret", f.tokenChecksum != "").
		Msg("Fetching data")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Trace().RawJSON("data", data).Msg("Response data")

	switch resp.StatusCode {
	case http.StatusNotModified:
		log.Trace().Msg("Data was not modified")
		return nil, ErrNotModified
	case http.StatusOK:
		// Decode JSON
		var properties nightscout.Properties
		if err := json.Unmarshal(data, &properties); err != nil {
			return nil, err
		}

		log.Trace().Msg("Parsed response")

		f.etag = resp.Header.Get("etag")
		return &properties, nil
	default:
		f.etag = ""
		return nil, fmt.Errorf("%w: %d", ErrHTTP, resp.StatusCode)
	}
}

func (f *Fetch) UpdateURL() error {
	u, err := BuildURL(f.config)
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, "api", "v2", "properties", "bgnow,buckets,delta,direction")
	f.url = u.String()
	log.Debug().Str("value", f.url).Msg("Generated URL")

	if token := f.config.Token; token != "" {
		rawChecksum := sha1.Sum([]byte(token))
		f.tokenChecksum = hex.EncodeToString(rawChecksum[:])
		log.Trace().Str("value", f.tokenChecksum).Msg("Generated token checksum")
	} else {
		f.tokenChecksum = ""
	}

	return nil
}

func (f *Fetch) Reset() {
	log.Debug().Msg("Resetting fetch cache")
	f.url = ""
	f.tokenChecksum = ""
	f.etag = ""
}

func BuildURL(conf *config.Config) (*url.URL, error) {
	if conf.URL == "" {
		return nil, ErrNoURL
	}

	return url.Parse(conf.URL)
}

func BuildURLWithToken(conf *config.Config) (*url.URL, error) {
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
