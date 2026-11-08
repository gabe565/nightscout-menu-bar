package systheme

import "image/color"

// IconColor contrasts with the taskbar or menu bar.
func IconColor() color.Color {
	if Light() {
		return color.Black
	}
	return color.White
}
