package render

import (
	"math"

	"github.com/thinkofdeath/steven/platform"
	"github.com/thinkofdeath/steven/platform/gl"
	"github.com/thinkofdeath/steven/vmath"
)

var (
	testProgram gl.Program
	test        *testShader

	lastWidth, lastHeight int = -1, -1
	perspectiveMatrix         = vmath.NewMatrix4()
	cameraMatrix              = vmath.NewMatrix4()

	syncChan = make(chan func(), 500)
)

// Start starts the renderer
func Start() {
	gl.ClearColor(0.0, 1.0, 1.0, 1.0)
	gl.Enable(gl.DepthTest)
	gl.Enable(gl.CullFaceFlag)
	gl.CullFace(gl.Back)
	gl.FrontFace(gl.ClockWise)

	testProgram = CreateProgram(vertex, fragment)
	test = &testShader{}
	InitStruct(test, testProgram)

	loadTextures()
}

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

	cameraMatrix.Identity()
	cameraMatrix.Translate(float32(Camera.X), float32(Camera.Y+1.62), float32(-Camera.Z))
	cameraMatrix.RotateY(float32(Camera.Yaw))
	cameraMatrix.RotateX(float32(Camera.Pitch))
	cameraMatrix.Scale(-1.0, 1.0, 1.0)

	test.CameraMatrix.Matrix4(cameraMatrix)

	test.Position.Enable()
	test.Color.Enable()

	for _, chunk := range buffers {
		if chunk.count == 0 {
			continue
		}
		test.Offset.Float3(float32(chunk.X), float32(chunk.Y), float32(chunk.Z))

		chunk.buffer.Bind(gl.ArrayBuffer)
		test.Position.Pointer(3, gl.Short, false, 9, 0)
		test.Color.Pointer(3, gl.UnsignedByte, true, 9, 6)
		gl.DrawArrays(gl.Triangles, 0, chunk.count)
	}

	test.Color.Disable()
	test.Position.Disable()
}

func renderSync(f func()) {
	syncChan <- f
}

type testShader struct {
	Position          gl.Attribute `gl:"aPosition"`
	Color             gl.Attribute `gl:"aColor"`
	PerspectiveMatrix gl.Uniform   `gl:"perspectiveMatrix"`
	CameraMatrix      gl.Uniform   `gl:"cameraMatrix"`
	Offset            gl.Uniform   `gl:"offset"`
}

var (
	vertex = `
attribute vec3 aPosition;
attribute vec3 aColor;

uniform mat4 perspectiveMatrix;
uniform mat4 cameraMatrix;
uniform vec3 offset;

varying vec3 vColor;

void main() {
	vec3 pos = vec3(aPosition.x, -aPosition.y, aPosition.z);
	vec3 o = vec3(offset.x, -offset.y, offset.z);
	gl_Position = perspectiveMatrix * cameraMatrix * vec4((pos / 256.0) + o * 16.0, 1.0);
	vColor = aColor;
}
`
	fragment = `

varying vec3 vColor;

void main() {
 	gl_FragColor = vec4(vColor, 1.0);
}
`
)
