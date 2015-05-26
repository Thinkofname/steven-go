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
	"io"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/render/gl"
)

var window *glfw.Window

func init() {
	runtime.LockOSThread()

	f, err := os.Create("steven-log.txt")
	if err != nil {
		panic(err)
	}
	log.SetOutput(io.MultiWriter(f, os.Stdout))
	log.SetFlags(log.Lshortfile | log.Ltime)
}

func startWindow() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, Config.Render.Samples)
	if os.Getenv("STEVEN_DEBUG") == "true" {
		glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)
	}
	render.MultiSample = Config.Render.Samples > 0

	var err error
	window, err = glfw.CreateWindow(854, 480, "Steven", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	if Config.Render.VSync {
		glfw.SwapInterval(1)
	} else {
		glfw.SwapInterval(0)
	}

	window.SetCursorPosCallback(onMouseMove)
	window.SetMouseButtonCallback(onMouseClick)
	window.SetKeyCallback(onKey)
	window.SetScrollCallback(onScroll)
	window.SetFocusCallback(onFocus)

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

	Client.network.Write(&protocol.HeldItemChange{Slot: int16(Client.currentHotbarSlot)})
}

var lockMouse bool

func onFocus(w *glfw.Window, focused bool) {
	if !focused {
		for i := range Client.KeyState {
			Client.KeyState[i] = false
		}
	} else if lockMouse {
		w.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	}
}

func onMouseMove(w *glfw.Window, xpos float64, ypos float64) {
	width, height := w.GetSize()
	if currentScreen != nil {
		fw, fh := w.GetFramebufferSize()
		currentScreen.hover(xpos*(float64(fw)/float64(width)), ypos*(float64(fh)/float64(height)), fw, fh)
		return
	}
	if !lockMouse {
		return
	}
	ww, hh := float64(width/2), float64(height/2)
	w.SetCursorPos(ww, hh)

	s := float64(Config.Game.MouseSensitivity)
	rotate((xpos-ww)/s, (ypos-hh)/s)
}

func onMouseClick(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if currentScreen != nil {
		if button != glfw.MouseButtonLeft || action == glfw.Repeat {
			return
		}
		width, height := w.GetSize()
		xpos, ypos := w.GetCursorPos()
		fw, fh := w.GetFramebufferSize()
		currentScreen.click(action == glfw.Press, xpos*(float64(fw)/float64(width)), ypos*(float64(fh)/float64(height)), fw, fh)
		return
	}
	if !Client.chat.enteringText && lockMouse && action != glfw.Repeat {
		Client.MouseAction(button, action == glfw.Press)
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
	KeyJump
)

var keyStateMap = map[glfw.Key]Key{
	glfw.KeyW:           KeyForward,
	glfw.KeyS:           KeyBackwards,
	glfw.KeyA:           KeyLeft,
	glfw.KeyD:           KeyRight,
	glfw.KeyLeftControl: KeySprint,
	glfw.KeySpace:       KeyJump,
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if currentScreen != nil {
		return
	}
	if Client.chat.enteringText {
		Client.chat.handleKey(w, key, scancode, action, mods)
		return
	}

	if k, ok := keyStateMap[key]; action != glfw.Repeat && ok {
		Client.KeyState[k] = action == glfw.Press
	}
	switch key {
	case glfw.KeyEscape:
		if action == glfw.Release {
			setScreen(newGameMenu())
		}
	case glfw.KeyF1:
		if action == glfw.Release {
			if Client.scene.IsVisible() {
				Client.scene.Hide()
				Client.hotbarScene.Hide()
			} else {
				Client.scene.Show()
				Client.hotbarScene.Show()
			}
		}
	case glfw.KeyF3:
		if action == glfw.Release {
			Client.toggleDebug()
		}
	case glfw.KeyF5:
		if action == glfw.Release {
			Client.cycleCamera()
		}
	case glfw.KeyTab:
		if action == glfw.Press {
			Client.playerList.set(true)
		} else if action == glfw.Release {
			Client.playerList.set(false)
		}
	case glfw.KeyT:
		state := w.GetKey(glfw.KeyF3)
		if action == glfw.Release && state == glfw.Press {
			reloadResources()
			return
		}
		fallthrough
	case glfw.KeySlash:
		if action != glfw.Release {
			return
		}
		for i := range Client.KeyState {
			Client.KeyState[i] = false
		}
		Client.chat.enteringText = true
		if key == glfw.KeySlash {
			Client.chat.inputLine = append(Client.chat.inputLine, '/')
		}
		lockMouse = false
		w.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		w.SetCharCallback(Client.chat.handleChar)
	}
}
