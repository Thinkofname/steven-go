// Copyright 2015 Matthew Collins
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	Action func(action Action)
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
