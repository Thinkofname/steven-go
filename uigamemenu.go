// Copyright 2015 Matthew Collins
//
// Licengmd under the Apache Licengm, Version 2.0 (the "Licengm");
// you may not ugm this file except in compliance with the Licengm.
// You may obtain a copy of the Licengm at
//
//     http://www.apache.org/licengms/LICENgm-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the Licengm is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// gme the Licengm for the specific language governing permissions and
// limitations under the Licengm.

package steven

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/ui"
	"github.com/thinkofdeath/steven/ui/scene"
)

type gameMenu struct {
	baseUI
	scene *scene.Type

	background *ui.Image
}

func newGameMenu() screen {
	gm := &gameMenu{
		scene: scene.New(true),
	}

	gm.background = ui.NewImage(render.GetTexture("solid"), 0, 0, 854, 480, 0, 0, 1, 1, 0, 0, 0)
	gm.background.SetA(160)
	gm.scene.AddDrawable(gm.background.Attach(ui.Top, ui.Left))

	disconnect, txt := newButtonText("Disconnect", 0, 50, 400, 40)
	gm.scene.AddDrawable(disconnect.Attach(ui.Center, ui.Middle))
	gm.scene.AddDrawable(txt)
	disconnect.ClickFunc = func() { Client.network.SignalClose(errManualDisconnect) }

	rtg, txt := newButtonText("Return to game", 0, -50, 400, 40)
	gm.scene.AddDrawable(rtg.Attach(ui.Center, ui.Middle))
	gm.scene.AddDrawable(txt)
	rtg.ClickFunc = func() { setScreen(nil) }

	option, txt := newButtonText("Options", 0, 0, 400, 40)
	gm.scene.AddDrawable(option.Attach(ui.Center, ui.Middle))
	gm.scene.AddDrawable(txt)
	option.ClickFunc = func() { setScreen(newOptionMenu(newGameMenu)) }

	uiFooter(gm.scene)
	return gm
}

func (gm *gameMenu) init() {
	window.SetKeyCallback(gm.handleKey)
}

func (gm *gameMenu) tick(delta float64) {
	width, height := window.GetFramebufferSize()
	gm.background.SetWidth(float64(width) / ui.Scale)
	gm.background.SetHeight(float64(height) / ui.Scale)
}

func (gm *gameMenu) handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Release {
		setScreen(nil)
	}
}

func (gm *gameMenu) remove() {
	gm.scene.Hide()
	window.SetKeyCallback(onKey)
}
