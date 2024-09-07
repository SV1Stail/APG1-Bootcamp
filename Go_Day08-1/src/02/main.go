package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "application.h"
#include "window.h"
*/

import "C"
import "unsafe"

func main() {

	C.InitApplication()
	title := C.CString("Школа 21")
	defer C.free(unsafe.Pointer(title))
	window := C.Window_Create(0, 0, 300, 200, title)
	C.Window_MakeKeyAndOrderFront(window)
	C.RunApplication()
}
