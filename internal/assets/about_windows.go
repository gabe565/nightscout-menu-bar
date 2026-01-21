package assets

import _ "embed"

//go:generate ./convert-icon.sh src/about.svg dist/about.ico
//go:embed dist/about.ico
var About []byte
