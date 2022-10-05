package assets

import _ "embed"

//go:generate ./convert-icon.sh src/square-up-right-solid.svg dist/square-up-right-solid.ico
//go:embed dist/square-up-right-solid.ico
var Open []byte
