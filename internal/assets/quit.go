//go:build !windows

package assets

import _ "embed"

//go:generate ./convert-icon.sh src/quit.svg dist/quit.png
//go:embed dist/quit.png
var Quit []byte
