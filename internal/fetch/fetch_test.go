package fetch

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/nightscout"
	"gabe565.com/nightscout-menu-bar/internal/nightscout/testproperties"
	"gabe565.com/nightscout-menu-bar/internal/oauth"
	"gabe565.com/nightscout-menu-bar/internal/secret"
	"github.com/hhsnopek/etag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFetch(t *testing.T) {
	t.Parallel()
	fetch := NewFetch(config.New(), nil)
	require.NotNil(t, fetch)
	assert.NotNil(t, fetch.config)
}

func TestFetch_Do(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		etag := etag.Generate(testproperties.JSON, true)

		if reqEtag := r.Header.Get("If-None-Match"); reqEtag == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		w.Header().Set("ETag", etag)
		_, _ = w.Write(testproperties.JSON)
	}))
	t.Cleanup(server.Close)

	type fields struct {
		config        *config.Config
		url           string
		tokenChecksum string
		etag          string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     *nightscout.Properties
		wantEtag string
		wantErr  require.ErrorAssertionFunc
	}{
		{
			"no url",
			fields{config: &config.Config{}},
			args{t.Context()},
			nil,
			"",
			require.Error,
		},
		{
			"success",
			fields{config: config.New(config.WithData(config.Data{URL: server.URL}))},
			args{t.Context()},
			testproperties.Properties,
			testproperties.Etag,
			require.NoError,
		},
		{
			"same etag",
			fields{config: config.New(config.WithData(config.Data{URL: server.URL})), etag: testproperties.Etag},
			args{t.Context()},
			nil,
			testproperties.Etag,
			require.Error,
		},
		{
			"different etag",
			fields{
				config: config.New(config.WithData(config.Data{URL: server.URL})),
				etag:   etag.Generate([]byte("test"), true),
			},
			args{t.Context()},
			testproperties.Properties,
			testproperties.Etag,
			require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := NewFetch(tt.fields.config, nil)
			f.url = tt.fields.url
			f.tokenChecksum = tt.fields.tokenChecksum
			f.etag = tt.fields.etag

			got, err := f.Do(tt.args.ctx)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantEtag, f.etag)
		})
	}
}

func TestFetch_Do_OAuth(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer access-1" || r.Header.Get("Api-Secret") != "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, _ = w.Write(testproperties.JSON)
	}))
	t.Cleanup(server.Close)

	dir := t.TempDir()
	creds, err := json.Marshal(oauth.Credentials{
		URL:      server.URL,
		ClientID: "client-1",
		Expiry:   time.Now().Add(time.Hour),
	})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "oauth.json"), creds, 0o600))

	conf := config.New(config.WithData(config.Data{URL: server.URL, Token: "ignored-when-connected"}))
	conf.File = filepath.Join(dir, "config.toml")
	conf.Secrets = secret.NewWithBackend(dir, &secret.Memory{})
	require.NoError(t, conf.Secrets.Set("nocturne", `{"access_token":"access-1","refresh_token":"refresh-1"}`))
	store := oauth.NewStore(conf)
	require.NoError(t, store.Load())

	got, err := NewFetch(conf, store).Do(t.Context())
	require.NoError(t, err)
	assert.Equal(t, testproperties.Properties, got)
}
