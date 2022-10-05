//go:build !windows

package assets

import _ "embed"

//go:generate ./convert-icon.sh src/rectangle-history-solid.svg dist/rectangle-history-solid.png
//go:embed dist/rectangle-history-solid.png
var History []byte
