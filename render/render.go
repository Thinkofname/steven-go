package render

import (
	"math"

	"github.com/thinkofdeath/steven/platform"
	"github.com/thinkofdeath/steven/platform/gl"
	"github.com/thinkofdeath/steven/type/direction"
	"github.com/thinkofdeath/steven/type/vmath"
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

var (
	textureIds    []int
	frameID       uint = 0
	nearestBuffer *ChunkBuffer
	viewVector    vmath.Vector3
)

// Draw draws a single frame
func Draw() {
	frameID++
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
	// Only update the viewport if the window was resized
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
	// Only update the texture ids if we have new
	// textures
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
	// +1.62 for the players height.
	// TODO(Think) Change this?
	cameraMatrix.Translate(float32(Camera.X), float32(Camera.Y+1.62), float32(-Camera.Z))
	cameraMatrix.RotateY(float32(Camera.Yaw))
	cameraMatrix.RotateX(float32(Camera.Pitch))
	cameraMatrix.Scale(-1.0, 1.0, 1.0)

	shaderChunk.CameraMatrix.Matrix4(cameraMatrix)
	shaderChunk.Textures.IntV(textureIds...)

	nearestBuffer = nil
	distance := math.MaxFloat64
	for _, chunk := range buffers {
		dx := Camera.X - float64((chunk.X<<4)+8)
		dy := Camera.Y - float64((chunk.Y<<4)+8)
		dz := Camera.Z - float64((chunk.Z<<4)+8)
		dist := dx*dx + dy*dy + dz*dz
		if nearestBuffer == nil || dist < distance {
			nearestBuffer = chunk
			distance = dist
		}
	}

	viewVector.X = float32(math.Cos(float64(Camera.Yaw-math.Pi/2)) * -math.Cos(float64(Camera.Pitch)))
	viewVector.Z = -float32(math.Sin(float64(Camera.Yaw-math.Pi/2)) * -math.Cos(float64(Camera.Pitch)))
	viewVector.Y = -float32(math.Sin(float64(Camera.Pitch)))

	colVisitMap = make(map[positionC]struct{})
	if nearestBuffer != nil {
		renderBuffer(nearestBuffer, nearestBuffer.position, direction.Invalid)
	}
}

var colVisitMap = make(map[positionC]struct{})

func renderBuffer(chunk *ChunkBuffer, pos position, from direction.Type) {
	v := vmath.Vector3{
		float32((pos.X<<4)+8) - float32(Camera.X),
		float32((pos.Y<<4)+8) - float32(Camera.Y),
		float32((pos.Z<<4)+8) - float32(Camera.Z),
	}
	if v.LengthSquared() > 40*40 && v.Dot(viewVector) < 0 {
		return
	}
	if chunk == nil {
		// Handle empty sections in columns
		if pos.Y >= 0 && pos.Y <= 16 {
			col := positionC{pos.X, pos.Z}
			if _, ok := colVisitMap[col]; !ok && bufferColumns[col] > 0 {
				colVisitMap[col] = struct{}{}
				for _, dir := range direction.Values {
					if dir != from {
						ox, oy, oz := dir.Offset()
						pos := position{pos.X + ox, pos.Y + oy, pos.Z + oz}
						renderBuffer(buffers[pos], pos, dir.Opposite())
					}
				}

			}
		}
		return
	}
	if chunk.renderedOn == frameID {
		return
	}
	chunk.renderedOn = frameID

	if chunk.count > 0 {
		shaderChunk.Offset.Float3(float32(chunk.X), float32(chunk.Y), float32(chunk.Z))

		chunk.array.Bind()
		gl.DrawArrays(gl.Triangles, 0, chunk.count)
	}

	for _, dir := range direction.Values {
		if dir != from && (from == direction.Invalid || chunk.IsVisible(from, dir)) {
			ox, oy, oz := dir.Offset()
			pos := position{pos.X + ox, pos.Y + oy, pos.Z + oz}
			renderBuffer(buffers[pos], pos, dir.Opposite())
		}
	}
}
func renderSync(f func()) {
	syncChan <- f
}
