package assets

import _ "embed"

//go:generate ./convert-icon.sh src/triangle-exclamation-solid.svg dist/triangle-exclamation-solid.ico
//go:embed dist/triangle-exclamation-solid.ico
var Error []byte
