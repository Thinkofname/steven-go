// +build mobile

package gl

import (
	"golang.org/x/mobile/gl"
)

const (
	ArrayBuffer BufferTarget = gl.ARRAY_BUFFER

	StaticDraw  BufferUsage = gl.STATIC_DRAW
	DynamicDraw BufferUsage = gl.DYNAMIC_DRAW
)

type Buffer gl.Buffer

type BufferTarget uint32
type BufferUsage uint32

func CreateBuffer() Buffer {
	return Buffer(gl.GenBuffer())
}

var (
	currentBuffer       Buffer
	currentBufferTarget BufferTarget
)

func (b Buffer) Bind(target BufferTarget) {
	if currentBuffer == b && currentBufferTarget == target {
		return
	}
	gl.BindBuffer(gl.Enum(target), gl.Buffer(b))
	currentBuffer = b
	currentBufferTarget = target
}

func (b Buffer) Data(data []byte, usage BufferUsage) {
	if currentBuffer != b {
		panic("buffer not bound")
	}
	if len(data) == 0 {
		return
	}
	gl.BufferData(gl.Enum(currentBufferTarget), gl.Enum(usage), data)
}

func (b Buffer) Delete() {
	gl.DeleteBuffer(gl.Buffer(b))
	if currentBuffer == b {
		currentBuffer = Buffer{}
	}
}
