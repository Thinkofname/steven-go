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
	"github.com/thinkofdeath/steven/render/gl"
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

	glTexture    gl.Texture
	textureDepth int
)

// Start starts the renderer
func Start(debug bool) {
	if debug {
		gl.Enable(gl.DebugOutput)
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

	textureLock.Lock()
	glTexture = gl.CreateTexture()
	glTexture.Bind(gl.Texture2DArray)
	textureDepth = len(textures)
	glTexture.Image3D(0, AtlasSize, AtlasSize, len(textures), gl.RGBA, gl.UnsignedByte, make([]byte, AtlasSize*AtlasSize*len(textures)*4))
	glTexture.Parameter(gl.TextureMagFilter, gl.Nearest)
	glTexture.Parameter(gl.TextureMinFilter, gl.Linear)
	glTexture.Parameter(gl.TextureWrapS, gl.ClampToEdge)
	glTexture.Parameter(gl.TextureWrapT, gl.ClampToEdge)
	for i, tex := range textures {
		glTexture.SubImage3D(0, 0, 0, i, AtlasSize, AtlasSize, 1, gl.RGBA, gl.UnsignedByte, tex.Buffer)
	}
	textureLock.Unlock()

	initUI()

	gl.BlendFunc(gl.SrcAlpha, gl.OneMinusSrcAlpha)
}

var (
	textureIds    []int
	frameID       uint
	nearestBuffer *ChunkBuffer
	viewVector    vmath.Vector3
)

// Draw draws a single frame
func Draw(delta float64) {
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

	// Textures
	textureLock.RLock()
	if textureDepth != len(textures) {
		glTexture.Bind(gl.Texture2DArray)
		textureDepth = len(textures)
		glTexture.Image3D(0, AtlasSize, AtlasSize, len(textures), gl.RGBA, gl.UnsignedByte, make([]byte, AtlasSize*AtlasSize*len(textures)*4))
		for i := range textureDirty {
			textureDirty[i] = true
		}
	}
	for i, tex := range textures {
		if textureDirty[i] {
			textureDirty[i] = true
			glTexture.SubImage3D(0, 0, 0, i, AtlasSize, AtlasSize, 1, gl.RGBA, gl.UnsignedByte, tex.Buffer)
		}
	}
	textureLock.RUnlock()

	glTexture.Bind(gl.Texture2DArray)
	gl.ActiveTexture(0)

	gl.Clear(gl.ColorBufferBit | gl.DepthBufferBit)

	chunkProgram.Use()

	cameraMatrix.Identity()
	// +1.62 for the players height.
	// TODO(Think) Change this?
	cameraMatrix.Translate(float32(Camera.X), float32(Camera.Y), float32(-Camera.Z))
	cameraMatrix.RotateY(float32(Camera.Yaw))
	cameraMatrix.RotateX(float32(Camera.Pitch))
	cameraMatrix.Scale(-1.0, 1.0, 1.0)

	shaderChunk.PerspectiveMatrix.Matrix4(perspectiveMatrix)
	shaderChunk.CameraMatrix.Matrix4(cameraMatrix)
	shaderChunk.Texture.Int(0)

	chunkPos := position{
		X: int(Camera.X) >> 4,
		Y: int(Camera.Y) >> 4,
		Z: int(Camera.Z) >> 4,
	}
	nearestBuffer = buffers[chunkPos]

	viewVector.X = math.Cos(Camera.Yaw-math.Pi/2) * -math.Cos(Camera.Pitch)
	viewVector.Z = -math.Sin(Camera.Yaw-math.Pi/2) * -math.Cos(Camera.Pitch)
	viewVector.Y = -math.Sin(Camera.Pitch)

	for _, dir := range direction.Values {
		validDirs[dir] = viewVector.Dot(dir.AsVector()) > -0.8
	}

	renderOrder = renderOrder[:0]
	renderBuffer(nearestBuffer, chunkPos, direction.Invalid)

	chunkProgramT.Use()
	shaderChunkT.PerspectiveMatrix.Matrix4(perspectiveMatrix)
	shaderChunkT.CameraMatrix.Matrix4(cameraMatrix)
	shaderChunkT.Texture.Int(0)

	gl.Enable(gl.Blend)
	for i := range renderOrder {
		chunk := renderOrder[len(renderOrder)-1-i]
		if chunk.countT > 0 && chunk.bufferT.IsValid() {
			shaderChunkT.Offset.Int3(chunk.X, chunk.Y, chunk.Z)

			chunk.arrayT.Bind()
			chunk.bufferT.Bind(gl.ArrayBuffer)
			data := chunk.transBuffer
			offset := 0
			sort.Sort(chunk.transInfo)
			for _, i := range chunk.transInfo {
				offset += copy(data[offset:], chunk.transData[i.Offset:i.Offset+i.Count])
			}
			chunk.bufferT.SubData(0, data)
			gl.DrawArrays(gl.Triangles, 0, chunk.countT)
		}
	}
	gl.Disable(gl.Blend)

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
	dist  int
}

const (
	renderQueueSize = 5000
)

var rQueue renderQueue

func renderBuffer(chunk *ChunkBuffer, pos position, from direction.Type) {
	rQueue.Append(renderRequest{chunk, pos, from, 1})
itQueue:
	for !rQueue.Empty() {
		req := rQueue.Take()
		chunk, pos, from = req.chunk, req.pos, req.from
		v := vmath.Vector3{
			float64((pos.X<<4)+8) - Camera.X,
			float64((pos.Y<<4)+8) - Camera.Y,
			float64((pos.Z<<4)+8) - Camera.Z,
		}
		if (v.LengthSquared() > 40*40 && v.Dot(viewVector) < 0) || req.dist > 20 {
			continue itQueue
		}
		if chunk == nil || chunk.renderedOn == frameID {
			continue itQueue
		}
		chunk.renderedOn = frameID
		renderOrder = append(renderOrder, chunk)

		if chunk.count > 0 && chunk.buffer.IsValid() {
			shaderChunk.Offset.Int3(chunk.X, chunk.Y, chunk.Z)

			chunk.array.Bind()
			gl.DrawArrays(gl.Triangles, 0, chunk.count)
		}

		for _, dir := range direction.Values {
			if dir != from && (from == direction.Invalid || (chunk.IsVisible(from, dir) && validDirs[dir])) {
				ox, oy, oz := dir.Offset()
				pos := position{pos.X + ox, pos.Y + oy, pos.Z + oz}
				rQueue.Append(renderRequest{chunk.neighborChunks[dir], pos, dir.Opposite(), req.dist + 1})
			}
		}
	}
}

// Sync runs the passed function on the next frame on the same goroutine
// as the renderer.
func Sync(f func()) {
	syncChan <- f
}
