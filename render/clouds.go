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
	"fmt"

	"github.com/thinkofdeath/steven/render/builder"
	"github.com/thinkofdeath/steven/render/gl"
	"github.com/thinkofdeath/steven/render/glsl"
)

var clouds cloudState
var DrawClouds = true

type cloudState struct {
	program gl.Program
	shader  cloudShader

	array  gl.VertexArray
	buffer gl.Buffer

	texture gl.Texture
	data    []byte

	numPoints int
}

type cloudShader struct {
	Position          gl.Attribute `gl:"aPosition"`
	PerspectiveMatrix gl.Uniform   `gl:"perspectiveMatrix"`
	CameraMatrix      gl.Uniform   `gl:"cameraMatrix"`
	LightLevel        gl.Uniform   `gl:"lightLevel"`
	SkyOffset         gl.Uniform   `gl:"skyOffset"`
	Offset            gl.Uniform   `gl:"offset"`
	CloudMap          gl.Uniform   `gl:"cloudMap"`
}

func CloudData() []byte {
	return clouds.data
}

func (c *cloudState) init() {
	program := gl.CreateProgram()

	v := gl.CreateShader(gl.VertexShader)
	v.Source(glsl.Get("cloud_vertex"))
	v.Compile()

	if v.Parameter(gl.CompileStatus) == 0 {
		panic(v.InfoLog())
	} else {
		log := v.InfoLog()
		if len(log) > 0 {
			fmt.Println(log)
		}
	}

	g := gl.CreateShader(gl.GeometryShader)
	g.Source(glsl.Get("cloud_geo"))
	g.Compile()

	if g.Parameter(gl.CompileStatus) == 0 {
		panic(g.InfoLog())
	} else {
		log := g.InfoLog()
		if len(log) > 0 {
			fmt.Println(log)
		}
	}

	f := gl.CreateShader(gl.FragmentShader)
	f.Source(glsl.Get("cloud_fragment"))
	f.Compile()

	if f.Parameter(gl.CompileStatus) == 0 {
		panic(f.InfoLog())
	} else {
		log := f.InfoLog()
		if len(log) > 0 {
			fmt.Println(log)
		}
	}

	program.AttachShader(v)
	program.AttachShader(g)
	program.AttachShader(f)
	program.Link()
	program.Use()

	c.program = program

	InitStruct(&c.shader, c.program)

	c.array = gl.CreateVertexArray()
	c.array.Bind()
	c.buffer = gl.CreateBuffer()
	c.buffer.Bind(gl.ArrayBuffer)
	c.shader.Position.Enable()
	c.shader.Position.Pointer(3, gl.Float, false, 12, 0)

	data := builder.New(builder.Float, builder.Float, builder.Float)
	for x := -160; x < 160; x++ {
		for y := -160; y < 160; y++ {
			data.Float(float32(x))
			data.Float(128)
			data.Float(float32(y))
			c.numPoints++
		}
	}

	c.buffer.Data(data.Data(), gl.StaticDraw)

	c.data = make([]byte, 512*512)
	c.texture = gl.CreateTexture()
	c.texture.Bind(gl.Texture2D)
	c.texture.Image2D(0, 512, 512, gl.Red, gl.UnsignedByte, c.data)
	c.texture.Parameter(gl.TextureMinFilter, gl.Nearest)
	c.texture.Parameter(gl.TextureMagFilter, gl.Nearest)
}

func (c *cloudState) tick(delta float64) {
	if !DrawClouds {
		return
	}
	c.program.Use()
	c.shader.PerspectiveMatrix.Matrix4(&perspectiveMatrix)
	c.shader.CameraMatrix.Matrix4(&cameraMatrix)
	c.shader.SkyOffset.Float(SkyOffset)
	c.shader.LightLevel.Float(LightLevel)
	c.shader.Offset.Float3(float32(int(Camera.X)), 0, float32(int(Camera.Z)))

	gl.ActiveTexture(0)
	c.texture.Bind(gl.Texture2D)
	c.texture.SubImage2D(0, 0, 0, 512, 512, gl.Red, gl.UnsignedByte, c.data)
	c.shader.CloudMap.Int(0)
	c.array.Bind()
	gl.DrawArrays(gl.Points, 0, c.numPoints)
}

func init() {
	glsl.Register("cloud_vertex", `	
in vec3 aPosition;

uniform float lightLevel;
uniform float skyOffset;

out vec3 vLighting;

#include get_light

void main() {
	vec3 pos = vec3(aPosition.x, -aPosition.y, aPosition.z);
	gl_Position = vec4(pos, 1.0);

	vLighting = getLight(vec2(0.0, 15.0));
}
`)
	glsl.Register("cloud_geo", `
layout(points) in;
layout(triangle_strip, max_vertices = 24) out;

uniform mat4 perspectiveMatrix;
uniform mat4 cameraMatrix;
uniform vec3 offset;
uniform sampler2D cloudMap;

in vec3 vLighting[];

out vec3 fLighting;
out vec4 fColor;

void setVertex(vec3 base, vec3 off, float color) {	
	gl_Position = perspectiveMatrix * cameraMatrix * vec4(base + off*vec3(1.0,-1.0,1.0), 1.0);
	fColor = vec4(color, color, color, 1.0);
	fLighting = vLighting[0];
	EmitVertex();
}

bool isSolid(ivec2 pos) {
	return texelFetch(cloudMap, pos, 0).r > 0.5;	
}

void main() {
	vec3 base = floor(offset) + gl_in[0].gl_Position.xyz;
	ivec2 texP = ivec2(gl_in[0].gl_Position.xz + 160.0);
	
	if (!isSolid(texP)) return;
	
	// Top
	setVertex(base, vec3(0.0, 1.0, 0.0), 1.0);
	setVertex(base, vec3(1.0, 1.0, 0.0), 1.0);
	setVertex(base, vec3(0.0, 1.0, 1.0), 1.0);
	setVertex(base, vec3(1.0, 1.0, 1.0), 1.0);
	EndPrimitive();	
	
	// Bottom
	setVertex(base, vec3(0.0, 0.0, 0.0), 0.7);
	setVertex(base, vec3(0.0, 0.0, 1.0), 0.7);
	setVertex(base, vec3(1.0, 0.0, 0.0), 0.7);
	setVertex(base, vec3(1.0, 0.0, 1.0), 0.7);
	EndPrimitive();	
	
	if (!isSolid(texP + ivec2(-1, 0))) {
		// -X
		setVertex(base, vec3(0.0, 0.0, 0.0), 0.8);
		setVertex(base, vec3(0.0, 1.0, 0.0), 0.8);
		setVertex(base, vec3(0.0, 0.0, 1.0), 0.8);
		setVertex(base, vec3(0.0, 1.0, 1.0), 0.8);
		EndPrimitive();
	}	
	
	if (!isSolid(texP + ivec2(1, 0))) {
		// +X
		setVertex(base, vec3(1.0, 0.0, 0.0), 0.8);
		setVertex(base, vec3(1.0, 0.0, 1.0), 0.8);
		setVertex(base, vec3(1.0, 1.0, 0.0), 0.8);
		setVertex(base, vec3(1.0, 1.0, 1.0), 0.8);
		EndPrimitive();
	}
	
	if (!isSolid(texP + ivec2(0, 1))) {
		// -Z
		setVertex(base, vec3(0.0, 0.0, 1.0), 0.8);
		setVertex(base, vec3(0.0, 1.0, 1.0), 0.8);
		setVertex(base, vec3(1.0, 0.0, 1.0), 0.8);
		setVertex(base, vec3(1.0, 1.0, 1.0), 0.8);
		EndPrimitive();
	}
	
	if (!isSolid(texP + ivec2(0, -1))) {
		// +Z
		setVertex(base, vec3(0.0, 0.0, 0.0), 0.8);
		setVertex(base, vec3(1.0, 0.0, 0.0), 0.8);
		setVertex(base, vec3(0.0, 1.0, 0.0), 0.8);
		setVertex(base, vec3(1.0, 1.0, 0.0), 0.8);
		EndPrimitive();
	}
}
`)
	glsl.Register("cloud_fragment", `
in vec4 fColor;
in vec3 fLighting;

out vec4 fragColor;

void main() {
	vec4 col = fColor;
	col.rgb *= fLighting;
	fragColor = col;
}
	`)
}
