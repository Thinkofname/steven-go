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
)

var staticState = struct {
	program     gl.Program
	shader      *staticShader
	models      []*StaticModel
	indexBuffer gl.Buffer
	indexType   gl.Type
	maxIndex    int
}{}

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
	R, G, B, A                 byte
	ID                         byte
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
}

func NewStaticModel(parts [][]*StaticVertex) *StaticModel {
	model := &StaticModel{}

	model.array = gl.CreateVertexArray()
	model.array.Bind()
	staticState.indexBuffer.Bind(gl.ElementArrayBuffer)
	model.buffer = gl.CreateBuffer()
	model.buffer.Bind(gl.ArrayBuffer)
	staticState.shader.Position.Enable()
	staticState.shader.TextureInfo.Enable()
	staticState.shader.TextureOffset.Enable()
	staticState.shader.Color.Enable()
	staticState.shader.ID.Enable()
	staticState.shader.Position.Pointer(3, gl.Float, false, 31, 0)
	staticState.shader.TextureInfo.Pointer(4, gl.UnsignedShort, false, 31, 12)
	staticState.shader.TextureOffset.PointerInt(3, gl.Short, 31, 20)
	staticState.shader.Color.Pointer(4, gl.UnsignedByte, true, 31, 26)
	staticState.shader.ID.PointerInt(1, gl.UnsignedByte, 31, 30)

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

	staticState.models = append(staticState.models, model)
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
	for i, s := range staticState.models {
		if s == sm {
			staticState.models = append(staticState.models[:i], staticState.models[i+1:]...)
			return
		}
	}
}

func RefreshStaticModels() {
	for _, mdl := range staticState.models {
		mdl.data()
	}
}

func initStatic() {
	staticState.program = CreateProgram(staticVertex, staticFragment)
	staticState.shader = &staticShader{}
	InitStruct(staticState.shader, staticState.program)

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
	staticState.program.Use()
	staticState.shader.Texture.Int(0)
	staticState.shader.PerspectiveMatrix.Matrix4(&perspectiveMatrix)
	staticState.shader.CameraMatrix.Matrix4(&cameraMatrix)

	offsetBuf := make([]uintptr, 10)

	for _, mdl := range staticState.models {
		if mdl.Radius != 0 && !frustum.IsSphereInside(mdl.X, mdl.Y, mdl.Z, mdl.Radius) {
			continue
		}
		mdl.array.Bind()
		if len(mdl.counts) > 1 {
			copy(offsetBuf, mdl.offsets)
			for i := range mdl.offsets {
				offsetBuf[i] *= uintptr(m)
			}
			staticState.shader.ModelMatrix.Matrix4Multi(mdl.Matrix)
			staticState.shader.ColorMul.FloatMutliRaw(mdl.Colors, len(mdl.Colors))
			gl.MultiDrawElements(gl.Triangles, mdl.counts, staticState.indexType, offsetBuf[:len(mdl.offsets)])
		} else {
			staticState.shader.ModelMatrix.Matrix4Multi(mdl.Matrix)
			staticState.shader.ColorMul.FloatMutliRaw(mdl.Colors, len(mdl.Colors))
			gl.DrawElements(gl.Triangles, int(mdl.counts[0]), staticState.indexType, int(mdl.offsets[0])*m)
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

const (
	staticVertex = `
#version 150
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
out float vLogDepth;
out float vID;

const float C = 0.01;
const float FC = 1.0/log(500.0*C + 1);

void main() {
	vec3 pos = vec3(aPosition.x, -aPosition.y, aPosition.z);
	gl_Position = perspectiveMatrix * cameraMatrix * modelMatrix[id] * vec4(pos, 1.0);

	vLogDepth = log(gl_Position.w*C + 1)*FC;
	gl_Position.z = (2*vLogDepth - 1)*gl_Position.w;

	vColor = aColor;
	vTextureInfo = aTextureInfo;
	vTextureOffset = aTextureOffset.xy / 16.0;
	vAtlas = aTextureOffset.z;
	vID = float(id);
}
`
	staticFragment = `
#version 150
#ifdef GL_ARB_conservative_depth
#extension GL_ARB_conservative_depth : enable
layout(depth_less) out float gl_FragDepth;
#endif

const float atlasSize = ` + atlasSizeStr + `;

uniform sampler2DArray textures;
uniform vec4 colorMul[10];

in vec4 vColor;
in vec4 vTextureInfo;
in vec2 vTextureOffset;
in float vAtlas;
in float vLogDepth;
in float vID;

out vec4 fragColor;

void main() {
	gl_FragDepth = vLogDepth;
	vec2 tPos = vTextureOffset;
	tPos = mod(tPos, vTextureInfo.zw);
	tPos += vTextureInfo.xy;
	tPos /= atlasSize;
	vec4 col = texture(textures, vec3(tPos, vAtlas));
	if (col.a <= 0.05) discard;
	col *= vColor;
	fragColor = col * colorMul[int(vID)];
}
`
)
