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
	"log"
	"unsafe"

	"github.com/thinkofdeath/gl/v3.2-core/gl"
)

// ClearFlags is a set of flags to mark what should be cleared during
// a Clear call.
type ClearFlags uint32

const (
	// ColorBufferBit marks the color buffer to be cleared
	ColorBufferBit ClearFlags = gl.COLOR_BUFFER_BIT
	// DepthBufferBit marks the depth buffer to be cleared
	DepthBufferBit ClearFlags = gl.DEPTH_BUFFER_BIT
	// StencilBufferBit marks the stencil buffer to be cleared
	StencilBufferBit ClearFlags = gl.STENCIL_BUFFER_BIT
)

// Flag is a setting that can be enabled or disabled on the context.
type Flag uint32

// Valid flags
const (
	DepthTest    Flag = gl.DEPTH_TEST
	CullFaceFlag Flag = gl.CULL_FACE
	StencilTest  Flag = gl.STENCIL_TEST
	Blend        Flag = gl.BLEND
	DebugOutput  Flag = gl.DEBUG_OUTPUT
	Multisample  Flag = gl.MULTISAMPLE
)

// Face specifies a face to act on.
type Face uint32

// Valid faces
const (
	Back  Face = gl.BACK
	Front Face = gl.FRONT
)

// FaceDirection is used to specify an order of vertices, normally
// used to set which is considered to be the front face.
type FaceDirection uint32

// Valid directions for vertices for faces.
const (
	ClockWise        FaceDirection = gl.CW
	CounterClockWise FaceDirection = gl.CCW
)

// DrawType is used to specify how the vertices will be handled
// to draw.
type DrawType uint32

const (
	// Triangles treats each set of 3 vertices as a triangle
	Triangles DrawType = gl.TRIANGLES
	// LineStrip means the previous vertex connects to the next
	// one in a continuous strip.
	LineStrip DrawType = gl.LINE_STRIP
	// Lines treats each set of 2 vertices as a line
	Lines DrawType = gl.LINES
)

// Func is a function to be preformed on two values.
type Func uint32

// Functions
const (
	Never       Func = gl.NEVER
	Less        Func = gl.LESS
	LessOrEqual Func = gl.LEQUAL
	Greater     Func = gl.GREATER
	Always      Func = gl.ALWAYS
	Equal       Func = gl.EQUAL
)

// Op is an operation to be applied (depending on the method used)
type Op uint32

// Valid operations
const (
	Replace Op = gl.REPLACE
	Keep    Op = gl.KEEP
	Zero    Op = gl.ZERO
)

// Factor is used in blending
type Factor uint32

// Valid factors
const (
	SrcAlpha         Factor = gl.SRC_ALPHA
	OneMinusSrcAlpha Factor = gl.ONE_MINUS_SRC_ALPHA
)

// Init inits the gl library. This should be called once a context is ready.
func Init() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
}

// Viewport sets the size of the viewport of this context.
func Viewport(x, y, width, height int) {
	gl.Viewport(int32(x), int32(y), int32(width), int32(height))
}

// ClearColor sets the color the color buffer should be cleared to
// when Clear is called with the color flag.
func ClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

// Clear clears the buffers specified by the passed flags.
func Clear(flags ClearFlags) {
	gl.Clear(uint32(flags))
}

// ActiveTexture sets the texture slot with the passed id as the
// currently active one.
func ActiveTexture(id int) {
	gl.ActiveTexture(gl.TEXTURE0 + uint32(id))
}

// Enable enables the passed flag.
func Enable(flag Flag) {
	gl.Enable(uint32(flag))
}

// Disable disables the passed flag.
func Disable(flag Flag) {
	gl.Disable(uint32(flag))
}

// CullFace sets the face to be culled by the gpu.
func CullFace(face Face) {
	gl.CullFace(uint32(face))
}

// FrontFace sets the direction of vertices used to specify the
// front face (e.g. for culling).
func FrontFace(dir FaceDirection) {
	gl.FrontFace(uint32(dir))
}

// DrawArrays draws the passed number of triangles starting at the
// passed offset using data from the currently bound buffer(s).
// The DrawType specifies how the shapes (triangles, lines etc)
// will be formed from the data.
func DrawArrays(ty DrawType, offset, count int) {
	gl.DrawArrays(uint32(ty), int32(offset), int32(count))
}

func DrawElements(ty DrawType, count int, dty Type, offset int) {
	gl.DrawElements(uint32(ty), int32(count), uint32(dty), uintptr(offset))
}

func MultiDrawElements(ty DrawType, count []int32, dty Type, offset []uintptr) {
	gl.MultiDrawElements(uint32(ty), (*int32)(gl.Ptr(count)), uint32(dty), (*uintptr)(gl.Ptr(offset)), int32(len(count)))
}

// CheckError panics if there has been an error reported to the
// context. This is normally not a cheap call so shouldn't be
// used in production.
func CheckError() {
	err := gl.GetError()
	if err != 0 {
		panic(fmt.Sprintf("gl error: %d", err))
	}
}

// Flush flushes all commands in the queue and waits for completion.
func Flush() {
	gl.Flush()
}

// DepthMask sets whether the depth buffer can be written too.
func DepthMask(f bool) {
	gl.DepthMask(f)
}

// DepthFunc sets the function used to decide whether to cull a pixel
// when depth testing
func DepthFunc(f Func) {
	gl.DepthFunc(uint32(f))
}

// ColorMask sets whether each color channel be the written too.
func ColorMask(r, g, b, a bool) {
	gl.ColorMask(r, g, b, a)
}

// StencilFunc sets the function to be used when comparing with the
// stencil buffer.
func StencilFunc(f Func, ref, mask int) {
	gl.StencilFunc(uint32(f), int32(ref), uint32(mask))
}

// StencilMask sets the value to be written to the stencil buffer on
// success.
func StencilMask(mask int) {
	gl.StencilMask(uint32(mask))
}

// StencilOp sets the operation to be executed depending on the result
// of the stencil test.
func StencilOp(op, fail, pass Op) {
	gl.StencilOp(uint32(op), uint32(fail), uint32(pass))
}

// ClearStencil clears the stencil buffer will the passed value.
func ClearStencil(i int) {
	gl.ClearStencil(int32(i))
}

// BlendFunc sets the factors to be used when blending.
func BlendFunc(sFactor, dFactor Factor) {
	gl.BlendFunc(uint32(sFactor), uint32(dFactor))
}

// DebugLog enables OpenGL's debug messages and logs them to stdout.
func DebugLog() {
	gl.DebugMessageCallback(func(
		source uint32,
		gltype uint32,
		id uint32,
		severity uint32,
		length int32,
		message string,
		userParam unsafe.Pointer) {
		// Source
		strSource := "unknown"
		switch source {
		case gl.DEBUG_SOURCE_API:
			strSource = "api"
		case gl.DEBUG_SOURCE_WINDOW_SYSTEM:
			strSource = "windowSystem"
		case gl.DEBUG_SOURCE_SHADER_COMPILER:
			strSource = "shaderCompiler"
		case gl.DEBUG_SOURCE_THIRD_PARTY:
			strSource = "thirdParty"
		case gl.DEBUG_SOURCE_APPLICATION:
			strSource = "application"
		case gl.DEBUG_SOURCE_OTHER:
			strSource = "other"
		}
		// Type
		strType := "unknown"
		switch gltype {
		case gl.DEBUG_TYPE_ERROR:
			strType = "error"
		case gl.DEBUG_TYPE_DEPRECATED_BEHAVIOR:
			strType = "deprecatedBehavior"
		case gl.DEBUG_TYPE_UNDEFINED_BEHAVIOR:
			strType = "undefinedBehavior"
		case gl.DEBUG_TYPE_PORTABILITY:
			strType = "portability"
		case gl.DEBUG_TYPE_PERFORMANCE:
			strType = "performance"
		case gl.DEBUG_TYPE_MARKER:
			strType = "marker"
		case gl.DEBUG_TYPE_PUSH_GROUP:
			strType = "pushGroup"
		case gl.DEBUG_TYPE_POP_GROUP:
			strType = "popGroup"
		case gl.DEBUG_TYPE_OTHER:
			strType = "other"
		}
		// Severity
		strSeverity := "unknown"
		switch severity {
		case gl.DEBUG_SEVERITY_HIGH:
			strSeverity = "high"
		case gl.DEBUG_SEVERITY_MEDIUM:
			strSeverity = "medium"
		case gl.DEBUG_SEVERITY_LOW:
			strSeverity = "low"
		case gl.DEBUG_SEVERITY_NOTIFICATION:
			return
		}
		log.Printf("[%s][%s][%s]: %s\n", strSource, strType, strSeverity, message)
	}, nil)
}
