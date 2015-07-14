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
	vec2 li = pow(vec2(lightLevel), 15.0 - light);
	float skyTint = skyOffset * 0.95 + 0.05;
	float bl = li.x;
	float sk = li.y * skyTint;

	float skyRed = sk * (skyOffset * 0.65 + 0.35);
	float skyGreen = sk * (skyOffset * 0.65 + 0.35);
	float blockGreen = bl * ((bl * 0.6 + 0.4) * 0.6 + 0.4);
	float blockBlue = bl * (bl * bl * 0.6 + 0.4);

	vec3 col = vec3(
		skyRed + bl,
		skyGreen + blockGreen,
		sk + blockBlue
	);

	col = col * 0.96 + 0.03;

	float gamma = 0.0;
	vec3 invCol = 1.0 - col;
	invCol = 1.0 - invCol * invCol * invCol * invCol;
	col = col * (1.0 - gamma) + invCol * gamma;
	col = col * 0.96 + 0.03;

	return clamp(col, 0.0, 1.0);
}	
`)
}
