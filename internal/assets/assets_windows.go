package assets

import _ "embed"

//go:generate ./convert-icon.sh icon-menu-bar.svg icon-menu-bar.ico
//go:embed icon-menu-bar.ico
var IconMenuBar []byte

//go:generate ./convert-icon.sh calendar-solid.svg calendar-solid.ico
//go:embed calendar-solid.ico
var Calendar []byte

//go:generate ./convert-icon.sh rectangle-history-solid.svg rectangle-history-solid.ico
//go:embed rectangle-history-solid.ico
var RectangleHistory []byte

//go:generate ./convert-icon.sh square-up-right-solid.svg square-up-right-solid.ico
//go:embed square-up-right-solid.ico
var SquareUpRight []byte

//go:generate ./convert-icon.sh gear-solid.svg gear-solid.ico
//go:embed gear-solid.ico
var Gear []byte

//go:generate ./convert-icon.sh xmark-solid.svg xmark-solid.ico
//go:embed xmark-solid.ico
var Xmark []byte
