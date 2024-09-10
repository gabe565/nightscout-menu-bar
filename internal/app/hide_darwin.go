package app

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void
SetActivationPolicy(void) {
    [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
}
*/
import "C"

func HideAppIcon() {
	C.SetActivationPolicy()
}
