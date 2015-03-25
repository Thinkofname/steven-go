package platform

import (
	"encoding/binary"
	"unsafe"
)

// Init sets the handler and starts platform specific code.
// This method blocks until the program ends.
func Init(handler Handler) {
	run(handler)
}

// Handler contains methods for handling platform events.
type Handler struct {
	Start func()
	Draw  func()

	Rotate func(x, y float64)
	Move   func(f, s float64)
}

// Size returns the size of the screen in pixels.
func Size() (width, height int) {
	return size()
}

// NativeOrder is the native byte order of the system
var NativeOrder binary.ByteOrder

func init() {
	check := uint32(1)
	c := (*[4]byte)(unsafe.Pointer(&check))
	NativeOrder = binary.BigEndian
	if binary.LittleEndian.Uint32(c[:]) == 1 {
		NativeOrder = binary.LittleEndian
	}
}
