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

package nibble

// Array is an array of 4 bit values.
type Array []byte

// New creates a new nibble array that can store at least
// the requested number of values.
func New(size int) []byte {
	return make(Array, (size+1)>>1)
}

// Get returns the value at the index.
func (a Array) Get(idx int) byte {
	val := a[idx>>1]
	if idx&1 == 0 {
		return val & 0xF
	}
	return val >> 4
}

// Set sets the value a the index.
func (a Array) Set(idx int, val byte) {
	i := idx >> 1
	o := a[i]
	if idx&1 == 0 {
		a[i] = (o & 0xF0) | (val & 0xF)
	} else {
		a[i] = (o & 0x0F) | ((val & 0xF) << 4)
	}
}
