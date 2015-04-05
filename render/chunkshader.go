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

import "github.com/thinkofdeath/steven/platform/gl"

type chunkShader struct {
	Position          gl.Attribute `gl:"aPosition"`
	TextureInfo       gl.Attribute `gl:"aTextureInfo"`
	TextureOffset     gl.Attribute `gl:"aTextureOffset"`
	Color             gl.Attribute `gl:"aColor"`
	Lighting          gl.Attribute `gl:"aLighting"`
	PerspectiveMatrix gl.Uniform   `gl:"perspectiveMatrix"`
	CameraMatrix      gl.Uniform   `gl:"cameraMatrix"`
	Offset            gl.Uniform   `gl:"offset"`
	Textures          gl.Uniform   `gl:"textures"`
}

var (
	vertex = `
#version 150
in ivec3 aPosition;
in vec4 aTextureInfo;
in vec2 aTextureOffset;
in vec3 aColor;
in vec2 aLighting;

uniform mat4 perspectiveMatrix;
uniform mat4 cameraMatrix;
uniform vec3 offset;

out vec3 vColor;
out vec4 vTextureInfo;
out vec2 vTextureOffset;
out float vLighting;
out float dist;

void main() {
	vec3 pos = vec3(aPosition.x, -aPosition.y, aPosition.z);
	vec3 o = vec3(offset.x, -offset.y, offset.z);
	gl_Position = perspectiveMatrix * cameraMatrix * vec4((pos / 256.0) + o * 16.0, 1.0);
	vColor = aColor;
	vTextureInfo = aTextureInfo;
	vTextureOffset = aTextureOffset;

	float light = max(aLighting.x, aLighting.y * 1.0);
	float val = pow(0.9, 16.0 - light) * 2.0;
	vLighting = clamp(pow(val, 1.5) * 0.5, 0.0, 1.0);
	dist = gl_Position.z;
}
`
	fragment = `
#version 150

uniform sampler2D textures[3];

in vec3 vColor;
in vec4 vTextureInfo;
in vec2 vTextureOffset;
in float vLighting;
in float dist;

out vec4 fragColor;

void main() {
	vec2 tPos = vTextureOffset / 16.0;
	tPos = mod(tPos, vTextureInfo.zw);
	vec2 offset = vec2(vTextureInfo.x, mod(vTextureInfo.y, 1024.0));
	tPos += offset;
	tPos /= 1024.0;
	int texID = int(floor(vTextureInfo.y / 1024.0));
	vec4 col = vec4(0.0);
	float mipLevel = max(0.0, (dist * 0.002) * 4.0);
	col += textureLod(textures[0], tPos, mipLevel) * float(0 == texID);
	col += textureLod(textures[1], tPos, mipLevel) * float(1 == texID);
	col += textureLod(textures[2], tPos, mipLevel) * float(2 == texID);
	#ifndef alpha
	if (col.a < 0.5) discard;
	#endif
	col *= vec4(vColor, 1.0);
	col.rgb *= vLighting;
	fragColor = col;
}
`
)
