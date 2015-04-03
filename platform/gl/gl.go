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

// Package gl provides a more Go friendly OpenGL API
package gl

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v3.2-core/gl"
)

const (
	ColorBufferBit   ClearFlags = gl.COLOR_BUFFER_BIT
	DepthBufferBit   ClearFlags = gl.DEPTH_BUFFER_BIT
	StencilBufferBit ClearFlags = gl.STENCIL_BUFFER_BIT

	DepthTest    Flag = gl.DEPTH_TEST
	CullFaceFlag Flag = gl.CULL_FACE
	StencilTest  Flag = gl.STENCIL_TEST
	Blend        Flag = gl.BLEND
	DebugOutput  Flag = gl.DEBUG_OUTPUT

	Back  Face = gl.BACK
	Front Face = gl.FRONT

	ClockWise        FaceDirection = gl.CW
	CounterClockWise FaceDirection = gl.CCW

	Triangles DrawType = gl.TRIANGLES
	LineStrip DrawType = gl.LINE_STRIP

	Never       Func = gl.NEVER
	Less        Func = gl.LESS
	LessOrEqual Func = gl.LEQUAL
	Greater     Func = gl.GREATER
	Always      Func = gl.ALWAYS
	Equal       Func = gl.EQUAL

	Replace Op = gl.REPLACE
	Keep    Op = gl.KEEP
	Zero    Op = gl.ZERO

	SrcAlpha         Factor = gl.SRC_ALPHA
	OneMinusSrcAlpha Factor = gl.ONE_MINUS_SRC_ALPHA
)

type (
	ClearFlags    uint32
	Flag          uint32
	Face          uint32
	FaceDirection uint32
	DrawType      uint32
	Func          uint32
	Op            uint32
	Factor        uint32
)

func Viewport(x, y, width, height int) {
	gl.Viewport(int32(x), int32(y), int32(width), int32(height))
}

func ClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func Clear(flags ClearFlags) {
	gl.Clear(uint32(flags))
}

func ActiveTexture(id int) {
	gl.ActiveTexture(gl.TEXTURE0 + uint32(id))
}

func Enable(flag Flag) {
	gl.Enable(uint32(flag))
}

func Disable(flag Flag) {
	gl.Disable(uint32(flag))
}

func CullFace(face Face) {
	gl.CullFace(uint32(face))
}

func FrontFace(dir FaceDirection) {
	gl.FrontFace(uint32(dir))
}

func DrawArrays(ty DrawType, offset, count int) {
	gl.DrawArrays(uint32(ty), int32(offset), int32(count))
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
	gl.StencilFunc(uint32(f), int32(ref), uint32(mask))
}

func StencilMask(mask int) {
	gl.StencilMask(uint32(mask))
}

func StencilOp(op, fail, pass Op) {
	gl.StencilOp(uint32(op), uint32(fail), uint32(pass))
}

func ClearStencil(i int) {
	gl.ClearStencil(int32(i))
}

func BlendFunc(sFactor, dFactor Factor) {
	gl.BlendFunc(uint32(sFactor), uint32(dFactor))
}

func DebugLog() {
	gl.DebugMessageCallback(func(
		source uint32,
		gltype uint32,
		id uint32,
		severity uint32,
		length int32,
		message string,
		userParam unsafe.Pointer) {
		fmt.Printf("GL: %d, %d, %d, %d: %s\n", source, gltype, id, severity, message)
	}, nil)
}
