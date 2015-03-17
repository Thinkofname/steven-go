// +build !mobile

package platform

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"runtime"
)

func run(handler Handler) {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMajor, 1)

	window, err := glfw.CreateWindow(800, 480, "Steven", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	handler.Start()

	for !window.ShouldClose() {
		handler.Draw()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
