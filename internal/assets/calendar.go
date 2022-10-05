//go:build !windows

package assets

import _ "embed"

//go:generate ./convert-icon.sh src/calendar-solid.svg dist/calendar-solid.png
//go:embed dist/calendar-solid.png
var Calendar []byte
