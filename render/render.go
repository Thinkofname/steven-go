package render

import (
	"math"

	"github.com/thinkofdeath/steven/platform"
	"github.com/thinkofdeath/steven/platform/gl"
	"github.com/thinkofdeath/steven/vmath"
)

var (
	chunkProgram gl.Program
	shaderChunk  *chunkShader

	lastWidth, lastHeight int = -1, -1
	perspectiveMatrix         = vmath.NewMatrix4()
	cameraMatrix              = vmath.NewMatrix4()

	syncChan = make(chan func(), 500)

	glTextures []gl.Texture
)

// Start starts the renderer
func Start() {
	gl.ClearColor(0.0, 1.0, 1.0, 1.0)
	gl.Enable(gl.DepthTest)
	gl.Enable(gl.CullFaceFlag)
	gl.CullFace(gl.Back)
	gl.FrontFace(gl.ClockWise)

	chunkProgram = CreateProgram(vertex, fragment)
	shaderChunk = &chunkShader{}
	InitStruct(shaderChunk, chunkProgram)

	loadTextures()

	for _, tex := range textures {
		glTextures = append(glTextures, createTexture(glTexture{
			Data:  tex.Buffer,
			Width: atlasSize, Height: atlasSize,
			Format: gl.RGBA,
		}))
	}
}

var textureIds []int

// Draw draws a single frame
func Draw() {
sync:
	for {
		select {
		case f := <-syncChan:
			f()
		default:
			break sync
		}
	}

	width, height := platform.Size()
	if lastHeight != height || lastWidth != width {
		lastWidth = width
		lastHeight = height

		perspectiveMatrix.Identity()
		perspectiveMatrix.Perspective(
			(math.Pi/180)*90,
			float32(width)/float32(height),
			0.1,
			10000.0,
		)
		gl.Viewport(0, 0, width, height)
	}
	if len(textureIds) != len(glTextures) {
		textureIds = make([]int, len(glTextures))
		for i, tex := range glTextures {
			tex.Bind(gl.Texture2D)
			gl.ActiveTexture(i)
			textureIds[i] = i
		}
	}

	gl.Clear(gl.ColorBufferBit | gl.DepthBufferBit)

	chunkProgram.Use()
	shaderChunk.PerspectiveMatrix.Matrix4(perspectiveMatrix)

	cameraMatrix.Identity()
	cameraMatrix.Translate(float32(Camera.X), float32(Camera.Y+1.62), float32(-Camera.Z))
	cameraMatrix.RotateY(float32(Camera.Yaw))
	cameraMatrix.RotateX(float32(Camera.Pitch))
	cameraMatrix.Scale(-1.0, 1.0, 1.0)

	shaderChunk.CameraMatrix.Matrix4(cameraMatrix)
	shaderChunk.Textures.IntV(textureIds...)

	for _, chunk := range buffers {
		if chunk.count == 0 {
			continue
		}
		shaderChunk.Offset.Float3(float32(chunk.X), float32(chunk.Y), float32(chunk.Z))

		chunk.array.Bind()
		gl.DrawArrays(gl.Triangles, 0, chunk.count)
	}
}

func renderSync(f func()) {
	syncChan <- f
}

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
attribute vec3 aPosition;
attribute vec4 aTextureInfo;
attribute vec2 aTextureOffset;
attribute vec3 aColor;
attribute vec2 aLighting;

uniform mat4 perspectiveMatrix;
uniform mat4 cameraMatrix;
uniform vec3 offset;

varying vec3 vColor;
varying vec4 vTextureInfo;
varying vec2 vTextureOffset;
varying float vLighting;

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

uniform sampler2D textures[5];

varying vec3 vColor;
varying vec4 vTextureInfo;
varying vec2 vTextureOffset;
varying float vLighting;

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
	gl_FragColor = col;
}
`
)
