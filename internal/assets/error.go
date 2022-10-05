//go:build !windows

package assets

import _ "embed"

//go:generate ./convert-icon.sh src/triangle-exclamation-solid.svg dist/triangle-exclamation-solid.png
//go:embed dist/triangle-exclamation-solid.png
var Error []byte
