//go:build !windows

package assets

import _ "embed"

//go:generate ./convert-icon.sh src/github-brands-solid.svg dist/github-brands-solid.png
//go:embed dist/github-brands-solid.png
var About []byte
