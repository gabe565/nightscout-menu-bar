//go:build !windows

package assets

import _ "embed"

//go:generate ./convert-icon.sh src/error.svg dist/error.png
//go:embed dist/error.png
var Error []byte
