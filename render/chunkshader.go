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
out vec2 vTextureOffset;
out float vAtlas;
out float vLighting;
out float vLogDepth;

const float C = 0.01;
const float FC = 1.0/log(500.0*C + 1);

void main() {
	ivec3 pos = ivec3(aPosition.x, -aPosition.y, aPosition.z);
	vec3 o = vec3(offset.x, -offset.y, offset.z);
	gl_Position = perspectiveMatrix * cameraMatrix * vec4((pos / 256.0) + o * 16.0, 1.0);

	vLogDepth = log(gl_Position.w*C + 1)*FC;
	gl_Position.z = (2*vLogDepth - 1)*gl_Position.w;

	vColor = aColor;
	vTextureInfo = aTextureInfo;
	vTextureOffset = aTextureOffset.xy / 16.0;
	vAtlas = aTextureOffset.z;

	float light = max(aLighting.x, aLighting.y * 1.0);
	vLighting = clamp(0.05 + pow(light / (4000.0 * 16.0), 1.5), 0.1, 1.0);
}
`
	fragment = `
#version 150
#ifdef GL_ARB_conservative_depth
#extension GL_ARB_conservative_depth : enable
layout(depth_less) out float gl_FragDepth;
#endif

const float atlasSize = ` + atlasSizeStr + `;

uniform sampler2DArray textures;

in vec3 vColor;
in vec4 vTextureInfo;
in vec2 vTextureOffset;
in float vAtlas;
in float vLighting;
in float vLogDepth;

out vec4 fragColor;

void main() {
	gl_FragDepth = vLogDepth;
	vec2 tPos = vTextureOffset;
	tPos = mod(tPos, vTextureInfo.zw);
	tPos += vTextureInfo.xy;
	tPos /= atlasSize;
	vec4 col = texture(textures, vec3(tPos, vAtlas));
	#ifndef alpha
	if (col.a < 0.5) discard;
	#endif
	col *= vec4(vColor, 1.0);
	col.rgb *= vLighting;
	fragColor = col;
}
`
)
