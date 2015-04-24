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

import "github.com/thinkofdeath/steven/render/gl"

type chunkShader struct {
	Position          gl.Attribute `gl:"aPosition"`
	TextureInfo       gl.Attribute `gl:"aTextureInfo"`
	TextureOffset     gl.Attribute `gl:"aTextureOffset"`
	Color             gl.Attribute `gl:"aColor"`
	Lighting          gl.Attribute `gl:"aLighting"`
	PerspectiveMatrix gl.Uniform   `gl:"perspectiveMatrix"`
	CameraMatrix      gl.Uniform   `gl:"cameraMatrix"`
	Offset            gl.Uniform   `gl:"offset"`
	Texture           gl.Uniform   `gl:"textures"`
}

const (
	vertex = `
#version 150
in ivec3 aPosition;
in vec4 aTextureInfo;
in vec3 aTextureOffset;
in vec3 aColor;
in vec2 aLighting;

uniform mat4 perspectiveMatrix;
uniform mat4 cameraMatrix;
uniform ivec3 offset;

out vec3 vColor;
out vec4 vTextureInfo;
out vec3 vTextureOffset;
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

const float atlasSize = ` + atlasSizeStr + `;

uniform sampler2DArray textures;

in vec3 vColor;
in vec4 vTextureInfo;
in vec3 vTextureOffset;
in float vLighting;
in float dist;

out vec4 fragColor;

void main() {
	vec2 tPos = vTextureOffset.xy / 16.0;
	tPos = mod(tPos, vTextureInfo.zw);
	vec2 offset = vTextureInfo.xy;
	tPos += offset;
	tPos /= atlasSize;
	vec4 col = texture(textures, vec3(tPos, vTextureOffset.z));
	#ifndef alpha
	if (col.a < 0.5) discard;
	#endif
	col *= vec4(vColor, 1.0);
	col.rgb *= vLighting;
	fragColor = col;
}
`
)
