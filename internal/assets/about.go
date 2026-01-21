//go:build !windows

package assets

import _ "embed"

//go:generate ./convert-icon.sh src/about.svg dist/about.png
//go:embed dist/about.png
var About []byte
