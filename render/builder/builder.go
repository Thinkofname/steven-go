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

// Package builder provides a simple way to create buffers to upload to the gpu
package builder

import (
	"bytes"
	"math"

	"github.com/thinkofdeath/steven/platform"
)

// Types allowed to be used in a buffer.
const (
	UnsignedByte  Type = 1
	Byte          Type = 1
	UnsignedShort Type = 2
	Short         Type = 2
	Float         Type = 4
)

// Type is a type that is allowed in a buffer.
type Type int

// Buffer is a dynamically sized byte buffer
// for creating data to upload to the gpu.
type Buffer struct {
	buf      bytes.Buffer
	elemSize int

	scratch [8]byte
}

// New creates a new buffer containing the passed
// types.
func New(types ...Type) *Buffer {
	elemSize := 0
	for _, t := range types {
		elemSize += int(t)
	}
	b := &Buffer{
		elemSize: elemSize,
	}
	b.buf.Grow(elemSize * 3000)
	return b
}

// UnsignedByte appends an unsigned byte to the
// buffer.
func (b *Buffer) UnsignedByte(i byte) {
	b.buf.WriteByte(i)
}

// Byte appends a signed byte to the buffer.
func (b *Buffer) Byte(i int8) {
	b.UnsignedByte(byte(i))
}

// UnsignedShort writes an unsigned short to the
// buffer
func (b *Buffer) UnsignedShort(i uint16) {
	d := b.scratch[:2]
	platform.NativeOrder.PutUint16(d, i)
	b.buf.Write(d)
}

// Short writes a short to the buffer.
func (b *Buffer) Short(i int16) {
	b.UnsignedShort(uint16(i))
}

// Float writes a float to the buffer
func (b *Buffer) Float(f float32) {
	d := b.scratch[:4]
	i := math.Float32bits(f)
	platform.NativeOrder.PutUint32(d, i)
	b.buf.Write(d)
}

// WriteBuffer copies the passed buffer to this buffer
func (b *Buffer) WriteBuffer(o *Buffer) {
	o.buf.WriteTo(&b.buf)
}

// Count returns the number of vertices in the buffer
func (b *Buffer) Count() int {
	return b.buf.Len() / b.elemSize
}

// Data returns a byte slice of the buffer
func (b *Buffer) Data() []byte {
	return b.buf.Bytes()
}

// ElementSize returns the size of a single vertex in the
// buffer
func (b *Buffer) ElementSize() int {
	return b.elemSize
}

// Reset resets the internal buffer ready for reuse
func (b *Buffer) Reset() {
	b.buf.Reset()
}
