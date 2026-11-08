package assets

import _ "embed"

//go:generate ./convert-icon.sh src/reading.svg dist/reading.ico
//go:embed dist/reading.ico
var Reading []byte

//go:generate ./convert-icon.sh src/reading.svg dist/reading-white.ico 8 #fff
//go:embed dist/reading-white.ico
var ReadingWhite []byte
