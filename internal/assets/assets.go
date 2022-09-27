package assets

import _ "embed"

//go:embed icon-32.png
var Icon32 []byte

//go:generate ./convert-images.sh calendar-solid.svg
//go:embed calendar-solid.png
var Calendar []byte

//go:generate ./convert-images.sh rectangle-history-solid.svg
//go:embed rectangle-history-solid.png
var RectangleHistory []byte

//go:generate ./convert-images.sh square-up-right-solid.svg
//go:embed square-up-right-solid.png
var SquareUpRight []byte
