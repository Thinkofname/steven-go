package render

import "github.com/thinkofdeath/steven/platform/gl"

var buffers []*ChunkBuffer

type ChunkBuffer struct {
	X, Y, Z int

	buffer       gl.Buffer
	count        int
	lastRenderID int
}

// AllocateChunkBuffer allocates a chunk buffer and adds it to the
// render list.
func AllocateChunkBuffer(x, y, z int) *ChunkBuffer {
	c := &ChunkBuffer{
		X: x, Y: y, Z: z,
		buffer: gl.CreateBuffer(),
	}
	buffers = append(buffers, c)
	return c
}

func (cb *ChunkBuffer) Upload(data []byte, count, rid int) {
	sync(func() {
		if rid < cb.lastRenderID {
			return
		}
		cb.lastRenderID = rid
		cb.buffer.Bind(gl.ArrayBuffer)
		cb.buffer.Data(data, gl.StaticDraw)
		cb.count = count
	})
}

func (cb *ChunkBuffer) Free() {
	for i, c := range buffers {
		if c == cb {
			buffers = append(buffers[:i], buffers[i+1:]...)
			return
		}
	}
	cb.buffer.Delete()
}
