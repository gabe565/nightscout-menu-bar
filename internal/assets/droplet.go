//go:build !windows

package assets

import _ "embed"

//go:generate ./convert-icon.sh src/droplet-solid.svg dist/droplet-solid.png
//go:embed dist/droplet-solid.png
var Droplet []byte
