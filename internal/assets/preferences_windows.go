package assets

import _ "embed"

//go:generate ./convert-icon.sh src/preferences.svg dist/preferences.ico
//go:embed dist/preferences.ico
var Preferences []byte

//go:generate ./convert-icon.sh src/preferences.svg dist/preferences-white.ico 8 #fff
//go:embed dist/preferences-white.ico
var PreferencesWhite []byte
