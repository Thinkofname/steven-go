// +build !mobile

package platform

import (
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

var window *glfw.Window
var handler Handler

func run(h Handler) {
	handler = h
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMajor, 1)

	var err error
	window, err = glfw.CreateWindow(800, 480, "Steven", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	window.SetCursorPosCallback(onMouseMove)
	window.SetMouseButtonCallback(onMouseClick)
	window.SetKeyCallback(onKey)

	handler.Start()

	for !window.ShouldClose() {
		handler.Draw()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func size() (int, int) {
	return window.GetFramebufferSize()
}

var lockMouse bool

func onMouseMove(w *glfw.Window, xpos float64, ypos float64) {
	if handler.Rotate == nil || !lockMouse {
		return
	}
	width, height := size()
	ww, hh := float64(width/2), float64(height/2)
	w.SetCursorPos(ww, hh)

	handler.Rotate((xpos-ww)/2000.0, (ypos-hh)/2000.0)
}

func onMouseClick(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		lockMouse = true
		w.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	}
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch key {
	case glfw.KeyEscape:
		lockMouse = false
		w.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	case glfw.KeyW:
		if handler.Move != nil {
			if action == glfw.Press {
				handler.Move(1, 0)
			} else if action == glfw.Release {
				handler.Move(0, 0)
			}
		}
	}
}
