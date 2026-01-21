package assets

import _ "embed"

//go:generate ./convert-icon.sh src/open.svg dist/open.ico
//go:embed dist/open.ico
var Open []byte
