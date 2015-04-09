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
	"math"

	"github.com/thinkofdeath/steven/render/gl"
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

	array       gl.VertexArray
	buffer      gl.Buffer
	bufferSize  int
	count       int
	arrayT      gl.VertexArray
	bufferT     gl.Buffer
	bufferTSize int
	countT      int
	cullBits    uint64

	renderedOn uint

	transData []byte
	transInfo objectInfoList

	neighborChunks [6]*ChunkBuffer
}

// IsVisible returns whether the 'to' face is visible through
// 'from' face.
func (cb *ChunkBuffer) IsVisible(from, to direction.Type) bool {
	return (cb.cullBits & (1 << (from*6 + to))) != 0
}

// AllocateChunkBuffer allocates a chunk buffer and adds it to the
// render list.
func AllocateChunkBuffer(x, y, z int) *ChunkBuffer {
	if _, ok := bufferColumns[positionC{x, z}]; !ok {
		for i := 0; i < 16; i++ {
			buffers[position{x, i, z}] = &ChunkBuffer{
				position: position{X: x, Y: i, Z: z},
				cullBits: math.MaxUint64,
				invalid:  true,
			}
		}
		// Update neighbors
		for i := 0; i < 16; i++ {
			c := buffers[position{x, i, z}]
			for _, d := range direction.Values {
				ox, oy, oz := d.Offset()
				o := buffers[position{x + ox, i + oy, z + oz}]
				if o != nil {
					c.neighborChunks[d] = o
					o.neighborChunks[d.Opposite()] = c
				}
			}
		}
	}
	bufferColumns[positionC{x, z}]++
	c := buffers[position{x, y, z}]
	c.invalid = false
	return c
}

// Upload uploads the passed vertex data to the buffer.
func (cb *ChunkBuffer) Upload(data []byte, count int, cullBits uint64) {
	if cb.invalid {
		return
	}
	cb.cullBits = cullBits
	var n bool

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
		n = true
	}

	cb.array.Bind()
	shaderChunk.Position.Enable()
	shaderChunk.TextureInfo.Enable()
	shaderChunk.TextureOffset.Enable()
	shaderChunk.Color.Enable()
	shaderChunk.Lighting.Enable()

	cb.buffer.Bind(gl.ArrayBuffer)
	if n || len(data) > cb.bufferSize {
		cb.bufferSize = len(data)
		cb.buffer.Data(data, gl.DynamicDraw)
	} else {
		target := cb.buffer.Map(gl.WriteOnly, len(data))
		copy(target, data)
		cb.buffer.Unmap()
	}
	shaderChunk.Position.PointerInt(3, gl.Short, 23, 0)
	shaderChunk.TextureInfo.Pointer(4, gl.UnsignedShort, false, 23, 6)
	shaderChunk.TextureOffset.Pointer(2, gl.Short, false, 23, 14)
	shaderChunk.Color.Pointer(3, gl.UnsignedByte, true, 23, 18)
	shaderChunk.Lighting.Pointer(2, gl.UnsignedByte, false, 23, 21)

	cb.count = count
}

// UploadTrans uploads the passed vertex data to the translucent buffer.
func (cb *ChunkBuffer) UploadTrans(info []ObjectInfo, data []byte, count int) {
	if cb.invalid {
		return
	}
	var n bool
	if count == 0 {
		if cb.arrayT.IsValid() {
			cb.arrayT.Delete()
			cb.bufferT.Delete()
		}
		cb.transData = nil
		cb.transInfo = nil
		return
	}

	if !cb.arrayT.IsValid() {
		cb.arrayT = gl.CreateVertexArray()
		cb.bufferT = gl.CreateBuffer()
		n = true
	}
	cb.transData = make([]byte, len(data))
	copy(cb.transData, data)

	cb.arrayT.Bind()
	shaderChunkT.Position.Enable()
	shaderChunkT.TextureInfo.Enable()
	shaderChunkT.TextureOffset.Enable()
	shaderChunkT.Color.Enable()
	shaderChunkT.Lighting.Enable()

	cb.bufferT.Bind(gl.ArrayBuffer)
	if n || len(cb.transData) > cb.bufferTSize {
		cb.bufferTSize = len(cb.transData)
		cb.bufferT.Data(cb.transData, gl.StreamDraw)
	}
	shaderChunkT.Position.PointerInt(3, gl.Short, 23, 0)
	shaderChunkT.TextureInfo.Pointer(4, gl.UnsignedShort, false, 23, 6)
	shaderChunkT.TextureOffset.Pointer(2, gl.Short, false, 23, 14)
	shaderChunkT.Color.Pointer(3, gl.UnsignedByte, true, 23, 18)
	shaderChunkT.Lighting.Pointer(2, gl.UnsignedByte, false, 23, 21)

	cb.countT = count
	cb.transInfo = info
}

// Free removes the buffer and frees related resources.
func (cb *ChunkBuffer) Free() {
	if cb.invalid {
		return
	}
	// Clear state
	cb.invalid = true
	cb.count = 0
	cb.countT = 0
	cb.transData = nil
	cb.transInfo = nil
	cpos := positionC{cb.position.X, cb.position.Z}
	val := bufferColumns[cpos]
	val--
	if val <= 0 {
		delete(bufferColumns, cpos)
		x, z := cb.X, cb.Z
		for i := 0; i < 16; i++ {
			// Update neighbors
			c := buffers[position{x, i, z}]
			for _, d := range direction.Values {
				ox, oy, oz := d.Offset()
				o := buffers[position{x + ox, i + oy, z + oz}]
				if o != nil {
					c.neighborChunks[d] = nil
					o.neighborChunks[d.Opposite()] = nil
				}
			}
			delete(buffers, position{x, i, z})
		}
	} else {
		bufferColumns[cpos] = val
	}

	if cb.buffer.IsValid() {
		cb.buffer.Delete()
	}
	if cb.array.IsValid() {
		cb.array.Delete()
	}
	if cb.bufferT.IsValid() {
		cb.bufferT.Delete()
	}
	if cb.arrayT.IsValid() {
		cb.arrayT.Delete()
	}
}

// ObjectInfo contains information about an renderable object that needs
// to be sorted before rendering.
type ObjectInfo struct {
	X, Y, Z       int
	Offset, Count int
}

type objectInfoList []ObjectInfo

func (o objectInfoList) Swap(a, b int) {
	o[a], o[b] = o[b], o[a]
}

func (o objectInfoList) Less(aa, bb int) bool {
	a := o[aa]
	b := o[bb]
	dx := float64(a.X) + 0.5 - Camera.X
	dy := float64(a.Y) + 0.5 - Camera.Y
	dz := float64(a.Z) + 0.5 - Camera.Z
	adist := dx*dx + dy*dy + dz*dz

	dx = float64(b.X) + 0.5 - Camera.X
	dy = float64(b.Y) + 0.5 - Camera.Y
	dz = float64(b.Z) + 0.5 - Camera.Z
	bdist := dx*dx + dy*dy + dz*dz
	return adist > bdist
}

func (o objectInfoList) Len() int {
	return len(o)
}
