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
	"github.com/go-gl/gl/v3.2-core/gl"
)

type VertexArray struct {
	internal uint32
}

var currentVertexArray VertexArray

func CreateVertexArray() VertexArray {
	var va VertexArray
	gl.GenVertexArrays(1, &va.internal)
	if va.internal == 0 {
		panic("failed to create vertex array")
	}
	return va
}

func (va VertexArray) Bind() {
	if currentVertexArray == va {
		return
	}
	gl.BindVertexArray(va.internal)
	currentVertexArray = va
}

func (va VertexArray) Delete() {
	gl.DeleteVertexArrays(1, &va.internal)
	if currentVertexArray == va {
		currentVertexArray = VertexArray{}
	}
}
