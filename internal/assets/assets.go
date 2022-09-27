package assets

import _ "embed"

//go:generate ./convert-images.sh icon-menu-bar.svg
//go:embed icon-menu-bar.png
var IconMenuBar []byte

//go:generate ./convert-images.sh calendar-solid.svg
//go:embed calendar-solid.png
var Calendar []byte

//go:generate ./convert-images.sh rectangle-history-solid.svg
//go:embed rectangle-history-solid.png
var RectangleHistory []byte

//go:generate ./convert-images.sh square-up-right-solid.svg
//go:embed square-up-right-solid.png
var SquareUpRight []byte
