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

package gl

import (
	"github.com/thinkofdeath/gl/v3.2-core/gl"
)

// VertexArray is used to store state needed to render vertices.
// This includes buffers, the format of the buffers and enabled
// attributes.
type VertexArray struct {
	internal uint32
}

var currentVertexArray VertexArray

// CreateVertexArray allocates a new VertexArray. If the allocation
// fails then IsValid will return false.
func CreateVertexArray() VertexArray {
	var va VertexArray
	gl.GenVertexArrays(1, &va.internal)
	return va
}

// Bind marks the VertexArray as the currently active one, this
// means buffers/the format of the buffers etc will be bound to
// this VertexArray. If this vertex array is already bound then
// this will do nothing.
func (va VertexArray) Bind() {
	if currentVertexArray == va {
		return
	}
	gl.BindVertexArray(va.internal)
	currentVertexArray = va
}

// Delete deallocates the VertexArray. This does not free any
// attached buffers. IsValid will return false after this call.
func (va *VertexArray) Delete() {
	gl.DeleteVertexArrays(1, &va.internal)
	if currentVertexArray == *va {
		currentVertexArray = VertexArray{}
	}
	va.internal = 0
}

// IsValid returns whether this VertexArray is still valid. A
// VertexArray will become invalid after Delete is called.
func (va VertexArray) IsValid() bool {
	return va.internal != 0
}
