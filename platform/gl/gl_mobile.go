// +build mobile

// Package gl provides a more Go friendly OpenGL API
package gl

import (
	"fmt"
	"golang.org/x/mobile/gl"
)

const (
	ColorBufferBit   ClearFlags = gl.COLOR_BUFFER_BIT
	DepthBufferBit   ClearFlags = gl.DEPTH_BUFFER_BIT
	StencilBufferBit ClearFlags = gl.STENCIL_BUFFER_BIT

	DepthTest    Flag = gl.DEPTH_TEST
	CullFaceFlag Flag = gl.CULL_FACE
	StencilTest  Flag = gl.STENCIL_TEST

	Back  Face = gl.BACK
	Front Face = gl.FRONT

	ClockWise        FaceDirection = gl.CW
	CounterClockWise FaceDirection = gl.CCW

	Triangles DrawType = gl.TRIANGLES

	Never       Func = gl.NEVER
	Less        Func = gl.LESS
	LessOrEqual Func = gl.LEQUAL
	Greater     Func = gl.GREATER
	Always      Func = gl.ALWAYS
	Equal       Func = gl.EQUAL

	Replace Op = gl.REPLACE
	Keep    Op = gl.KEEP
	Zero    Op = gl.ZERO
)

type (
	ClearFlags    uint32
	Flag          uint32
	Face          uint32
	FaceDirection uint32
	DrawType      uint32
	Func          uint32
	Op            uint32
)

func Viewport(x, y, width, height int) {
	gl.Viewport(x, y, width, height)
}

func ClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func Clear(flags ClearFlags) {
	gl.Clear(gl.Enum(flags))
}

func ActiveTexture(id int) {
	gl.ActiveTexture(gl.TEXTURE0 + gl.Enum(id))
}

func Enable(flag Flag) {
	gl.Enable(gl.Enum(flag))
}

func Disable(flag Flag) {
	gl.Disable(gl.Enum(flag))
}

func CullFace(face Face) {
	gl.CullFace(gl.Enum(face))
}

func FrontFace(dir FaceDirection) {
	gl.FrontFace(gl.Enum(dir))
}

func DrawArrays(ty DrawType, offset, count int) {
	gl.DrawArrays(gl.Enum(ty), offset, count)
}

func checkError() {
	err := gl.GetError()
	if err != 0 {
		panic(fmt.Sprintf("gl error: %d", err))
	}
}

func Flush() {
	gl.Flush()
}

func DepthMask(f bool) {
	gl.DepthMask(f)
}

func ColorMask(r, g, b, a bool) {
	gl.ColorMask(r, g, b, a)
}

func StencilFunc(f Func, ref, mask int) {
	gl.StencilFunc(gl.Enum(f), ref, uint32(mask))
}

func StencilMask(mask int) {
	gl.StencilMask(uint32(mask))
}

func StencilOp(op, fail, pass Op) {
	gl.StencilOp(gl.Enum(op), gl.Enum(fail), gl.Enum(pass))
}

func ClearStencil(i int) {
	gl.ClearStencil(i)
}
