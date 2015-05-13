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
	X, Y, Z                    float32
	TX, TY, TW, TH             uint16
	TOffsetX, TOffsetY, TAtlas int16
	R, G, B, A                 byte
}

var staticFunc, staticTypes = builder.Struct(&StaticVertex{})

type StaticModel struct {
	// For culling only
	X, Y, Z float32
	Radius  float32
	// Per a part matrix
	Matrix     []mgl32.Mat4
	array      gl.VertexArray
	buffer     gl.Buffer
	bufferSize int
	count      int
	ranges     [][2]int
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
	staticState.shader.Position.Pointer(3, gl.Float, false, 30, 0)
	staticState.shader.TextureInfo.Pointer(4, gl.UnsignedShort, false, 30, 12)
	staticState.shader.TextureOffset.PointerInt(3, gl.Short, 30, 20)
	staticState.shader.Color.Pointer(4, gl.UnsignedByte, true, 30, 26)

	model.Matrix = make([]mgl32.Mat4, len(parts))
	model.ranges = make([][2]int, len(parts))
	var all []*StaticVertex
	for i, p := range parts {
		model.Matrix[i] = mgl32.Ident4()
		model.ranges[i] = [2]int{
			(len(all) / 4) * 6,
			(len(p) / 4) * 6,
		}
		all = append(all, p...)
	}
	model.data(all)

	staticState.models = append(staticState.models, model)
	return model
}

func (sm *StaticModel) data(verts []*StaticVertex) {
	sm.count = (len(verts) / 4) * 6
	if staticState.maxIndex < sm.count {
		var data []byte
		data, staticState.indexType = genElementBuffer(sm.count)
		staticState.indexBuffer.Bind(gl.ElementArrayBuffer)
		staticState.indexBuffer.Data(data, gl.DynamicDraw)
		staticState.maxIndex = sm.count
	}
	buf := builder.New(staticTypes...)
	for _, v := range verts {
		staticFunc(buf, v)
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
	for _, mdl := range staticState.models {
		if mdl.Radius != 0 && !frustum.IsSphereInside(mdl.X, mdl.Y, mdl.Z, mdl.Radius) {
			continue
		}
		mdl.array.Bind()
		for i := range mdl.Matrix {
			staticState.shader.ModelMatrix.Matrix4(&mdl.Matrix[i])
			gl.DrawElements(gl.Triangles, mdl.ranges[i][1], staticState.indexType, mdl.ranges[i][0]*m)
		}
	}

	gl.Disable(gl.Blend)
}

type staticShader struct {
	Position          gl.Attribute `gl:"aPosition"`
	TextureInfo       gl.Attribute `gl:"aTextureInfo"`
	TextureOffset     gl.Attribute `gl:"aTextureOffset"`
	Color             gl.Attribute `gl:"aColor"`
	PerspectiveMatrix gl.Uniform   `gl:"perspectiveMatrix"`
	CameraMatrix      gl.Uniform   `gl:"cameraMatrix"`
	ModelMatrix       gl.Uniform   `gl:"modelMatrix"`
	Texture           gl.Uniform   `gl:"textures"`
}

const (
	staticVertex = `
#version 150
in vec3 aPosition;
in vec4 aTextureInfo;
in ivec3 aTextureOffset;
in vec4 aColor;

uniform mat4 perspectiveMatrix;
uniform mat4 cameraMatrix;
uniform mat4 modelMatrix;

out vec4 vColor;
out vec4 vTextureInfo;
out vec2 vTextureOffset;
out float vAtlas;

void main() {
	vec3 pos = vec3(aPosition.x, -aPosition.y, aPosition.z);
	gl_Position = perspectiveMatrix * cameraMatrix * modelMatrix * vec4(pos, 1.0);
	vColor = aColor;
	vTextureInfo = aTextureInfo;
	vTextureOffset = aTextureOffset.xy / 16.0;
	vAtlas = aTextureOffset.z;
}
`
	staticFragment = `
#version 150

const float atlasSize = ` + atlasSizeStr + `;

uniform sampler2DArray textures;

in vec4 vColor;
in vec4 vTextureInfo;
in vec2 vTextureOffset;
in float vAtlas;

out vec4 fragColor;

void main() {
	vec2 tPos = vTextureOffset;
	tPos = mod(tPos, vTextureInfo.zw);
	tPos += vTextureInfo.xy;
	tPos /= atlasSize;
	vec4 col = texture(textures, vec3(tPos, vAtlas));
	if (col.a == 0.0) discard;
	col *= vColor;
	fragColor = col;
}
`
)
