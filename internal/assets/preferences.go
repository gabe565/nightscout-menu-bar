//go:build !windows

package assets

import _ "embed"

//go:generate ./convert-icon.sh src/preferences.svg dist/preferences.png
//go:embed dist/preferences.png
var Preferences []byte
