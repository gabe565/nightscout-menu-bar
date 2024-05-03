package fetch

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/hhsnopek/etag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed fetch_test_properties.json
var propertiesJson []byte

var (
	properties     = &nightscout.Properties{Bgnow: nightscout.Reading{Mean: 123, Last: 123, Mills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 30, 58, 417000000, time.Local)}, Index: 0, FromMills: nightscout.Mills{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}, ToMills: nightscout.Mills{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}, Sgvs: []nightscout.SGV{{ID: "633a49639fc610138697ba4d", Device: "xDrip-DexcomG5", Direction: "Flat", Filtered: 0, Mgdl: 123, Mills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 30, 58, 417000000, time.Local)}, Noise: 1, Rssi: 100, Scaled: "123", Type: "sgv", Unfiltered: 0}}}, Buckets: []nightscout.Reading{{Mean: 123, Last: 123, Mills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 30, 58, 417000000, time.Local)}, Index: 0, FromMills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 28, 28, 417000000, time.Local)}, ToMills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 33, 28, 417000000, time.Local)}, Sgvs: []nightscout.SGV{{ID: "633a49639fc610138697ba4d", Device: "xDrip-DexcomG5", Direction: "Flat", Filtered: 0, Mgdl: 123, Mills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 30, 58, 417000000, time.Local)}, Noise: 1, Rssi: 100, Scaled: "123", Type: "sgv", Unfiltered: 0}}}, {Mean: 122, Last: 122, Mills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 25, 59, 159000000, time.Local)}, Index: 1, FromMills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 23, 28, 417000000, time.Local)}, ToMills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 28, 28, 417000000, time.Local)}, Sgvs: []nightscout.SGV{{ID: "633a48389fc610138697b95b", Device: "xDrip-DexcomG5", Direction: "Flat", Filtered: 0, Mgdl: 122, Mills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 25, 59, 159000000, time.Local)}, Noise: 1, Rssi: 100, Scaled: "122", Type: "sgv", Unfiltered: 0}}}, {Mean: 119, Last: 119, Mills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 20, 59, 528000000, time.Local)}, Index: 2, FromMills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 18, 28, 417000000, time.Local)}, ToMills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 23, 28, 417000000, time.Local)}, Sgvs: []nightscout.SGV{{ID: "633a470d9fc610138697b86a", Device: "xDrip-DexcomG5", Direction: "Flat", Filtered: 0, Mgdl: 119, Mills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 20, 59, 528000000, time.Local)}, Noise: 1, Rssi: 100, Scaled: "119", Type: "sgv", Unfiltered: 0}}}, {Mean: 116, Last: 116, Mills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 15, 59, 256000000, time.Local)}, Index: 3, FromMills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 13, 28, 417000000, time.Local)}, ToMills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 18, 28, 417000000, time.Local)}, Sgvs: []nightscout.SGV{{ID: "633a45e09fc610138697b779", Device: "xDrip-DexcomG5", Direction: "Flat", Filtered: 0, Mgdl: 116, Mills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 15, 59, 256000000, time.Local)}, Noise: 1, Rssi: 100, Scaled: "116", Type: "sgv", Unfiltered: 0}}}}, Delta: nightscout.Delta{Absolute: 1, DisplayVal: "+1", ElapsedMins: 4.987633333333333, Interpolated: false, Mean5MinsAgo: 122, Mgdl: 1, Previous: nightscout.Reading{Mean: 122, Last: 122, Mills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 25, 59, 159000000, time.Local)}, Index: 0, FromMills: nightscout.Mills{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}, ToMills: nightscout.Mills{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}, Sgvs: []nightscout.SGV{{ID: "633a48389fc610138697b95b", Device: "xDrip-DexcomG5", Direction: "Flat", Filtered: 0, Mgdl: 122, Mills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 25, 59, 159000000, time.Local)}, Noise: 1, Rssi: 100, Scaled: "122", Type: "sgv", Unfiltered: 0}}}, Scaled: 1, Times: nightscout.Times{Previous: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 25, 59, 159000000, time.Local)}, Recent: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 30, 58, 417000000, time.Local)}}}, Direction: nightscout.Direction{Entity: "&#8594;", Label: "→", Value: "Flat"}}
	propertiesEtag = `W/"20-8b9f9edb2e2b1a9f5a8ffbf92a1a1c42f170a654"`
	differentEtag  = `W/"7cb-pLFn++MnPyzFsPH7e6MXNVFr2KU"`
)

func TestFetch_Do(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		etag := etag.Generate(propertiesJson, true)

		if reqEtag := r.Header.Get("If-None-Match"); reqEtag == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		w.Header().Set("Etag", etag)
		_, _ = w.Write(propertiesJson)
	}))
	defer server.Close()

	type fields struct {
		config        *config.Config
		url           *url.URL
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
			properties,
			propertiesEtag,
			require.NoError,
		},
		{
			"same etag",
			fields{config: &config.Config{URL: server.URL}, etag: propertiesEtag},
			args{context.Background()},
			nil,
			propertiesEtag,
			require.Error,
		},
		{
			"different etag",
			fields{config: &config.Config{URL: server.URL}, etag: differentEtag},
			args{context.Background()},
			properties,
			propertiesEtag,
			require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFetch(tt.fields.config)
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
