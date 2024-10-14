//nolint:gochecknoglobals,gosmopolitan
package testproperties

import (
	_ "embed"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/nightscout"
)

var (
	//go:embed fetch_test_properties.json
	JSON []byte

	Properties = &nightscout.Properties{
		Bgnow: nightscout.Reading{
			Mean:      "123",
			Last:      123,
			Mills:     nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 30, 58, 417000000, time.Local)},
			Index:     "",
			FromMills: nightscout.Mills{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
			ToMills:   nightscout.Mills{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
			Sgvs: []nightscout.SGV{{
				ID:         "633a49639fc610138697ba4d",
				Device:     "xDrip-DexcomG5",
				Direction:  "Flat",
				Filtered:   "0",
				Mgdl:       123,
				Mills:      nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 30, 58, 417000000, time.Local)},
				Noise:      "1",
				Rssi:       "100",
				Scaled:     "123",
				Type:       "sgv",
				Unfiltered: "0",
			}},
		},
		Buckets: []nightscout.Reading{
			{
				Mean:      "123",
				Last:      123,
				Mills:     nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 30, 58, 417000000, time.Local)},
				Index:     "0",
				FromMills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 28, 28, 417000000, time.Local)},
				ToMills:   nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 33, 28, 417000000, time.Local)},
				Sgvs: []nightscout.SGV{{
					ID:         "633a49639fc610138697ba4d",
					Device:     "xDrip-DexcomG5",
					Direction:  "Flat",
					Filtered:   "0",
					Mgdl:       123,
					Mills:      nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 30, 58, 417000000, time.Local)},
					Noise:      "1",
					Rssi:       "100",
					Scaled:     "123",
					Type:       "sgv",
					Unfiltered: "0",
				}},
			},
			{
				Mean:      "122",
				Last:      122,
				Mills:     nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 25, 59, 159000000, time.Local)},
				Index:     "1",
				FromMills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 23, 28, 417000000, time.Local)},
				ToMills:   nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 28, 28, 417000000, time.Local)},
				Sgvs: []nightscout.SGV{{
					ID:         "633a48389fc610138697b95b",
					Device:     "xDrip-DexcomG5",
					Direction:  "Flat",
					Filtered:   "0",
					Mgdl:       122,
					Mills:      nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 25, 59, 159000000, time.Local)},
					Noise:      "1",
					Rssi:       "100",
					Scaled:     "122",
					Type:       "sgv",
					Unfiltered: "0",
				}},
			},
			{
				Mean:      "119",
				Last:      119,
				Mills:     nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 20, 59, 528000000, time.Local)},
				Index:     "2",
				FromMills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 18, 28, 417000000, time.Local)},
				ToMills:   nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 23, 28, 417000000, time.Local)},
				Sgvs: []nightscout.SGV{{
					ID:         "633a470d9fc610138697b86a",
					Device:     "xDrip-DexcomG5",
					Direction:  "Flat",
					Filtered:   "0",
					Mgdl:       119,
					Mills:      nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 20, 59, 528000000, time.Local)},
					Noise:      "1",
					Rssi:       "100",
					Scaled:     "119",
					Type:       "sgv",
					Unfiltered: "0",
				}},
			},
			{
				Mean:      "116",
				Last:      116,
				Mills:     nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 15, 59, 256000000, time.Local)},
				Index:     "3",
				FromMills: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 13, 28, 417000000, time.Local)},
				ToMills:   nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 18, 28, 417000000, time.Local)},
				Sgvs: []nightscout.SGV{{
					ID:         "633a45e09fc610138697b779",
					Device:     "xDrip-DexcomG5",
					Direction:  "Flat",
					Filtered:   "0",
					Mgdl:       116,
					Mills:      nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 15, 59, 256000000, time.Local)},
					Noise:      "1",
					Rssi:       "100",
					Scaled:     "116",
					Type:       "sgv",
					Unfiltered: "0",
				}},
			},
		},
		Delta: nightscout.Delta{
			Absolute:     "1",
			DisplayVal:   "+1",
			ElapsedMins:  "4.987633333333333",
			Interpolated: false,
			Mean5MinsAgo: "122",
			Mgdl:         "1",
			Previous: nightscout.Reading{
				Mean:      "122",
				Last:      122,
				Mills:     nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 25, 59, 159000000, time.Local)},
				Index:     "",
				FromMills: nightscout.Mills{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
				ToMills:   nightscout.Mills{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
				Sgvs: []nightscout.SGV{{
					ID:         "633a48389fc610138697b95b",
					Device:     "xDrip-DexcomG5",
					Direction:  "Flat",
					Filtered:   "0",
					Mgdl:       122,
					Mills:      nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 25, 59, 159000000, time.Local)},
					Noise:      "1",
					Rssi:       "100",
					Scaled:     "122",
					Type:       "sgv",
					Unfiltered: "0",
				}},
			},
			Scaled: 1,
			Times: nightscout.Times{
				Previous: nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 25, 59, 159000000, time.Local)},
				Recent:   nightscout.Mills{Time: time.Date(2022, time.October, 2, 21, 30, 58, 417000000, time.Local)},
			},
		},
		Direction: nightscout.Direction{
			Entity: "&#8594;",
			Label:  "→",
			Value:  "Flat",
		},
	}

	Etag = `W/"20-8b9f9edb2e2b1a9f5a8ffbf92a1a1c42f170a654"`
)
