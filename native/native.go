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

package native

import (
	"encoding/binary"
	"unsafe"
)

// Order is the native byte order of the system
var Order binary.ByteOrder

func init() {
	check := uint32(1)
	c := (*[4]byte)(unsafe.Pointer(&check))
	Order = binary.BigEndian
	if binary.LittleEndian.Uint32(c[:]) == 1 {
		Order = binary.LittleEndian
	}
}
