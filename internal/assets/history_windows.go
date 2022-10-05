package assets

import _ "embed"

//go:generate ./convert-icon.sh src/rectangle-history-solid.svg dist/rectangle-history-solid.ico
//go:embed dist/rectangle-history-solid.ico
var History []byte
