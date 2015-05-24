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
	"strings"
	"unsafe"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/thinkofdeath/gl/v3.2-core/gl"
)

// ShaderType is type of shader to be used, different types run
// at different stages in the pipeline.
type ShaderType uint32

// Valid shader types.
const (
	VertexShader   ShaderType = gl.VERTEX_SHADER
	FragmentShader ShaderType = gl.FRAGMENT_SHADER
)

// ShaderParameter is a parameter that can set or read from a
// shader.
type ShaderParameter uint32

// Valid shader parameters.
const (
	CompileStatus ShaderParameter = gl.COMPILE_STATUS
	InfoLogLength ShaderParameter = gl.INFO_LOG_LENGTH
)

// Program is a collection of shaders which will be run on draw
// operations.
type Program uint32

// CreateProgram allocates a new program.
func CreateProgram() Program {
	return Program(gl.CreateProgram())
}

// AttachShader attaches the passed shader to the program.
func (p Program) AttachShader(s Shader) {
	gl.AttachShader(uint32(p), uint32(s))
}

// Link links the program's shaders.
func (p Program) Link() {
	gl.LinkProgram(uint32(p))
}

var (
	currentProgram Program
)

// Use sets this shader as the active shader. If this shader is
// already active this does nothing.
func (p Program) Use() {
	if p == currentProgram {
		return
	}
	gl.UseProgram(uint32(p))
	currentProgram = p
}

// Uniform is a per-a-draw value that can be passed into the
// program.
type Uniform int32

// UniformLocation returns the uniform with the given name in
// the program.
func (p Program) UniformLocation(name string) Uniform {
	n := gl.Str(name + "\x00")
	return Uniform(gl.GetUniformLocation(uint32(p), n))
}

// Matrix4 sets the value of the uniform to the passed matrix.
func (u Uniform) Matrix4(matrix *mgl32.Mat4) {
	gl.UniformMatrix4fv(int32(u), 1, false, (*float32)(unsafe.Pointer(matrix)))
}

// Matrix4 sets the value of the uniform to the passed matrix.
func (u Uniform) Matrix4Multi(matrix []mgl32.Mat4) {
	gl.UniformMatrix4fv(int32(u), int32(len(matrix)), false, (*float32)(gl.Ptr(matrix)))
}

// Int sets the value of the uniform to the passed integer.
func (u Uniform) Int(val int) {
	gl.Uniform1i(int32(u), int32(val))
}

// Int3 sets the value of the uniform to the passed integers.
func (u Uniform) Int3(x, y, z int) {
	gl.Uniform3i(int32(u), int32(x), int32(y), int32(z))
}

// IntV sets the value of the uniform to the passed integer array.
func (u Uniform) IntV(v ...int) {
	gl.Uniform1iv(int32(u), int32(len(v)), (*int32)(gl.Ptr(v)))
}

// Float sets the value of the uniform to the passed float.
func (u Uniform) Float(val float32) {
	gl.Uniform1f(int32(u), val)
}

// Float2 sets the value of the uniform to the passed floats.
func (u Uniform) Float2(x, y float32) {
	gl.Uniform2f(int32(u), x, y)
}

// Float3 sets the value of the uniform to the passed floats.
func (u Uniform) Float3(x, y, z float32) {
	gl.Uniform3f(int32(u), x, y, z)
}

// Float3 sets the value of the uniform to the passed floats.
func (u Uniform) Float4(x, y, z, w float32) {
	gl.Uniform4f(int32(u), x, y, z, w)
}

func (u Uniform) FloatMutli(a []float32) {
	gl.Uniform4fv(int32(u), int32(len(a)), (*float32)(gl.Ptr(a)))
}

func (u Uniform) FloatMutliRaw(a interface{}, l int) {
	gl.Uniform4fv(int32(u), int32(l), (*float32)(gl.Ptr(a)))
}

// Attribute is a per-a-vertex value that can be passed into the
// program.
type Attribute int32

// AttributeLocation returns the attribute with the given name in
// the program.
func (p Program) AttributeLocation(name string) Attribute {
	n := gl.Str(name + "\x00")
	return Attribute(gl.GetAttribLocation(uint32(p), n))
}

// Enable enables the attribute for use in rendering.
func (a Attribute) Enable() {
	gl.EnableVertexAttribArray(uint32(a))
}

// Disable disables the attribute for use in rendering.
func (a Attribute) Disable() {
	gl.DisableVertexAttribArray(uint32(a))
}

// Pointer is used to specify the format of the data in the buffer. The data will
// be uploaded as floats.
func (a Attribute) Pointer(size int, ty Type, normalized bool, stride, offset int) {
	gl.VertexAttribPointer(uint32(a), int32(size), uint32(ty), normalized, int32(stride), uintptr(offset))
}

// PointerInt is used to specify the format of the data in the buffer. The data will
// be uploaded as integers.
func (a Attribute) PointerInt(size int, ty Type, stride, offset int) {
	gl.VertexAttribIPointer(uint32(a), int32(size), uint32(ty), int32(stride), uintptr(offset))
}

// Shader is code to be run on the gpu at a specific stage in the
// pipeline.
type Shader uint32

// CreateShader creates a new shader of the specifed type.
func CreateShader(t ShaderType) Shader {
	return Shader(gl.CreateShader(uint32(t)))
}

// Source sets the source of the shader.
func (s Shader) Source(src string) {
	ss := gl.Str(src + "\x00")
	gl.ShaderSource(uint32(s), 1, &ss, nil)
}

// Compile compiles the shader.
func (s Shader) Compile() {
	gl.CompileShader(uint32(s))
}

// Parameter returns the integer value of the parameter for
// this shader.
func (s Shader) Parameter(param ShaderParameter) int {
	var p int32
	gl.GetShaderiv(uint32(s), uint32(param), &p)
	return int(p)
}

// InfoLog returns the log from compiling the shader.
func (s Shader) InfoLog() string {
	l := s.Parameter(InfoLogLength)

	var ptr unsafe.Pointer
	var buf []byte
	if l > 0 {
		buf = make([]byte, l)
		ptr = gl.Ptr(buf)
	}

	gl.GetShaderInfoLog(uint32(s), int32(l), nil, (*uint8)(ptr))
	return strings.TrimRight(string(buf), "\x00")
}
