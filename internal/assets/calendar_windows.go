package assets

import _ "embed"

//go:generate ./convert-icon.sh src/calendar-solid.svg dist/calendar-solid.ico
//go:embed dist/calendar-solid.ico
var Calendar []byte
