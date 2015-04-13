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

package render

import "github.com/thinkofdeath/steven/render/gl"

var lineState = struct {
	program  gl.Program
	shader   *lineShader
	array    gl.VertexArray
	buffer   gl.Buffer
	count    int
	data     []byte
	prevSize int
}{
	prevSize: -1,
}

func initLineDraw() {
	lineState.program = CreateProgram(vertexLine, fragmentLine)
	lineState.shader = &lineShader{}
	InitStruct(lineState.shader, lineState.program)

	lineState.array = gl.CreateVertexArray()
	lineState.array.Bind()
	lineState.buffer = gl.CreateBuffer()
	lineState.buffer.Bind(gl.ArrayBuffer)
	lineState.shader.Position.Enable()
	lineState.shader.Color.Enable()
	lineState.shader.Position.Pointer(3, gl.Float, false, 16, 0)
	lineState.shader.Color.Pointer(4, gl.UnsignedByte, true, 16, 12)
}

func drawLines() {
	if lineState.count > 0 {
		lineState.program.Use()
		lineState.shader.PerspectiveMatrix.Matrix4(perspectiveMatrix)
		lineState.shader.CameraMatrix.Matrix4(cameraMatrix)
		lineState.array.Bind()
		lineState.buffer.Bind(gl.ArrayBuffer)
		if len(lineState.data) > lineState.prevSize {
			lineState.prevSize = len(lineState.data)
			lineState.buffer.Data(lineState.data, gl.DynamicDraw)
		} else {
			target := lineState.buffer.Map(gl.WriteOnly, len(lineState.data))
			copy(target, lineState.data)
			lineState.buffer.Unmap()
		}
		gl.DrawArrays(gl.Triangles, 0, lineState.count)
		lineState.count = 0
		lineState.data = lineState.data[:0]
	}
}

func DrawBox(x1, y1, z1, x2, y2, z2 float64, r, g, b, a byte) {
	for _, f := range faceVertices {
		for _, v := range f {
			val := v[0]*x2 + (1.0-v[0])*x1
			lineState.data = appendFloat(lineState.data, float32(val))
			val = v[1]*y2 + (1.0-v[1])*y1
			lineState.data = appendFloat(lineState.data, float32(val))
			val = v[2]*z2 + (1.0-v[2])*z1
			lineState.data = appendFloat(lineState.data, float32(val))
			lineState.data = append(lineState.data, r, g, b, a)
			lineState.count++
		}
	}
}

// Precomputed face vertices
var faceVertices = [6][6][3]float64{
	{ // Up
		{0, 1, 0},
		{1, 1, 0},
		{0, 1, 1},

		{1, 1, 1},
		{0, 1, 1},
		{1, 1, 0},
	},
	{ // Down
		{0, 0, 0},
		{0, 0, 1},
		{1, 0, 0},

		{1, 0, 1},
		{1, 0, 0},
		{0, 0, 1},
	},
	{ // North
		{0, 0, 0},
		{1, 0, 0},
		{0, 1, 0},

		{1, 1, 0},
		{0, 1, 0},
		{1, 0, 0},
	},
	{ // South
		{0, 0, 1},
		{0, 1, 1},
		{1, 0, 1},

		{1, 1, 1},
		{1, 0, 1},
		{0, 1, 1},
	},
	{ // West
		{0, 0, 0},
		{0, 1, 0},
		{0, 0, 1},

		{0, 1, 1},
		{0, 0, 1},
		{0, 1, 0},
	},
	{ // East
		{1, 0, 0},
		{1, 0, 1},
		{1, 1, 0},

		{1, 1, 1},
		{1, 1, 0},
		{1, 0, 1},
	},
}
