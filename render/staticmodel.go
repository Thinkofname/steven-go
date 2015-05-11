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
	"github.com/thinkofdeath/steven/render/gl"
)

var staticState = struct {
	program gl.Program
	shader  *staticShader
	models  []*StaticModel
}{}

type StaticModel struct {
	Matrix mgl32.Mat4
}

func initStatic() {
	staticState.program = CreateProgram(staticVertex, staticFragment)
	staticState.shader = &staticShader{}
	InitStruct(staticState.shader, staticState.program)
}

func drawStatic() {

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
in vec3 aTextureOffset;
in vec3 aColor;

uniform mat4 perspectiveMatrix;
uniform mat4 cameraMatrix;
uniform mat4 modelMatrix;

out vec3 vColor;
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

in vec3 vColor;
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
	col *= vec4(vColor, 1.0);
	fragColor = col;
}
`
)
