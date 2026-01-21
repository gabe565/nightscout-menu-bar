package assets

import _ "embed"

//go:generate ./convert-icon.sh src/preferences.svg dist/preferences.ico
//go:embed dist/preferences.ico
var Preferences []byte
