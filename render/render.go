package render

import (
	"math"

	"github.com/thinkofdeath/steven/platform"
	"github.com/thinkofdeath/steven/platform/gl"
	"github.com/thinkofdeath/steven/vmath"
)

var (
	testProgram gl.Program
	test        *testShader

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

	testProgram = CreateProgram(vertex, fragment)
	test = &testShader{}
	InitStruct(test, testProgram)

	loadTextures()

	for _, tex := range textures {
		glTextures = append(glTextures, createTexture(glTexture{
			Data:  tex.Buffer,
			Width: atlasSize, Height: atlasSize,
			Format: gl.RGBA,
		}))
	}
}

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
			(math.Pi/180)*75,
			float32(width)/float32(height),
			0.1,
			10000.0,
		)
	}

	gl.Clear(gl.ColorBufferBit | gl.DepthBufferBit)
	gl.Viewport(0, 0, width, height)

	testProgram.Use()
	test.PerspectiveMatrix.Matrix4(perspectiveMatrix)

	cameraMatrix.Identity()
	cameraMatrix.Translate(float32(Camera.X), float32(Camera.Y+1.62), float32(-Camera.Z))
	cameraMatrix.RotateY(float32(Camera.Yaw))
	cameraMatrix.RotateX(float32(Camera.Pitch))
	cameraMatrix.Scale(-1.0, 1.0, 1.0)

	test.CameraMatrix.Matrix4(cameraMatrix)

	ids := make([]int, len(glTextures))
	for i, tex := range glTextures {
		tex.Bind(gl.Texture2D)
		gl.ActiveTexture(i)
		ids[i] = i
	}
	test.Textures.IntV(ids...)
	test.Position.Enable()
	test.TextureInfo.Enable()
	test.TextureOffset.Enable()

	for _, chunk := range buffers {
		if chunk.count == 0 {
			continue
		}
		test.Offset.Float3(float32(chunk.X), float32(chunk.Y), float32(chunk.Z))

		chunk.buffer.Bind(gl.ArrayBuffer)
		test.Position.Pointer(3, gl.Short, false, 18, 0)
		test.TextureInfo.Pointer(4, gl.UnsignedShort, false, 18, 6)
		test.TextureOffset.Pointer(2, gl.Short, false, 18, 14)
		gl.DrawArrays(gl.Triangles, 0, chunk.count)
	}

	test.TextureOffset.Disable()
	test.TextureInfo.Disable()
	test.Position.Disable()
}

func renderSync(f func()) {
	syncChan <- f
}

type testShader struct {
	Position          gl.Attribute `gl:"aPosition"`
	TextureInfo       gl.Attribute `gl:"aTextureInfo"`
	TextureOffset     gl.Attribute `gl:"aTextureOffset"`
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

uniform mat4 perspectiveMatrix;
uniform mat4 cameraMatrix;
uniform vec3 offset;

varying vec3 vPosition;
varying vec4 vTextureInfo;
varying vec2 vTextureOffset;

void main() {
	vec3 pos = vec3(aPosition.x, -aPosition.y, aPosition.z);
	vec3 o = vec3(offset.x, -offset.y, offset.z);
	gl_Position = perspectiveMatrix * cameraMatrix * vec4((pos / 256.0) + o * 16.0, 1.0);
	vPosition = aPosition / (256.0 * 16.0);
	vTextureInfo = aTextureInfo;
	vTextureOffset = aTextureOffset;
}
`
	fragment = `

uniform sampler2D textures[5];

varying vec3 vPosition;
varying vec4 vTextureInfo;
varying vec2 vTextureOffset;

void main() {
	vec2 tPos = vTextureOffset / 16.0;
	tPos = mod(tPos, vTextureInfo.zw);
	vec2 offset = vec2(vTextureInfo.x, mod(vTextureInfo.y, 1024.0));
	tPos += offset;
	tPos /= 1024.0;
 	vec4 col = texture2D(textures[int(floor(vTextureInfo.y / 1024.0))], vec2(tPos.x, tPos.y));
	if (col.a < 0.5) discard;
	gl_FragColor = col;
}
`
)
