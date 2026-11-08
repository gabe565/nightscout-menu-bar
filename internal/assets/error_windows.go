package assets

import _ "embed"

//go:generate ./convert-icon.sh src/error.svg dist/error.ico
//go:embed dist/error.ico
var Error []byte

//go:generate ./convert-icon.sh src/error.svg dist/error-white.ico 8 #fff
//go:embed dist/error-white.ico
var ErrorWhite []byte
