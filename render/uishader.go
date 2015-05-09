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

import "github.com/thinkofdeath/phteven/render/gl"

type uiShader struct {
	Position      gl.Attribute `gl:"aPosition"`
	TextureInfo   gl.Attribute `gl:"aTextureInfo"`
	TextureOffset gl.Attribute `gl:"aTextureOffset"`
	Color         gl.Attribute `gl:"aColor"`
	Texture       gl.Uniform   `gl:"textures"`
	ScreenSize    gl.Uniform   `gl:"screenSize"`
}

const (
	vertexUI = `
#version 150
in ivec2 aPosition;
in vec4 aTextureInfo;
in ivec3 aTextureOffset;
in vec4 aColor;

out vec4 vColor;
out vec4 vTextureInfo;
out vec2 vTextureOffset;
out float vAtlas;

uniform vec2 screenSize;

void main() {
	vec2 pos = aPosition / screenSize;
	gl_Position = vec4((pos.x-0.5)*2.0, -(pos.y-0.5)*2.0, 0.0, 1.0);
	vColor = aColor;
	vTextureInfo = aTextureInfo;
	vTextureOffset = aTextureOffset.xy / 16.0;
	vAtlas = aTextureOffset.z;
}
`
	fragmentUI = `
#version 150

const float atlasSize = ` + atlasSizeStr + `;

uniform sampler2DArray textures;

in vec4 vColor;
in vec4 vTextureInfo;
in vec2 vTextureOffset;
in float vAtlas;

out vec4 fragColor;

void main() {
	vec2 tPos = floor(vTextureOffset);
	tPos = mod(tPos, vTextureInfo.zw);
	tPos += vTextureInfo.xy;
	tPos /= atlasSize;
	vec4 col = texture(textures, vec3(tPos, vAtlas));
	col *= vColor;
	if (col.a == 0.0) discard;
	fragColor = col;
}
`
)
