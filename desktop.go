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

package main

import (
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/gl/v3.2-core/gl"
)

var window *glfw.Window

func startWindow() {
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

	start()

	for !window.ShouldClose() {
		draw()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

var lockMouse bool

func onMouseMove(w *glfw.Window, xpos float64, ypos float64) {
	if !lockMouse {
		return
	}
	width, height := w.GetFramebufferSize()
	ww, hh := float64(width/2), float64(height/2)
	w.SetCursorPos(ww, hh)

	rotate((xpos-ww)/2000.0, (ypos-hh)/2000.0)
}

func onMouseClick(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		lockMouse = true
		w.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	}
}

type Key int

const (
	KeyForward Key = iota
	KeyBackwards
	KeyLeft
	KeyRight
)

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch key {
	case glfw.KeyEscape:
		lockMouse = false
		w.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	case glfw.KeyW:
		if action != glfw.Repeat {
			Client.KeyState[KeyForward] = action == glfw.Press
		}
	case glfw.KeyS:
		if action != glfw.Repeat {
			Client.KeyState[KeyBackwards] = action == glfw.Press
		}
	case glfw.KeyA:
		if action != glfw.Repeat {
			Client.KeyState[KeyLeft] = action == glfw.Press
		}
	case glfw.KeyD:
		if action != glfw.Repeat {
			Client.KeyState[KeyRight] = action == glfw.Press
		}
	case glfw.KeySpace:
		if action != glfw.Repeat {
			Client.Jumping = action == glfw.Press
		}
	}
}
