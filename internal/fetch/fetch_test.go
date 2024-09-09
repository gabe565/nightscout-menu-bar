package fetch

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout/testproperties"
	"github.com/hhsnopek/etag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFetch(t *testing.T) {
	t.Parallel()
	fetch := NewFetch(config.New(), "")
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

		w.Header().Set("Etag", etag)
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
			args{context.Background()},
			nil,
			"",
			require.Error,
		},
		{
			"success",
			fields{config: &config.Config{URL: server.URL}},
			args{context.Background()},
			testproperties.Properties,
			testproperties.Etag,
			require.NoError,
		},
		{
			"same etag",
			fields{config: &config.Config{URL: server.URL}, etag: testproperties.Etag},
			args{context.Background()},
			nil,
			testproperties.Etag,
			require.Error,
		},
		{
			"different etag",
			fields{config: &config.Config{URL: server.URL}, etag: etag.Generate([]byte("test"), true)},
			args{context.Background()},
			testproperties.Properties,
			testproperties.Etag,
			require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := NewFetch(tt.fields.config, "")
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
