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
	"sort"
	"strings"

	"github.com/thinkofdeath/steven/platform"
	"github.com/thinkofdeath/steven/platform/gl"
	"github.com/thinkofdeath/steven/type/direction"
	"github.com/thinkofdeath/steven/type/vmath"
)

var (
	chunkProgram  gl.Program
	shaderChunk   *chunkShader
	chunkProgramT gl.Program
	shaderChunkT  *chunkShader

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

	chunkProgramT = CreateProgram(vertex, strings.Replace(fragment, "#version 150", "#version 150\n#define alpha", 1))
	shaderChunkT = &chunkShader{}
	InitStruct(shaderChunkT, chunkProgramT)

	textureLock.Lock()
	for _, tex := range textures {
		glTextures = append(glTextures, createTexture(glTexture{
			Data:  tex.Buffer,
			Width: atlasSize, Height: atlasSize,
			Format: gl.RGBA,
		}))
	}
	textureLock.Unlock()

	gl.BlendFunc(gl.SrcAlpha, gl.OneMinusSrcAlpha)
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

	cameraMatrix.Identity()
	// +1.62 for the players height.
	// TODO(Think) Change this?
	cameraMatrix.Translate(float32(Camera.X), float32(Camera.Y+1.62), float32(-Camera.Z))
	cameraMatrix.RotateY(float32(Camera.Yaw))
	cameraMatrix.RotateX(float32(Camera.Pitch))
	cameraMatrix.Scale(-1.0, 1.0, 1.0)

	shaderChunk.PerspectiveMatrix.Matrix4(perspectiveMatrix)
	shaderChunk.CameraMatrix.Matrix4(cameraMatrix)
	shaderChunk.Textures.IntV(textureIds...)

	nearestBuffer = buffers[position{
		X: int(Camera.X) >> 4,
		Y: int(Camera.Y) >> 4,
		Z: int(Camera.Z) >> 4,
	}]
	if nearestBuffer == nil {
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
	}

	viewVector.X = float32(math.Cos(float64(Camera.Yaw-math.Pi/2)) * -math.Cos(float64(Camera.Pitch)))
	viewVector.Z = -float32(math.Sin(float64(Camera.Yaw-math.Pi/2)) * -math.Cos(float64(Camera.Pitch)))
	viewVector.Y = -float32(math.Sin(float64(Camera.Pitch)))

	airVisitMap = make(map[position]struct{})
	renderOrder = renderOrder[:0]
	if nearestBuffer != nil {
		renderBuffer(nearestBuffer, nearestBuffer.position, direction.Invalid)
	}

	chunkProgramT.Use()
	shaderChunkT.PerspectiveMatrix.Matrix4(perspectiveMatrix)
	shaderChunkT.CameraMatrix.Matrix4(cameraMatrix)
	shaderChunkT.Textures.IntV(textureIds...)
	sort.Sort(renderOrder)

	gl.Enable(gl.Blend)
	for _, pos := range renderOrder {
		chunk := buffers[pos]
		if chunk != nil && chunk.countT > 0 {
			shaderChunkT.Offset.Float3(float32(chunk.X), float32(chunk.Y), float32(chunk.Z))

			chunk.arrayT.Bind()
			gl.DrawArrays(gl.Triangles, 0, chunk.countT)
		}
	}
	gl.Disable(gl.Blend)
}

var (
	airVisitMap = make(map[position]struct{})
	renderOrder transList
)

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
		if pos.Y >= 0 && pos.Y <= 15 {
			col := positionC{pos.X, pos.Z}
			if _, ok := airVisitMap[pos]; !ok && bufferColumns[col] > 0 {
				airVisitMap[pos] = struct{}{}
				renderOrder = append(renderOrder, pos)
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
	renderOrder = append(renderOrder, pos)

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

// Sync runs the passed function on the next frame on the same goroutine
// as the renderer.
func Sync(f func()) {
	syncChan <- f
}
