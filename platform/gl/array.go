package gl

import (
	"github.com/go-gl/gl/v3.2-core/gl"
)

type VertexArray struct {
	internal uint32
}

var currentVertexArray VertexArray

func CreateVertexArray() VertexArray {
	var va VertexArray
	gl.GenVertexArrays(1, &va.internal)
	return va
}

func (va VertexArray) Bind() {
	if currentVertexArray == va {
		return
	}
	gl.BindVertexArray(va.internal)
	currentVertexArray = va
}

func (va VertexArray) Delete() {
	gl.DeleteVertexArrays(1, &va.internal)
	if currentVertexArray == va {
		currentVertexArray = VertexArray{}
	}
}
