// +build mobile

package gl

import (
	"github.com/thinkofdeath/steven/vmath"
	"golang.org/x/mobile/gl"
)

const (
	VertexShader   ShaderType = gl.VERTEX_SHADER
	FragmentShader ShaderType = gl.FRAGMENT_SHADER

	CompileStatus ShaderParameter = gl.COMPILE_STATUS
	InfoLogLength ShaderParameter = gl.INFO_LOG_LENGTH
)

type (
	Program   gl.Program
	Attribute gl.Attrib
	Uniform   gl.Uniform
)

func CreateProgram() Program {
	return Program(gl.CreateProgram())
}

func (p Program) AttachShader(s Shader) {
	gl.AttachShader(gl.Program(p), gl.Shader(s))
}

func (p Program) Link() {
	gl.LinkProgram(gl.Program(p))
}

var (
	currentProgram Program
)

func (p Program) Use() {
	if p == currentProgram {
		return
	}
	gl.UseProgram(gl.Program(p))
	currentProgram = p
}
func (p Program) AttributeLocation(name string) Attribute {
	return Attribute(gl.GetAttribLocation(gl.Program(p), name))
}

func (p Program) UniformLocation(name string) Uniform {
	return Uniform(gl.GetUniformLocation(gl.Program(p), name))
}

func (u Uniform) Matrix4(transpose bool, matrix *vmath.Matrix4) {
	gl.UniformMatrix4fv(gl.Uniform(u), matrix[:])
}

func (u Uniform) Int(val int) {
	gl.Uniform1i(gl.Uniform(u), val)
}

func (u Uniform) Float(val float32) {
	gl.Uniform1f(gl.Uniform(u), val)
}

func (a Attribute) Enable() {
	gl.EnableVertexAttribArray(gl.Attrib(a))
}

func (a Attribute) Disable() {
	gl.DisableVertexAttribArray(gl.Attrib(a))
}

func (a Attribute) Pointer(size int, ty Type, normalized bool, stride, offset int) {
	gl.VertexAttribPointer(
		gl.Attrib(a),
		size,
		gl.Enum(ty),
		normalized,
		stride,
		offset,
	)
}

type (
	Shader          gl.Shader
	ShaderType      uint32
	ShaderParameter uint32
)

func CreateShader(t ShaderType) Shader {
	return Shader(gl.CreateShader(gl.Enum(t)))
}

func (s Shader) Source(src string) {
	gl.ShaderSource(gl.Shader(s), src)
}

func (s Shader) Compile() {
	gl.CompileShader(gl.Shader(s))
}

func (s Shader) Parameter(param ShaderParameter) int {
	return gl.GetShaderi(gl.Shader(s), gl.Enum(param))
}

func (s Shader) InfoLog() string {
	return gl.GetShaderInfoLog(gl.Shader(s))
}
