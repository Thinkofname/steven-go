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

	"github.com/thinkofdeath/steven/native"
	"github.com/thinkofdeath/steven/render/gl"
	"github.com/thinkofdeath/steven/type/direction"
)

var (
	buffers           = make(map[position]*ChunkBuffer)
	elementBuffer     gl.Buffer
	elementBufferSize int
	elementBufferType gl.Type = gl.UnsignedShort
)

// ChunkBuffer is a renderable chunk section
type ChunkBuffer struct {
	position
	invalid bool

	array        gl.VertexArray
	buffer       gl.Buffer
	bufferSize   int
	count        int
	arrayT       gl.VertexArray
	bufferT      gl.Buffer
	bufferTI     gl.Buffer
	bufferTIType gl.Type
	bufferTSize  int
	countT       int
	cullBits     uint64

	renderedOn uint

	transInfo objectInfoList
	transData []byte

	neighborChunks [6]*ChunkBuffer
}

// IsVisible returns whether the 'to' face is visible through
// 'from' face.
func (cb *ChunkBuffer) IsVisible(from, to direction.Type) bool {
	return (cb.cullBits & (1 << (from*6 + to))) != 0
}

// AllocateColumn ensures the column's buffers are allocated.
func AllocateColumn(x, z int) {
	for i := 0; i < 16; i++ {
		if _, ok := buffers[position{x, i, z}]; !ok {
			buffers[position{x, i, z}] = &ChunkBuffer{
				position: position{X: x, Y: i, Z: z},
				cullBits: math.MaxUint64,
				invalid:  true,
			}
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

// FreeColumn deallocates the column's buffers.
func FreeColumn(x, z int) {
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
}

// AllocateChunkBuffer allocates a chunk buffer and adds it to the
// render list.
func AllocateChunkBuffer(x, y, z int) *ChunkBuffer {
	c := buffers[position{x, y, z}]
	c.invalid = false
	return c
}

func ensureElementBuffer(size int) {
	if elementBufferSize < size {
		data, ty := genElementBuffer(size)
		elementBufferType = ty
		elementBuffer.Bind(gl.ElementArrayBuffer)
		elementBuffer.Data(data, gl.DynamicDraw)
		elementBufferSize = size
	}
}

func genElementBuffer(size int) ([]byte, gl.Type) {
	data := make([]byte, size*4)
	offset := 0
	ty := gl.UnsignedShort
	if uint32(size/6)*4+3 >= math.MaxUint16 {
		ty = gl.UnsignedInt
	}
	for i := 0; i < size/6; i++ {
		for _, val := range []uint32{0, 1, 2, 3, 2, 1} {
			if ty == gl.UnsignedInt {
				native.Order.PutUint32(data[offset:], uint32(i)*4+val)
				offset += 4
			} else {
				native.Order.PutUint16(data[offset:], uint16(uint32(i)*4+val))
				offset += 2
			}
		}
	}
	return data, ty
}

// Upload uploads the passed vertex data to the buffer.
func (cb *ChunkBuffer) Upload(data []byte, indices int, cullBits uint64) {
	if cb.invalid {
		return
	}
	cb.cullBits = cullBits
	var n bool

	if indices == 0 {
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

	ensureElementBuffer(indices)
	elementBuffer.Bind(gl.ElementArrayBuffer)

	cb.buffer.Bind(gl.ArrayBuffer)
	if n || len(data) > cb.bufferSize {
		cb.bufferSize = len(data)
		cb.buffer.Data(data, gl.DynamicDraw)
	} else {
		target := cb.buffer.Map(gl.WriteOnly, len(data))
		copy(target, data)
		cb.buffer.Unmap()
	}
	shaderChunk.Position.PointerInt(3, gl.Short, 28, 0)
	shaderChunk.TextureInfo.Pointer(4, gl.UnsignedShort, false, 28, 6)
	shaderChunk.TextureOffset.Pointer(3, gl.Short, false, 28, 14)
	shaderChunk.Color.Pointer(3, gl.UnsignedByte, true, 28, 20)
	shaderChunk.Lighting.Pointer(2, gl.UnsignedShort, false, 28, 24)

	cb.count = indices
}

// UploadTrans uploads the passed vertex data to the translucent buffer.
func (cb *ChunkBuffer) UploadTrans(info []ObjectInfo, data []byte, indices int) {
	if cb.invalid {
		return
	}
	var n bool
	if indices == 0 {
		if cb.arrayT.IsValid() {
			cb.arrayT.Delete()
			cb.bufferT.Delete()
			cb.bufferTI.Delete()
		}
		cb.transInfo = nil
		return
	}

	if !cb.arrayT.IsValid() {
		cb.arrayT = gl.CreateVertexArray()
		cb.bufferT = gl.CreateBuffer()
		cb.bufferTI = gl.CreateBuffer()
		n = true
	}

	cb.arrayT.Bind()
	shaderChunkT.Position.Enable()
	shaderChunkT.TextureInfo.Enable()
	shaderChunkT.TextureOffset.Enable()
	shaderChunkT.Color.Enable()
	shaderChunkT.Lighting.Enable()

	cb.bufferTI.Bind(gl.ElementArrayBuffer)
	cb.transData, cb.bufferTIType = genElementBuffer(indices)
	cb.bufferTI.Data(cb.transData, gl.StreamDraw)

	cb.bufferT.Bind(gl.ArrayBuffer)
	if n || len(data) > cb.bufferTSize {
		cb.bufferTSize = len(data)
		cb.bufferT.Data(data, gl.DynamicDraw)
	} else {
		target := cb.bufferT.Map(gl.WriteOnly, len(data))
		copy(target, data)
		cb.bufferT.Unmap()
	}
	shaderChunkT.Position.PointerInt(3, gl.Short, 28, 0)
	shaderChunkT.TextureInfo.Pointer(4, gl.UnsignedShort, false, 28, 6)
	shaderChunkT.TextureOffset.Pointer(3, gl.Short, false, 28, 14)
	shaderChunkT.Color.Pointer(3, gl.UnsignedByte, true, 28, 20)
	shaderChunkT.Lighting.Pointer(2, gl.UnsignedShort, false, 28, 24)

	cb.countT = indices
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
	cb.transInfo = nil
	cb.transData = nil
	cb.cullBits = math.MaxUint64

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
	if cb.bufferTI.IsValid() {
		cb.bufferTI.Delete()
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
