//go:build !windows

package assets

import _ "embed"

//go:generate ./convert-icon.sh src/reading.svg dist/reading.png
//go:embed dist/reading.png
var Reading []byte
