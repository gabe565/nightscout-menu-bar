// Package assets embeds the icons generated from src by generate.sh.
package assets

import _ "embed"

//go:generate ./generate.sh

// App is the full application icon, used for windows.
//
//go:embed dist/app.png
var App []byte
