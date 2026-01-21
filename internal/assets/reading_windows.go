package assets

import _ "embed"

//go:generate ./convert-icon.sh src/reading.svg dist/reading.ico
//go:embed dist/reading.ico
var Reading []byte
