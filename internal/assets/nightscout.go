//go:build !(darwin || windows)

package assets

import _ "embed"

//go:generate ./convert-icon.sh src/nightscout.svg dist/nightscout.png
//go:embed dist/nightscout.png
var Nightscout []byte
