//go:build !windows

package assets

import _ "embed"

//go:generate ./convert-icon.sh src/open.svg dist/open.png
//go:embed dist/open.png
var Open []byte
