//go:build !windows

package assets

import _ "embed"

//go:generate ./convert-icon.sh src/xmark-solid.svg dist/xmark-solid.png
//go:embed dist/xmark-solid.png
var Quit []byte
