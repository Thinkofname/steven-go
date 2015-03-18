package render

import (
	"github.com/thinkofdeath/steven/platform/gl"
	"reflect"
)

// CreateProgram creates an OpenGL shader program from the
// passed shader sources. Panics if the shader is invalid.
func CreateProgram(vertex, fragment string) gl.Program {
	program := gl.CreateProgram()

	v := gl.CreateShader(gl.VertexShader)
	v.Source(vertex)
	v.Compile()

	if v.Parameter(gl.CompileStatus) == 0 {
		panic(v.InfoLog())
	}

	f := gl.CreateShader(gl.FragmentShader)
	f.Source(fragment)
	f.Compile()

	if f.Parameter(gl.CompileStatus) == 0 {
		panic(f.InfoLog())
	}

	program.AttachShader(v)
	program.AttachShader(f)
	program.Link()
	program.Use()
	return program
}

// InitStruct loads the Uniform and Attribute fields
// of a struct with values from the program using
// the name from the 'gl' struct tag.
func InitStruct(inter interface{}, program gl.Program) {
	v := reflect.ValueOf(inter).Elem()
	t := v.Type()
	l := t.NumField()

	var a gl.Attribute
	var u gl.Uniform
	gla := reflect.TypeOf(a)
	glu := reflect.TypeOf(u)

	for i := 0; i < l; i++ {
		f := t.Field(i)
		name := f.Tag.Get("gl")
		if f.Type == gla {
			v.Field(i).Set(reflect.ValueOf(
				program.AttributeLocation(name),
			))
		} else if f.Type == glu {
			v.Field(i).Set(reflect.ValueOf(
				program.UniformLocation(name),
			))
		}
	}
}
