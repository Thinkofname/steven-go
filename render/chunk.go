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

import (
	"github.com/thinkofdeath/steven/platform/gl"
	"github.com/thinkofdeath/steven/type/direction"
)

var (
	buffers       = make(map[position]*ChunkBuffer)
	bufferColumns = make(map[positionC]int)
)

// ChunkBuffer is a renderable chunk section
type ChunkBuffer struct {
	position
	invalid bool

	array    gl.VertexArray
	buffer   gl.Buffer
	count    int
	arrayT   gl.VertexArray
	bufferT  gl.Buffer
	countT   int
	cullBits uint64

	renderedOn uint
}

// IsVisible returns whether the 'to' face is visible through
// 'from' face.
func (cb *ChunkBuffer) IsVisible(from, to direction.Type) bool {
	return (cb.cullBits & (1 << (from*6 + to))) != 0
}

// AllocateChunkBuffer allocates a chunk buffer and adds it to the
// render list.
func AllocateChunkBuffer(x, y, z int) *ChunkBuffer {
	c := &ChunkBuffer{
		position: position{X: x, Y: y, Z: z},
	}
	buffers[c.position] = c
	bufferColumns[positionC{x, z}]++
	return c
}

// Upload uploads the passed vertex data to the buffer.
func (cb *ChunkBuffer) Upload(data []byte, count int, cullBits uint64) {
	if cb.invalid {
		return
	}
	cb.cullBits = cullBits

	if count == 0 {
		if cb.array.IsValid() {
			cb.array.Delete()
			cb.buffer.Delete()
		}
		return
	}

	if !cb.array.IsValid() {
		cb.array = gl.CreateVertexArray()
		cb.buffer = gl.CreateBuffer()
	}

	cb.array.Bind()
	shaderChunk.Position.Enable()
	shaderChunk.TextureInfo.Enable()
	shaderChunk.TextureOffset.Enable()
	shaderChunk.Color.Enable()
	shaderChunk.Lighting.Enable()

	cb.buffer.Bind(gl.ArrayBuffer)
	cb.buffer.Data(data, gl.StaticDraw)
	shaderChunk.Position.PointerInt(3, gl.Short, 23, 0)
	shaderChunk.TextureInfo.Pointer(4, gl.UnsignedShort, false, 23, 6)
	shaderChunk.TextureOffset.Pointer(2, gl.Short, false, 23, 14)
	shaderChunk.Color.Pointer(3, gl.UnsignedByte, true, 23, 18)
	shaderChunk.Lighting.Pointer(2, gl.UnsignedByte, false, 23, 21)

	cb.count = count
}

// UploadTrans uploads the passed vertex data to the translucent buffer.
func (cb *ChunkBuffer) UploadTrans(data []byte, count int) {
	if cb.invalid {
		return
	}
	if count == 0 {
		if cb.arrayT.IsValid() {
			cb.arrayT.Delete()
			cb.bufferT.Delete()
		}
		return
	}

	if !cb.arrayT.IsValid() {
		cb.arrayT = gl.CreateVertexArray()
		cb.bufferT = gl.CreateBuffer()
	}
	cb.arrayT.Bind()
	shaderChunkT.Position.Enable()
	shaderChunkT.TextureInfo.Enable()
	shaderChunkT.TextureOffset.Enable()
	shaderChunkT.Color.Enable()
	shaderChunkT.Lighting.Enable()

	cb.bufferT.Bind(gl.ArrayBuffer)
	cb.bufferT.Data(data, gl.StaticDraw)
	shaderChunkT.Position.PointerInt(3, gl.Short, 23, 0)
	shaderChunkT.TextureInfo.Pointer(4, gl.UnsignedShort, false, 23, 6)
	shaderChunkT.TextureOffset.Pointer(2, gl.Short, false, 23, 14)
	shaderChunkT.Color.Pointer(3, gl.UnsignedByte, true, 23, 18)
	shaderChunkT.Lighting.Pointer(2, gl.UnsignedByte, false, 23, 21)

	cb.countT = count
}

// Free removes the buffer and frees related resources.
func (cb *ChunkBuffer) Free() {
	if cb.invalid {
		return
	}
	cb.invalid = true
	delete(buffers, cb.position)
	cpos := positionC{cb.position.X, cb.position.Z}
	val := bufferColumns[cpos]
	val--
	if val <= 0 {
		delete(bufferColumns, cpos)
	} else {
		bufferColumns[cpos] = val
	}

	cb.buffer.Delete()
	cb.array.Delete()
	cb.bufferT.Delete()
	cb.arrayT.Delete()
}
