package assets

import _ "embed"

//go:generate ./convert-icon.sh src/history.svg dist/history.ico
//go:embed dist/history.ico
var History []byte

//go:generate ./convert-icon.sh src/history.svg dist/history-white.ico 8 #fff
//go:embed dist/history-white.ico
var HistoryWhite []byte
