//go:build !darwin && !windows

package assets

import _ "embed"

//go:embed dist/nightscout.png
var Nightscout []byte

//go:embed dist/about.png
var About []byte

//go:embed dist/error.png
var Error []byte

//go:embed dist/history.png
var History []byte

//go:embed dist/open.png
var Open []byte

//go:embed dist/preferences.png
var Preferences []byte

//go:embed dist/quit.png
var Quit []byte

//go:embed dist/reading.png
var Reading []byte

//go:embed dist/update.png
var Update []byte

// White variants for dark menus.

//go:embed dist/about-white.png
var AboutWhite []byte

//go:embed dist/error-white.png
var ErrorWhite []byte

//go:embed dist/history-white.png
var HistoryWhite []byte

//go:embed dist/open-white.png
var OpenWhite []byte

//go:embed dist/preferences-white.png
var PreferencesWhite []byte

//go:embed dist/quit-white.png
var QuitWhite []byte

//go:embed dist/reading-white.png
var ReadingWhite []byte

//go:embed dist/update-white.png
var UpdateWhite []byte
