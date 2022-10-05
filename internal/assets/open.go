//go:build !windows

package assets

import _ "embed"

//go:generate ./convert-icon.sh src/square-up-right-solid.svg dist/square-up-right-solid.png
//go:embed dist/square-up-right-solid.png
var Open []byte
