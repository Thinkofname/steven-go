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

import (
	"math"
	"os"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/thinkofdeath/steven/render/gl"
	"github.com/thinkofdeath/steven/type/direction"
	"github.com/thinkofdeath/steven/type/vmath"
)

var (
	chunkProgram  gl.Program
	shaderChunk   *chunkShader
	chunkProgramT gl.Program
	shaderChunkT  *chunkShader
	lineProgram   gl.Program
	shaderLine    *lineShader

	FOV, lastFOV          int = 90, 90
	lastWidth, lastHeight int = -1, -1
	perspectiveMatrix         = mgl32.Mat4{}
	cameraMatrix              = mgl32.Mat4{}
	frustum                   = vmath.NewFrustum()

	syncChan = make(chan func(), 500)

	glTexture       gl.Texture
	textureDepth    int
	texturesCreated bool

	MultiSample bool
)

// Start starts the renderer
func Start() {
	if os.Getenv("STEVEN_DEBUG") == "true" {
		gl.DebugLog()
	}

	gl.ClearColor(122.0/255.0, 165.0/255.0, 247.0/255.0, 1.0)
	gl.Enable(gl.DepthTest)
	gl.Enable(gl.CullFaceFlag)
	gl.CullFace(gl.Back)
	gl.FrontFace(gl.ClockWise)

	chunkProgram = CreateProgram(vertex, fragment)
	shaderChunk = &chunkShader{}
	InitStruct(shaderChunk, chunkProgram)

	chunkProgramT = CreateProgram(vertex, strings.Replace(fragment, "#version 150", "#version 150\n#define alpha", 1))
	shaderChunkT = &chunkShader{}
	InitStruct(shaderChunkT, chunkProgramT)

	initUI()
	initLineDraw()
	initStatic()

	gl.BlendFunc(gl.SrcAlpha, gl.OneMinusSrcAlpha)

	elementBuffer = gl.CreateBuffer()
}

var (
	textureIds    []int
	frameID       uint
	nearestBuffer *ChunkBuffer
	viewVector    mgl32.Vec3
)

// Draw draws a single frame
func Draw(width, height int, delta float64) {
	tickAnimatedTextures(delta)
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

	// Only update the viewport if the window was resized
	if lastHeight != height || lastWidth != width || lastFOV != FOV {
		lastWidth = width
		lastHeight = height
		lastFOV = FOV

		perspectiveMatrix = mgl32.Perspective(
			(math.Pi/180)*float32(FOV),
			float32(width)/float32(height),
			1,
			500.0,
		)
		gl.Viewport(0, 0, width, height)
		frustum.SetPerspective(
			(math.Pi/180)*float32(FOV),
			float32(width)/float32(height),
			1,
			500.0,
		)
	}

	if MultiSample {
		gl.Enable(gl.Multisample)
	}

	glTexture.Bind(gl.Texture2DArray)
	gl.ActiveTexture(0)

	gl.Clear(gl.ColorBufferBit | gl.DepthBufferBit)

	chunkProgram.Use()

	viewVector = mgl32.Vec3{
		float32(math.Cos(Camera.Yaw-math.Pi/2) * -math.Cos(Camera.Pitch)),
		float32(-math.Sin(Camera.Pitch)),
		float32(-math.Sin(Camera.Yaw-math.Pi/2) * -math.Cos(Camera.Pitch)),
	}
	cam := mgl32.Vec3{-float32(Camera.X), -float32(Camera.Y), float32(Camera.Z)}
	cameraMatrix = mgl32.LookAtV(
		cam,
		cam.Add(mgl32.Vec3{-viewVector.X(), -viewVector.Y(), viewVector.Z()}),
		mgl32.Vec3{0, -1, 0},
	)
	cameraMatrix = cameraMatrix.Mul4(mgl32.Scale3D(-1.0, 1.0, 1.0))

	frustum.SetCamera(
		cam,
		cam.Add(mgl32.Vec3{-viewVector.X(), -viewVector.Y(), viewVector.Z()}),
		mgl32.Vec3{0, -1, 0},
	)

	shaderChunk.PerspectiveMatrix.Matrix4(&perspectiveMatrix)
	shaderChunk.CameraMatrix.Matrix4(&cameraMatrix)
	shaderChunk.Texture.Int(0)

	chunkPos := position{
		X: int(Camera.X) >> 4,
		Y: int(Camera.Y) >> 4,
		Z: int(Camera.Z) >> 4,
	}
	nearestBuffer = buffers[chunkPos]

	for _, dir := range direction.Values {
		validDirs[dir] = viewVector.Dot(dir.AsVec()) > -0.8
	}

	renderOrder = renderOrder[:0]
	renderBuffer(nearestBuffer, chunkPos, direction.Invalid)

	drawLines()
	drawStatic()

	chunkProgramT.Use()
	shaderChunkT.PerspectiveMatrix.Matrix4(&perspectiveMatrix)
	shaderChunkT.CameraMatrix.Matrix4(&cameraMatrix)
	shaderChunkT.Texture.Int(0)

	gl.Enable(gl.Blend)
	for i := range renderOrder {
		chunk := renderOrder[len(renderOrder)-1-i]
		if chunk.countT > 0 && chunk.bufferT.IsValid() {
			shaderChunkT.Offset.Int3(chunk.X, chunk.Y, chunk.Z)

			chunk.arrayT.Bind()
			chunk.bufferTI.Bind(gl.ElementArrayBuffer)
			insertSort(chunk.transInfo)

			offset := 0
			data := chunk.bufferTI.Map(gl.WriteOnly, len(chunk.transData))
			m := 2
			if chunk.bufferTIType == gl.UnsignedInt {
				m = 4
			}
			for _, i := range chunk.transInfo {
				offset += copy(data[offset:], chunk.transData[i.Offset*m:(i.Offset+i.Count)*m])
			}
			chunk.bufferTI.Unmap()
			gl.DrawElements(gl.Triangles, chunk.countT, chunk.bufferTIType, 0)
		}
	}
	gl.Disable(gl.Blend)

	if MultiSample {
		gl.Disable(gl.Multisample)
	}

	drawUI()
}

var (
	renderOrder []*ChunkBuffer
	validDirs   = make([]bool, len(direction.Values))
)

type renderRequest struct {
	chunk *ChunkBuffer
	pos   position
	from  direction.Type
}

const (
	renderQueueSize = 5000
)

var rQueue renderQueue

func renderBuffer(ch *ChunkBuffer, po position, fr direction.Type) {
	if ch == nil {
		return
	}
	rQueue.Append(renderRequest{ch, po, fr})
itQueue:
	for !rQueue.Empty() {
		req := rQueue.Take()
		if req.chunk.renderedOn == frameID {
			continue itQueue
		}
		aabb := vmath.NewAABB(
			-float32((req.pos.X<<4)+16), -float32((req.pos.Y<<4)+16), float32((req.pos.Z<<4)),
			-float32((req.pos.X<<4)), -float32((req.pos.Y<<4)), float32((req.pos.Z<<4)+16),
		).Grow(1, 1, 1)
		if !frustum.IsAABBInside(aabb) {
			req.chunk.renderedOn = frameID
			continue itQueue
		}
		req.chunk.renderedOn = frameID
		renderOrder = append(renderOrder, req.chunk)

		if req.chunk.count != 0 && req.chunk.buffer.IsValid() {
			shaderChunk.Offset.Int3(req.chunk.X, req.chunk.Y, req.chunk.Z)

			req.chunk.array.Bind()
			gl.DrawElements(gl.Triangles, req.chunk.count, elementBufferType, 0)
		}

		for _, dir := range direction.Values {
			c := req.chunk.neighborChunks[dir]
			if dir != req.from && c != nil && c.renderedOn != frameID &&
				(req.from == direction.Invalid || (req.chunk.IsVisible(req.from, dir) && validDirs[dir])) {
				ox, oy, oz := dir.Offset()
				pos := position{req.pos.X + ox, req.pos.Y + oy, req.pos.Z + oz}
				rQueue.Append(renderRequest{c, pos, dir.Opposite()})
			}
		}
	}
}

// Sync runs the passed function on the next frame on the same goroutine
// as the renderer.
func Sync(f func()) {
	syncChan <- f
}
