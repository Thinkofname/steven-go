// +build !mobile

package gl

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/thinkofdeath/steven/vmath"
	"unsafe"
)

const (
	VertexShader   ShaderType = gl.VERTEX_SHADER
	FragmentShader ShaderType = gl.FRAGMENT_SHADER

	CompileStatus ShaderParameter = gl.COMPILE_STATUS
	InfoLogLength ShaderParameter = gl.INFO_LOG_LENGTH
)

type (
	Program   uint32
	Attribute int32
	Uniform   int32
)

func CreateProgram() Program {
	return Program(gl.CreateProgram())
}

func (p Program) AttachShader(s Shader) {
	gl.AttachShader(uint32(p), uint32(s))
}

func (p Program) Link() {
	gl.LinkProgram(uint32(p))
}

var (
	currentProgram Program
)

func (p Program) Use() {
	if p == currentProgram {
		return
	}
	gl.UseProgram(uint32(p))
	currentProgram = p
}

func (p Program) AttributeLocation(name string) Attribute {
	n := gl.Str(name + "\x00")
	return Attribute(gl.GetAttribLocation(uint32(p), n))
}

func (p Program) UniformLocation(name string) Uniform {
	n := gl.Str(name + "\x00")
	return Uniform(gl.GetUniformLocation(uint32(p), n))
}

func (u Uniform) Matrix4(transpose bool, matrix *vmath.Matrix4) {
	gl.UniformMatrix4fv(int32(u), 1, transpose, (*float32)(unsafe.Pointer(matrix)))
}

func (u Uniform) Int(val int) {
	gl.Uniform1i(int32(u), int32(val))
}

func (u Uniform) Float(val float32) {
	gl.Uniform1f(int32(u), val)
}

func (a Attribute) Enable() {
	gl.EnableVertexAttribArray(uint32(a))
}

func (a Attribute) Disable() {
	gl.DisableVertexAttribArray(uint32(a))
}

func (a Attribute) Pointer(size int, ty Type, normalized bool, stride, offset int) {
	gl.VertexAttribPointer(
		uint32(a),
		int32(size),
		uint32(ty),
		normalized,
		int32(stride),
		gl.PtrOffset(offset),
	)
}

type (
	Shader          uint32
	ShaderType      uint32
	ShaderParameter uint32
)

func CreateShader(t ShaderType) Shader {
	return Shader(gl.CreateShader(uint32(t)))
}

func (s Shader) Source(src string) {
	ss := gl.Str(src + "\x00")
	gl.ShaderSource(uint32(s), 1, &ss, nil)
}

func (s Shader) Compile() {
	gl.CompileShader(uint32(s))
}

func (s Shader) Parameter(param ShaderParameter) int {
	var p int32
	gl.GetShaderiv(uint32(s), uint32(param), &p)
	return int(p)
}

func (s Shader) InfoLog() string {
	l := s.Parameter(InfoLogLength)

	buf := make([]byte, l)

	gl.GetShaderInfoLog(uint32(s), int32(l), nil, (*uint8)(gl.Ptr(buf)))
	return string(buf)
}
