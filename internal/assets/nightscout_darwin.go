package assets

import _ "embed"

//go:generate ./convert-icon.sh src/nightscout-transparent.svg dist/nightscout-transparent.png 2
//go:embed dist/nightscout-transparent.png
var Nightscout []byte
