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
	SkyOffset         gl.Uniform   `gl:"skyOffset"`
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
uniform float skyOffset;

out vec3 vColor;
out vec4 vTextureInfo;
out vec2 vTextureOffset;
out float vAtlas;
out vec3 vLighting;
out float vDepth;

vec3 getLight(vec2 light);

void main() {
	vec3 pos = vec3(aPosition.x, -aPosition.y, aPosition.z);
	vec3 o = vec3(offset.x, -offset.y, offset.z);
	gl_Position = perspectiveMatrix * cameraMatrix * vec4(pos + o * 16.0, 1.0);
	vDepth = (cameraMatrix * vec4(pos + o * 16.0, 1.0)).z;

	vColor = aColor;
	vTextureInfo = aTextureInfo;
	vTextureOffset = aTextureOffset.xy / 16.0;
	vAtlas = aTextureOffset.z;

	vLighting = getLight(aLighting / (4000.0));
}

// TODO Pre compute this? 3D texture?
vec3 getLight(vec2 light) {
	vec2 li = pow(vec2(lightLevel), 15.0 - light) * 15.0 + 1.0;
	float bl = li.x;
	float sk = li.y;

	float br = (0.879552 * pow(bl, 2.0) + 0.871148 * bl + 32.9821);
	float bg = (1.22181 * pow(bl, 2.0) - 4.78113 * bl + 36.7125);
	float bb = (1.67612 * pow(bl, 2.0) - 12.9764 * bl + 48.8321);

	float sr = (0.131653 * pow(sk, 2.0) - 0.761625 * sk + 35.0393);
	float sg = (0.136555 * pow(sk, 2.0) - 0.853782 * sk + 29.6143);
	float sb = (0.277311 * pow(sk, 2.0) - 1.62017 * sk + 28.0929);
	float srl = (0.996148 * pow(sk, 2) - 4.19629 * sk + 51.4036);
	float sgl = (1.03904 * pow(sk, 2) - 4.81516 * sk + 47.0911);
	float sbl = (1.076164 * pow(sk, 2) - 5.36376 * sk + 43.9089);

	sr = srl * skyOffset + sr * (1.0 - skyOffset);
	sg = sgl * skyOffset + sg * (1.0 - skyOffset);
	sb = sbl * skyOffset + sb * (1.0 - skyOffset);

	return clamp(vec3(
		sqrt((br*br + sr*sr) / 2) / 255.0,
		sqrt((bg*bg + sg*sg) / 2) / 255.0,
		sqrt((bb*bb + sb*sb) / 2) / 255.0
	), 0.0, 1.0);
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
in vec3 vLighting;
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
