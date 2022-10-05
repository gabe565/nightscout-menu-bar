//go:build !windows

package assets

import _ "embed"

//go:generate ./convert-icon.sh src/gear-solid.svg dist/gear-solid.png
//go:embed dist/gear-solid.png
var Gear []byte
