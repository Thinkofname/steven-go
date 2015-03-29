package render

import "github.com/thinkofdeath/steven/platform/gl"

var buffers []*ChunkBuffer

// ChunkBuffer is a renderable chunk section
type ChunkBuffer struct {
	X, Y, Z int

	array  gl.VertexArray
	buffer gl.Buffer
	count  int
}

// AllocateChunkBuffer allocates a chunk buffer and adds it to the
// render list.
func AllocateChunkBuffer(x, y, z int) *ChunkBuffer {
	c := &ChunkBuffer{
		X: x, Y: y, Z: z,
		array:  gl.CreateVertexArray(),
		buffer: gl.CreateBuffer(),
	}
	buffers = append(buffers, c)
	return c
}

// Upload uploads the passed vertex data to the buffer.
func (cb *ChunkBuffer) Upload(data []byte, count int) {
	renderSync(func() {
		cb.array.Bind()
		cb.buffer.Bind(gl.ArrayBuffer)
		cb.buffer.Data(data, gl.DynamicDraw)
		shaderChunk.Position.Enable()
		shaderChunk.TextureInfo.Enable()
		shaderChunk.TextureOffset.Enable()
		shaderChunk.Color.Enable()
		shaderChunk.Lighting.Enable()

		cb.buffer.BindVertex(0, 0, 23)

		shaderChunk.Position.Format(3, gl.Short, false, 0)
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
	})
}

// Free removes the buffer and frees related resources.
func (cb *ChunkBuffer) Free() {
	for i, c := range buffers {
		if c == cb {
			buffers = append(buffers[:i], buffers[i+1:]...)
			return
		}
	}
	cb.buffer.Delete()
	cb.array.Delete()
}
