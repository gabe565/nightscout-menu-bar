package assets

import _ "embed"

//go:generate ./convert-icon.sh src/error.svg dist/error.ico
//go:embed dist/error.ico
var Error []byte
