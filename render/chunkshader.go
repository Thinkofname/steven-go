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
in vec3 aPosition;
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
}
`
	fragment = `
#version 150

uniform sampler2D textures[5];

in vec3 vColor;
in vec4 vTextureInfo;
in vec2 vTextureOffset;
in float vLighting;

out vec4 fragColor;

void main() {
	vec2 tPos = vTextureOffset / 16.0;
	tPos = mod(tPos, vTextureInfo.zw);
	vec2 offset = vec2(vTextureInfo.x, mod(vTextureInfo.y, 1024.0));
	tPos += offset;
	tPos /= 1024.0;
 	vec4 col = texture2D(textures[int(floor(vTextureInfo.y / 1024.0))], vec2(tPos.x, tPos.y));
	if (col.a < 0.5) discard;
	col *= vec4(vColor, 1.0);
	col.rgb *= vLighting;
	fragColor = col;
}
`
)
