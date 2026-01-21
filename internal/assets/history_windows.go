package assets

import _ "embed"

//go:generate ./convert-icon.sh src/history.svg dist/history.ico
//go:embed dist/history.ico
var History []byte
