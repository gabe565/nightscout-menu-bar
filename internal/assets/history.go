//go:build !windows

package assets

import _ "embed"

//go:generate ./convert-icon.sh src/history.svg dist/history.png
//go:embed dist/history.png
var History []byte
