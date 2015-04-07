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

package platform

import (
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/gl/v3.2-core/gl"
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

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var err error
	window, err = glfw.CreateWindow(800, 480, "Steven", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	window.SetCursorPosCallback(onMouseMove)
	window.SetMouseButtonCallback(onMouseClick)
	window.SetKeyCallback(onKey)

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
		if action == glfw.Press {
			handler.Move(1, 0)
		} else if action == glfw.Release {
			handler.Move(0, 0)
		}
	case glfw.KeyH:
		if action == glfw.Release {
			handler.Action(Debug)
		}
	case glfw.KeySpace:
		if action == glfw.Release || action == glfw.Press {
			handler.Action(JumpToggle)
		}
	}
}
