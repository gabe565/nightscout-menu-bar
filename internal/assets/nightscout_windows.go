package assets

import _ "embed"

//go:generate ./convert-icon.sh src/nightscout.svg dist/nightscout.ico 0
//go:embed dist/nightscout.ico
var Nightscout []byte
