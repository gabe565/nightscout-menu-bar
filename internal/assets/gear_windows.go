package assets

import _ "embed"

//go:generate ./convert-icon.sh src/gear-solid.svg dist/gear-solid.ico
//go:embed dist/gear-solid.ico
var Gear []byte
