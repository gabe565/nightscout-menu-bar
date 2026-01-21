package assets

import _ "embed"

//go:generate ./convert-icon.sh src/quit.svg dist/quit.ico
//go:embed dist/quit.ico
var Quit []byte
