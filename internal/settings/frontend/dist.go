package frontend

import (
	"embed"
	"io/fs"
)

//go:generate go tool wails3 generate bindings -silent -ts -i -names -d src/bindings ..
//go:generate pnpm install
//go:generate pnpm run build

//go:embed all:dist
var dist embed.FS

func FS() (fs.FS, error) {
	return fs.Sub(dist, "dist")
}
