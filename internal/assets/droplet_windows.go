package assets

import _ "embed"

//go:generate ./convert-icon.sh src/droplet-solid.svg dist/droplet-solid.ico
//go:embed dist/droplet-solid.ico
var Droplet []byte
