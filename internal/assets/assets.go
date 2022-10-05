package assets

import _ "embed"

//go:generate ./convert-icon.sh icon-menu-bar.svg icon-menu-bar.png
//go:embed icon-menu-bar.png
var IconMenuBar []byte

//go:generate ./convert-icon.sh calendar-solid.svg calendar-solid.png
//go:embed calendar-solid.png
var Calendar []byte

//go:generate ./convert-icon.sh rectangle-history-solid.svg rectangle-history-solid.png
//go:embed rectangle-history-solid.png
var RectangleHistory []byte

//go:generate ./convert-icon.sh square-up-right-solid.svg square-up-right-solid.png
//go:embed square-up-right-solid.png
var SquareUpRight []byte

//go:generate ./convert-icon.sh gear-solid.svg gear-solid.png
//go:embed gear-solid.png
var Gear []byte

//go:generate ./convert-icon.sh xmark-solid.svg xmark-solid.png
//go:embed xmark-solid.png
var Xmark []byte
