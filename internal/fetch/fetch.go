package fetch

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
)

var (
	ErrHttp        = errors.New("unexpected HTTP error")
	ErrNotModified = errors.New("not modified")
)

func NewFetch(conf *config.Config) *Fetch {
	return &Fetch{
		config: conf,
		client: &http.Client{
			Timeout: time.Minute,
		},
	}
}

type Fetch struct {
	config        *config.Config
	client        *http.Client
	url           *url.URL
	tokenChecksum string
	etag          string
}

func (f *Fetch) Do() (*nightscout.Properties, error) {
	if f.url == nil {
		if err := f.UpdateUrl(); err != nil {
			return nil, err
		}
	}

	// Fetch JSON
	req, err := http.NewRequest("GET", f.url.String(), nil)
	if err != nil {
		return nil, err
	}
	if f.etag != "" {
		req.Header.Set("If-None-Match", f.etag)
	}

	if f.tokenChecksum != "" {
		req.Header.Set("Api-Secret", f.tokenChecksum)
	}

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
		return nil, ErrNotModified
	case http.StatusOK:
		// Decode JSON
		var properties nightscout.Properties
		if err := json.NewDecoder(resp.Body).Decode(&properties); err != nil {
			return nil, err
		}

		f.etag = resp.Header.Get("etag")
		return &properties, nil
	default:
		f.etag = ""
		return nil, fmt.Errorf("%w: %d", ErrHttp, resp.StatusCode)
	}
}

func (f *Fetch) UpdateUrl() error {
	u, err := BuildUrl(f.config)
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, "api", "v2", "properties", "bgnow,buckets,delta,direction")
	f.url = u

	if token := f.config.Token; token != "" {
		rawChecksum := sha1.Sum([]byte(token))
		f.tokenChecksum = hex.EncodeToString(rawChecksum[:])
	} else {
		f.tokenChecksum = ""
	}

	return nil
}

func (f *Fetch) Reset() {
	f.url = nil
	f.tokenChecksum = ""
	f.etag = ""
}

func BuildUrl(conf *config.Config) (*url.URL, error) {
	if conf.URL == "" {
		return nil, errors.New("please configure your Nightscout URL")
	}

	return url.Parse(conf.URL)
}

func BuildUrlWithToken(conf *config.Config) (*url.URL, error) {
	u, err := BuildUrl(conf)
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
