package render

import (
	"bytes"
	"encoding/binary"
	"github.com/thinkofdeath/steven/platform"
	"github.com/thinkofdeath/steven/platform/gl"
	"github.com/thinkofdeath/steven/vmath"
	"math"
)

var (
	testProgram gl.Program
	test        *testShader
	testBuffer  gl.Buffer

	lastWidth, lastHeight int = -1, -1
	perspectiveMatrix         = vmath.NewMatrix4()
)

func Start() {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Enable(gl.DepthTest)
	gl.Enable(gl.CullFaceFlag)
	gl.CullFace(gl.Back)
	gl.FrontFace(gl.CounterClockWise)

	testProgram = CreateProgram(vertex, fragment)
	test = &testShader{}
	InitStruct(test, testProgram)

	testBuffer = gl.CreateBuffer()
	testBuffer.Bind(gl.ArrayBuffer)
	var buf bytes.Buffer
	binary.Write(&buf, platform.NativeOrder, []float32{
		0.0, 1.0, 0.0,
		-1.0, -1.0, 0.0,
		1.0, -1.0, 0.0,
	})
	testBuffer.Data(buf.Bytes(), gl.StaticDraw)
}

var (
	offset float32
	dir    float32 = 0.1
)

func Draw() {
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
	test.Offset.Float(offset)
	offset += dir
	if offset > 20.0 {
		dir = -0.1
	} else if offset < 1.0 {
		dir = 0.1
	}

	test.Position.Enable()
	testBuffer.Bind(gl.ArrayBuffer)
	test.Position.Pointer(3, gl.Float, false, 12, 0)
	gl.DrawArrays(gl.Triangles, 0, 3)

	test.Position.Disable()
}

type testShader struct {
	Position          gl.Attribute `gl:"aPosition"`
	PerspectiveMatrix gl.Uniform   `gl:"perspectiveMatrix"`
	Offset            gl.Uniform   `gl:"offset"`
}

var (
	vertex = `

attribute vec3 aPosition;

uniform mat4 perspectiveMatrix;
uniform float offset;

void main() {
	gl_Position = perspectiveMatrix * vec4(aPosition - vec3(0.0, 0.0, offset), 1.0);
}
`
	fragment = `
void main() {
 	gl_FragColor = vec4(0.0, 1.0, 0.0, 1.0);
}
`
)
