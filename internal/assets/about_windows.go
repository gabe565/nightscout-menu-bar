package assets

import _ "embed"

//go:generate ./convert-icon.sh src/github-brands-solid.svg dist/github-brands-solid.ico
//go:embed dist/github-brands-solid.ico
var About []byte
