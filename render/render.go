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
