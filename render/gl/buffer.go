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
	"unsafe"

	"github.com/thinkofdeath/gl/v3.2-core/gl"
)

// BufferTarget is a target for a buffer to be bound to.
type BufferTarget uint32

const (
	// ArrayBuffer is a buffer target intended for vertex data.
	ArrayBuffer        BufferTarget = gl.ARRAY_BUFFER
	ElementArrayBuffer BufferTarget = gl.ELEMENT_ARRAY_BUFFER
)

// BufferUsage states how a buffer is going to be used by the program.
type BufferUsage uint32

const (
	// StaticDraw marks the buffer as 'not going to change' after the
	// initial data upload to be rendered by the gpu.
	StaticDraw BufferUsage = gl.STATIC_DRAW
	// DynamicDraw marks the buffer as 'changed frequently' during the
	// course of the program whilst being rendered by the gpu.
	DynamicDraw BufferUsage = gl.DYNAMIC_DRAW
	// StreamDraw marks the buffer as 'changed every frame' whilst being
	// rendered by the gpu.
	StreamDraw BufferUsage = gl.STREAM_DRAW
)

// Access states how a value will be accesed by the program.
type Access uint32

const (
	// ReadOnly states that the returned value will only be read.
	ReadOnly Access = gl.READ_ONLY
	// WriteOnly states that the returned value will only be written
	// to.
	WriteOnly Access = gl.WRITE_ONLY
)

// Buffer is a storage for vertex data.
type Buffer struct {
	internal uint32
}

// CreateBuffer allocates a new Buffer. If the allocation fails IsValid
// will return false.
func CreateBuffer() Buffer {
	var buffer Buffer
	gl.GenBuffers(1, &buffer.internal)
	return buffer
}

var (
	currentBufferTarget BufferTarget
)

// Bind makes the buffer the currently active one for the given target.
// This will allow it to be the source of operations that act on a buffer
// (Data, Map etc). If the buffer is already bound then this does nothing.
func (b Buffer) Bind(target BufferTarget) {
	gl.BindBuffer(uint32(target), b.internal)
	currentBufferTarget = target
}

// Data uploads the passed data to the gpu to be placed in this buffer.
// The usage specifies how the program plans to use this buffer.
func (b Buffer) Data(data []byte, usage BufferUsage) {
	var ptr unsafe.Pointer
	if len(data) != 0 {
		ptr = gl.Ptr(data)
	}
	gl.BufferData(uint32(currentBufferTarget), len(data), ptr, uint32(usage))
}

// Map maps the memory in the buffer on the gpu to memory which the program
// can access. The access flag will specify how the program plans to use the
// returned data. Unmapped must be called to return the data to the gpu.
//
// Warning: the passed length value is not checked in anyway so it is
// possible to overrun the memory. It is up to the program to ensure this
// length is valid.
func (b Buffer) Map(access Access, length int) []byte {
	ptr := gl.MapBuffer(uint32(currentBufferTarget), uint32(access))
	return (*[1 << 30]byte)(ptr)[:length:length]
}

// Unmap returns mapped memory for the buffer, as returned by Map, to the gpu
// so that it can be used in draw operations.
func (b Buffer) Unmap() {
	gl.UnmapBuffer(uint32(currentBufferTarget))
}

// Delete deallocates the buffer and any stored data. IsValid will
// return false after this call.
func (b *Buffer) Delete() {
	gl.DeleteBuffers(1, &b.internal)
	b.internal = 0
}

// IsValid returns whether this Buffer is still valid. A
// Buffer will become invalid after Delete is called.
func (b Buffer) IsValid() bool {
	return b.internal != 0
}
