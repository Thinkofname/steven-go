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
	LightLevel        gl.Uniform   `gl:"lightLevel"`
}

const (
	vertex = `
#version 150
in vec3 aPosition;
in vec4 aTextureInfo;
in vec3 aTextureOffset;
in vec3 aColor;
in vec2 aLighting;

uniform mat4 perspectiveMatrix;
uniform mat4 cameraMatrix;
uniform ivec3 offset;
uniform float lightLevel;

out vec3 vColor;
out vec4 vTextureInfo;
out vec2 vTextureOffset;
out float vAtlas;
out float vLighting;
out float vDepth;

void main() {
	vec3 pos = vec3(aPosition.x, -aPosition.y, aPosition.z);
	vec3 o = vec3(offset.x, -offset.y, offset.z);
	gl_Position = perspectiveMatrix * cameraMatrix * vec4(pos + o * 16.0, 1.0);
	vDepth = (cameraMatrix * vec4(pos + o * 16.0, 1.0)).z;

	vColor = aColor;
	vTextureInfo = aTextureInfo;
	vTextureOffset = aTextureOffset.xy / 16.0;
	vAtlas = aTextureOffset.z;

	float light = max(aLighting.x, aLighting.y * 1.0) / 4000.0;
	vLighting = 1.0 * pow(lightLevel, 15.0 - light); // clamp(0.05 + pow(light / (4000.0 * 16.0), 1.5), 0.1, 1.0);
}
`
	fragment = `
#version 150

const float invAtlasSize = 1.0 / ` + atlasSizeStr + `;

uniform sampler2DArray textures;

in vec3 vColor;
in vec4 vTextureInfo;
in vec2 vTextureOffset;
in float vAtlas;
in float vLighting;
in float vDepth;

#ifndef alpha
out vec4 fragColor;
#else
out vec4 accum;
out float revealage;
#endif

void main() {
	vec2 tPos = vTextureOffset;
	tPos = mod(tPos, vTextureInfo.zw);
	tPos += vTextureInfo.xy;
	tPos *= invAtlasSize;
	vec4 col = texture(textures, vec3(tPos, vAtlas));
	#ifndef alpha
	if (col.a < 0.5) discard;
	#endif
	col *= vec4(vColor, 1.0);
	col.rgb *= vLighting;
	
	#ifndef alpha
	fragColor = col;
	#else
	float z = vDepth;
	float al = col.a;	
    float weight = pow(alpha + 0.01f, 4.0f) +
                   max(0.01f, min(3000.0f, 0.3f / (0.00001f + pow(abs(z) / 800.0f, 4.0f))));
	accum = vec4(col.rgb * al * weight, al);
	revealage = weight * al;
	#endif
}
`
)
