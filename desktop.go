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

package steven

import (
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/render/gl"
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
	glfw.WindowHint(glfw.Samples, Config.Render.Samples)
	render.MultiSample = Config.Render.Samples > 0

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
	window.SetScrollCallback(onScroll)

	gl.Init()

	start()

	for !window.ShouldClose() {
		draw()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func onScroll(w *glfw.Window, xoff float64, yoff float64) {
	if currentScreen != nil || !ready {
		return
	}
	if yoff < 0 {
		Client.currentHotbarSlot++
	} else {
		Client.currentHotbarSlot--
	}
	if Client.currentHotbarSlot < 0 {
		Client.currentHotbarSlot = 0
	} else if Client.currentHotbarSlot > 8 {
		Client.currentHotbarSlot = 8
	}

	writeChan <- &protocol.HeldItemChange{Slot: int16(Client.currentHotbarSlot)}
}

var lockMouse bool

func onMouseMove(w *glfw.Window, xpos float64, ypos float64) {
	width, height := w.GetFramebufferSize()
	if currentScreen != nil {
		currentScreen.hover(xpos, ypos, width, height)
		return
	}
	if !lockMouse {
		return
	}
	ww, hh := float64(width/2), float64(height/2)
	w.SetCursorPos(ww, hh)

	rotate((xpos-ww)/2000.0, (ypos-hh)/2000.0)
}

func onMouseClick(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if currentScreen != nil {
		if button != glfw.MouseButtonLeft || action == glfw.Repeat {
			return
		}
		width, height := w.GetFramebufferSize()
		xpos, ypos := w.GetCursorPos()
		currentScreen.click(action == glfw.Press, xpos, ypos, width, height)
		return
	}
	if button == glfw.MouseButtonLeft && action == glfw.Press && !Client.chat.enteringText {
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
	KeySprint
)

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if currentScreen != nil {
		return
	}
	if Client.chat.enteringText {
		Client.chat.handleKey(w, key, scancode, action, mods)
		return
	}
	switch key {
	case glfw.KeyEscape:
		if action == glfw.Release {
			setScreen(newGameMenu())
		}
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
	case glfw.KeyLeftControl:
		if action != glfw.Repeat {
			Client.KeyState[KeySprint] = action == glfw.Press
		}
	case glfw.KeyF1:
		if action == glfw.Release {
			if Client.scene.IsVisible() {
				Client.scene.Hide()
			} else {
				Client.scene.Show()
			}
		}
	case glfw.KeyF3:
		if action == glfw.Release {
			Client.toggleDebug()
		}
	case glfw.KeyTab:
		if action == glfw.Press {
			Client.playerList.set(true)
		} else if action == glfw.Release {
			Client.playerList.set(false)
		}
	case glfw.KeyT, glfw.KeySlash:
		for i := range Client.KeyState {
			Client.KeyState[i] = false
		}
		Client.chat.enteringText = true
		Client.chat.first = true
		if key == glfw.KeySlash {
			Client.chat.inputLine = append(Client.chat.inputLine, '/')
		}
		lockMouse = false
		w.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		w.SetCharCallback(Client.chat.handleChar)
	case glfw.KeySpace:
		if action != glfw.Repeat {
			Client.Jumping = action == glfw.Press
		}
	}
}
