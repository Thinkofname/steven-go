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

import "github.com/thinkofdeath/steven/render/glsl"

func init() {
	glsl.Register("lookup_texture", `
const float invAtlasSize = 1.0 / `+atlasSizeStr+`;
vec4 atlasTexture() {
	vec2 tPos = vTextureOffset;
	tPos = mod(tPos, vTextureInfo.zw);
	tPos += vTextureInfo.xy;
	tPos *= invAtlasSize;
	return texture(textures, vec3(tPos, vAtlas));
}	
`)
	glsl.Register("get_light", `
vec3 getLight(vec2 light) {
	vec2 li = pow(vec2(lightLevel), 15.0 - light) * 15.0 + 1.0;
	float bl = li.x;
	float sk = li.y;

	float br = (0.879552 * pow(bl, 2.0) + 0.871148 * bl + 32.9821);
	float bg = (1.22181 * pow(bl, 2.0) - 4.78113 * bl + 36.7125);
	float bb = (1.67612 * pow(bl, 2.0) - 12.9764 * bl + 48.8321);

	float sr = (0.131653 * pow(sk, 2.0) - 0.761625 * sk + 35.0393);
	float sg = (0.136555 * pow(sk, 2.0) - 0.853782 * sk + 29.6143);
	float sb = (0.327311 * pow(sk, 2.0) - 1.62017 * sk + 28.0929);
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
`)
}
