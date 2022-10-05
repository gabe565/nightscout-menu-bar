package assets

import _ "embed"

//go:generate ./convert-icon.sh src/xmark-solid.svg dist/xmark-solid.ico
//go:embed dist/xmark-solid.ico
var Quit []byte
