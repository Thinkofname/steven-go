package render

import (
	"github.com/thinkofdeath/steven/platform/gl"
	"github.com/thinkofdeath/steven/type/direction"
)

var (
	buffers       = make(map[position]*ChunkBuffer)
	bufferColumns = make(map[positionC]int)
)

// ChunkBuffer is a renderable chunk section
type ChunkBuffer struct {
	position

	array    gl.VertexArray
	buffer   gl.Buffer
	count    int
	arrayT   gl.VertexArray
	bufferT  gl.Buffer
	countT   int
	cullBits uint64

	renderedOn uint
}

func (cb *ChunkBuffer) IsVisible(from, to direction.Type) bool {
	return (cb.cullBits & (1 << (from*6 + to))) != 0
}

// AllocateChunkBuffer allocates a chunk buffer and adds it to the
// render list.
func AllocateChunkBuffer(x, y, z int) *ChunkBuffer {
	c := &ChunkBuffer{
		position: position{X: x, Y: y, Z: z},
		array:    gl.CreateVertexArray(),
		buffer:   gl.CreateBuffer(),
		arrayT:   gl.CreateVertexArray(),
		bufferT:  gl.CreateBuffer(),
	}
	buffers[c.position] = c
	bufferColumns[positionC{x, z}]++
	return c
}

// Upload uploads the passed vertex data to the buffer.
func (cb *ChunkBuffer) Upload(data []byte, count int, cullBits uint64) {
	chunkProgram.Use()
	cb.cullBits = cullBits
	cb.array.Bind()
	cb.buffer.Bind(gl.ArrayBuffer)
	cb.buffer.Data(data, gl.DynamicDraw)
	shaderChunk.Position.Enable()
	shaderChunk.TextureInfo.Enable()
	shaderChunk.TextureOffset.Enable()
	shaderChunk.Color.Enable()
	shaderChunk.Lighting.Enable()

	cb.buffer.BindVertex(0, 0, 23)

	shaderChunk.Position.FormatInt(3, gl.Short, 0)
	shaderChunk.Position.Binding(0)

	shaderChunk.TextureInfo.Format(4, gl.UnsignedShort, false, 6)
	shaderChunk.TextureInfo.Binding(0)

	shaderChunk.TextureOffset.Format(2, gl.Short, false, 14)
	shaderChunk.TextureOffset.Binding(0)

	shaderChunk.Color.Format(3, gl.UnsignedByte, true, 18)
	shaderChunk.Color.Binding(0)

	shaderChunk.Lighting.Format(2, gl.UnsignedByte, false, 21)
	shaderChunk.Lighting.Binding(0)

	cb.count = count
}

// UploadTrans uploads the passed vertex data to the translucent buffer.
func (cb *ChunkBuffer) UploadTrans(data []byte, count int) {
	chunkProgramT.Use()
	cb.arrayT.Bind()
	cb.bufferT.Bind(gl.ArrayBuffer)
	cb.bufferT.Data(data, gl.DynamicDraw)
	shaderChunk.Position.Enable()
	shaderChunk.TextureInfo.Enable()
	shaderChunk.TextureOffset.Enable()
	shaderChunk.Color.Enable()
	shaderChunk.Lighting.Enable()

	cb.bufferT.BindVertex(0, 0, 23)

	shaderChunkT.Position.FormatInt(3, gl.Short, 0)
	shaderChunkT.Position.Binding(0)

	shaderChunkT.TextureInfo.Format(4, gl.UnsignedShort, false, 6)
	shaderChunkT.TextureInfo.Binding(0)

	shaderChunkT.TextureOffset.Format(2, gl.Short, false, 14)
	shaderChunkT.TextureOffset.Binding(0)

	shaderChunkT.Color.Format(3, gl.UnsignedByte, true, 18)
	shaderChunkT.Color.Binding(0)

	shaderChunkT.Lighting.Format(2, gl.UnsignedByte, false, 21)
	shaderChunkT.Lighting.Binding(0)

	cb.countT = count
}

// Free removes the buffer and frees related resources.
func (cb *ChunkBuffer) Free() {
	delete(buffers, cb.position)
	cpos := positionC{cb.position.X, cb.position.Z}
	val := bufferColumns[cpos]
	val--
	if val <= 0 {
		delete(bufferColumns, cpos)
	} else {
		bufferColumns[cpos] = val
	}

	cb.buffer.Delete()
	cb.array.Delete()
}
