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

package render

import (
	"reflect"

	"github.com/thinkofdeath/phteven/render/gl"
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
