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
	"github.com/go-gl/mathgl/mgl32"
	"github.com/thinkofdeath/steven/render/builder"
	"github.com/thinkofdeath/steven/render/gl"
	"github.com/thinkofdeath/steven/render/glsl"
)

// Really static is a bad name for this but i'm
// too lazy to think of a better one

var staticState = struct {
	models         []*staticCollection
	baseCollection *staticCollection
	indexBuffer    gl.Buffer
	indexType      gl.Type
	maxIndex       int
}{}

type staticCollection struct {
	program gl.Program
	shader  *staticShader
	models  []*StaticModel
}

type StaticVertex struct {
	X, Y, Z            float32
	Texture            TextureInfo
	TextureX, TextureY float64
	R, G, B, A         byte
	id                 byte
}
type staticVertexInternal struct {
	X, Y, Z                    float32
	TX, TY, TW, TH             uint16
	TOffsetX, TOffsetY, TAtlas int16
	Pad0                       int16
	R, G, B, A                 byte
	ID                         byte
	Pad1, Pad2, Pad3           byte
}

var staticFunc, staticTypes = builder.Struct(&staticVertexInternal{})

type StaticModel struct {
	// For culling only
	X, Y, Z float32
	Radius  float32
	// Per a part matrix
	Matrix     []mgl32.Mat4
	Colors     [][4]float32
	array      gl.VertexArray
	buffer     gl.Buffer
	bufferSize int
	count      int

	counts  []int32
	offsets []uintptr

	all []*StaticVertex

	c *staticCollection
}

func NewStaticModel(parts [][]*StaticVertex) *StaticModel {
	return newStaticModel(parts, staticState.baseCollection)
}

func newStaticModel(parts [][]*StaticVertex, c *staticCollection) *StaticModel {
	model := &StaticModel{
		c: c,
	}

	model.array = gl.CreateVertexArray()
	model.array.Bind()
	staticState.indexBuffer.Bind(gl.ElementArrayBuffer)
	model.buffer = gl.CreateBuffer()
	model.buffer.Bind(gl.ArrayBuffer)
	c.shader.Position.Enable()
	c.shader.TextureInfo.Enable()
	c.shader.TextureOffset.Enable()
	c.shader.Color.Enable()
	c.shader.ID.Enable()
	c.shader.Position.Pointer(3, gl.Float, false, 36, 0)
	c.shader.TextureInfo.Pointer(4, gl.UnsignedShort, false, 36, 12)
	c.shader.TextureOffset.PointerInt(3, gl.Short, 36, 20)
	c.shader.Color.Pointer(4, gl.UnsignedByte, true, 36, 28)
	c.shader.ID.PointerInt(4, gl.UnsignedByte, 36, 32)

	model.Matrix = make([]mgl32.Mat4, len(parts))
	model.Colors = make([][4]float32, len(parts))
	model.counts = make([]int32, len(parts))
	model.offsets = make([]uintptr, len(parts))
	var all []*StaticVertex
	for i, p := range parts {
		model.Matrix[i] = mgl32.Ident4()
		model.Colors[i] = [4]float32{1.0, 1.0, 1.0, 1.0}
		model.counts[i] = int32((len(p) / 4) * 6)
		model.offsets[i] = uintptr((len(all) / 4) * 6)
		for _, pp := range p {
			pp.id = byte(i)
		}
		all = append(all, p...)
	}
	model.all = all
	model.data()

	c.models = append(c.models, model)
	return model
}

func (sm *StaticModel) data() {
	verts := sm.all
	sm.array.Bind()
	sm.count = (len(verts) / 4) * 6
	if staticState.maxIndex < sm.count {
		var data []byte
		data, staticState.indexType = genElementBuffer(sm.count)
		staticState.indexBuffer.Bind(gl.ElementArrayBuffer)
		staticState.indexBuffer.Data(data, gl.DynamicDraw)
		staticState.maxIndex = sm.count
	}
	buf := builder.New(staticTypes...)
	in := staticVertexInternal{}
	for _, v := range verts {
		in.X = v.X
		in.Y = v.Y
		in.Z = v.Z
		rect := v.Texture.Rect()
		in.TX = uint16(rect.X)
		in.TY = uint16(rect.Y)
		in.TW = uint16(rect.Width)
		in.TH = uint16(rect.Height)
		in.TOffsetX = int16(16 * float64(rect.Width) * v.TextureX)
		in.TOffsetY = int16(16 * float64(rect.Height) * v.TextureY)
		in.TAtlas = int16(v.Texture.Atlas())
		in.R = v.R
		in.G = v.G
		in.B = v.B
		in.A = v.A
		in.ID = v.id
		staticFunc(buf, &in)
	}
	sm.buffer.Bind(gl.ArrayBuffer)
	data := buf.Data()
	if len(data) < sm.bufferSize {
		gb := sm.buffer.Map(gl.WriteOnly, sm.bufferSize)
		copy(gb, data)
		sm.buffer.Unmap()
	} else {
		sm.buffer.Data(data, gl.DynamicDraw)
		sm.bufferSize = len(data)
	}
}

func (sm *StaticModel) Free() {
	sm.array.Delete()
	sm.buffer.Delete()
	c := sm.c
	for i, s := range c.models {
		if s == sm {
			c.models = append(c.models[:i], c.models[i+1:]...)
			return
		}
	}
}

func RefreshStaticModels() {
	for _, c := range staticState.models {
		for _, mdl := range c.models {
			mdl.data()
		}
	}
}

func initStatic() {
	c := &staticCollection{}
	c.program = CreateProgram(glsl.Get("static_vertex"), glsl.Get("static_frag"))
	c.shader = &staticShader{}
	InitStruct(c.shader, c.program)

	staticState.baseCollection = c
	staticState.models = append(staticState.models, c)

	staticState.indexBuffer = gl.CreateBuffer()
}

func drawStatic() {
	if len(staticState.models) == 0 {
		return
	}
	m := 4
	if staticState.indexType == gl.UnsignedShort {
		m = 2
	}

	gl.Enable(gl.Blend)

	offsetBuf := make([]uintptr, 10)

	for _, c := range staticState.models {
		c.program.Use()
		c.shader.Texture.Int(0)
		c.shader.PerspectiveMatrix.Matrix4(&perspectiveMatrix)
		c.shader.CameraMatrix.Matrix4(&cameraMatrix)
		for _, mdl := range c.models {
			if mdl.Radius != 0 && !frustum.IsSphereInside(mdl.X, mdl.Y, mdl.Z, mdl.Radius) {
				continue
			}
			mdl.array.Bind()
			if len(mdl.counts) > 1 {
				copy(offsetBuf, mdl.offsets)
				for i := range mdl.offsets {
					offsetBuf[i] *= uintptr(m)
				}
				c.shader.ModelMatrix.Matrix4Multi(mdl.Matrix)
				c.shader.ColorMul.FloatMutliRaw(mdl.Colors, len(mdl.Colors))
				gl.MultiDrawElements(gl.Triangles, mdl.counts, staticState.indexType, offsetBuf[:len(mdl.offsets)])
			} else {
				c.shader.ModelMatrix.Matrix4Multi(mdl.Matrix)
				c.shader.ColorMul.FloatMutliRaw(mdl.Colors, len(mdl.Colors))
				gl.DrawElements(gl.Triangles, int(mdl.counts[0]), staticState.indexType, int(mdl.offsets[0])*m)
			}
		}
	}

	gl.Disable(gl.Blend)
}

type staticShader struct {
	Position          gl.Attribute `gl:"aPosition"`
	TextureInfo       gl.Attribute `gl:"aTextureInfo"`
	TextureOffset     gl.Attribute `gl:"aTextureOffset"`
	Color             gl.Attribute `gl:"aColor"`
	ID                gl.Attribute `gl:"id"`
	PerspectiveMatrix gl.Uniform   `gl:"perspectiveMatrix"`
	CameraMatrix      gl.Uniform   `gl:"cameraMatrix"`
	ModelMatrix       gl.Uniform   `gl:"modelMatrix[]"`
	Texture           gl.Uniform   `gl:"textures"`
	ColorMul          gl.Uniform   `gl:"colorMul[]"`
}

func init() {
	glsl.Register("static_vertex", `
in vec3 aPosition;
in vec4 aTextureInfo;
in ivec3 aTextureOffset;
in vec4 aColor;
in int id;

uniform mat4 perspectiveMatrix;
uniform mat4 cameraMatrix;
uniform mat4 modelMatrix[10];

out vec4 vColor;
out vec4 vTextureInfo;
out vec2 vTextureOffset;
out float vAtlas;
out float vID;

void main() {
	vec3 pos = vec3(aPosition.x, -aPosition.y, aPosition.z);
	gl_Position = perspectiveMatrix * cameraMatrix * modelMatrix[id] * vec4(pos, 1.0);

	vColor = aColor;
	vTextureInfo = aTextureInfo;
	vTextureOffset = aTextureOffset.xy / 16.0;
	vAtlas = aTextureOffset.z;
	vID = float(id);
}
`)
	glsl.Register("static_frag", `

uniform sampler2DArray textures;
uniform vec4 colorMul[10];

in vec4 vColor;
in vec4 vTextureInfo;
in vec2 vTextureOffset;
in float vAtlas;
in float vID;

out vec4 fragColor;

#include lookup_texture

void main() {
	vec4 col = atlasTexture();
	if (col.a <= 0.05) discard;
	col *= vColor;
	fragColor = col * colorMul[int(vID)];
}
`)
}
