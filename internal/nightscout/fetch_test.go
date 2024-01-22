package nightscout

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/hhsnopek/etag"
	"github.com/spf13/viper"
)

//go:embed fetch_test_properties.json
var propertiesJson []byte

var (
	properties     = &Properties{Bgnow: Reading{Mean: 123, Last: 123, Mills: Mills{Time: time.Date(2022, time.October, 2, 21, 30, 58, 417000000, time.Local)}, Index: 0, FromMills: Mills{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}, ToMills: Mills{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}, Sgvs: []SGV{{ID: "633a49639fc610138697ba4d", Device: "xDrip-DexcomG5", Direction: "Flat", Filtered: 0, Mgdl: 123, Mills: Mills{Time: time.Date(2022, time.October, 2, 21, 30, 58, 417000000, time.Local)}, Noise: 1, Rssi: 100, Scaled: "123", Type: "sgv", Unfiltered: 0}}}, Buckets: []Reading{{Mean: 123, Last: 123, Mills: Mills{Time: time.Date(2022, time.October, 2, 21, 30, 58, 417000000, time.Local)}, Index: 0, FromMills: Mills{Time: time.Date(2022, time.October, 2, 21, 28, 28, 417000000, time.Local)}, ToMills: Mills{Time: time.Date(2022, time.October, 2, 21, 33, 28, 417000000, time.Local)}, Sgvs: []SGV{{ID: "633a49639fc610138697ba4d", Device: "xDrip-DexcomG5", Direction: "Flat", Filtered: 0, Mgdl: 123, Mills: Mills{Time: time.Date(2022, time.October, 2, 21, 30, 58, 417000000, time.Local)}, Noise: 1, Rssi: 100, Scaled: "123", Type: "sgv", Unfiltered: 0}}}, {Mean: 122, Last: 122, Mills: Mills{Time: time.Date(2022, time.October, 2, 21, 25, 59, 159000000, time.Local)}, Index: 1, FromMills: Mills{Time: time.Date(2022, time.October, 2, 21, 23, 28, 417000000, time.Local)}, ToMills: Mills{Time: time.Date(2022, time.October, 2, 21, 28, 28, 417000000, time.Local)}, Sgvs: []SGV{{ID: "633a48389fc610138697b95b", Device: "xDrip-DexcomG5", Direction: "Flat", Filtered: 0, Mgdl: 122, Mills: Mills{Time: time.Date(2022, time.October, 2, 21, 25, 59, 159000000, time.Local)}, Noise: 1, Rssi: 100, Scaled: "122", Type: "sgv", Unfiltered: 0}}}, {Mean: 119, Last: 119, Mills: Mills{Time: time.Date(2022, time.October, 2, 21, 20, 59, 528000000, time.Local)}, Index: 2, FromMills: Mills{Time: time.Date(2022, time.October, 2, 21, 18, 28, 417000000, time.Local)}, ToMills: Mills{Time: time.Date(2022, time.October, 2, 21, 23, 28, 417000000, time.Local)}, Sgvs: []SGV{{ID: "633a470d9fc610138697b86a", Device: "xDrip-DexcomG5", Direction: "Flat", Filtered: 0, Mgdl: 119, Mills: Mills{Time: time.Date(2022, time.October, 2, 21, 20, 59, 528000000, time.Local)}, Noise: 1, Rssi: 100, Scaled: "119", Type: "sgv", Unfiltered: 0}}}, {Mean: 116, Last: 116, Mills: Mills{Time: time.Date(2022, time.October, 2, 21, 15, 59, 256000000, time.Local)}, Index: 3, FromMills: Mills{Time: time.Date(2022, time.October, 2, 21, 13, 28, 417000000, time.Local)}, ToMills: Mills{Time: time.Date(2022, time.October, 2, 21, 18, 28, 417000000, time.Local)}, Sgvs: []SGV{{ID: "633a45e09fc610138697b779", Device: "xDrip-DexcomG5", Direction: "Flat", Filtered: 0, Mgdl: 116, Mills: Mills{Time: time.Date(2022, time.October, 2, 21, 15, 59, 256000000, time.Local)}, Noise: 1, Rssi: 100, Scaled: "116", Type: "sgv", Unfiltered: 0}}}}, Delta: Delta{Absolute: 1, DisplayVal: "+1", ElapsedMins: 4.987633333333333, Interpolated: false, Mean5MinsAgo: 122, Mgdl: 1, Previous: Reading{Mean: 122, Last: 122, Mills: Mills{Time: time.Date(2022, time.October, 2, 21, 25, 59, 159000000, time.Local)}, Index: 0, FromMills: Mills{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}, ToMills: Mills{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}, Sgvs: []SGV{{ID: "633a48389fc610138697b95b", Device: "xDrip-DexcomG5", Direction: "Flat", Filtered: 0, Mgdl: 122, Mills: Mills{Time: time.Date(2022, time.October, 2, 21, 25, 59, 159000000, time.Local)}, Noise: 1, Rssi: 100, Scaled: "122", Type: "sgv", Unfiltered: 0}}}, Scaled: 1, Times: Times{Previous: Mills{Time: time.Date(2022, time.October, 2, 21, 25, 59, 159000000, time.Local)}, Recent: Mills{Time: time.Date(2022, time.October, 2, 21, 30, 58, 417000000, time.Local)}}}, Direction: Direction{Entity: "&#8594;", Label: "→", Value: "Flat"}}
	propertiesEtag = `W/"20-8b9f9edb2e2b1a9f5a8ffbf92a1a1c42f170a654"`
	differentEtag  = `W/"7cb-pLFn++MnPyzFsPH7e6MXNVFr2KU"`
)

func TestFetch(t *testing.T) {
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

	tests := []struct {
		name     string
		url      string
		etag     string
		want     *Properties
		wantEtag string
		wantErr  bool
	}{
		{"no url", "", "", nil, "", true},
		{"success", server.URL, "", properties, propertiesEtag, false},
		{"same etag", server.URL, propertiesEtag, nil, propertiesEtag, true},
		{"different etag", server.URL, differentEtag, properties, propertiesEtag, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("url", tt.url)
			lastEtag = tt.etag

			got, err := Fetch()
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fetch() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(lastEtag, tt.wantEtag) {
				t.Errorf("Fetch() got etag = %v, want etag %v", lastEtag, tt.wantEtag)
			}
		})
	}
}
