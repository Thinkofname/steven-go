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
	"github.com/thinkofdeath/steven/resource"
	"github.com/thinkofdeath/steven/ui"
	"github.com/thinkofdeath/steven/ui/scene"
)

type gameMenu struct {
	scene *scene.Type

	background *ui.Image
}

func newGameMenu() *gameMenu {
	gm := &gameMenu{
		scene: scene.New(true),
	}
	Client.scene.Hide()
	window.SetKeyCallback(gm.handleKey)

	gm.background = ui.NewImage(render.GetTexture("solid"), 0, 0, 800, 480, 0, 0, 1, 1, 0, 0, 0)
	gm.background.A = 160
	gm.scene.AddDrawable(gm.background.Attach(ui.Top, ui.Left))

	disconnect, txt := newButtonText("Disconnect", 0, 50, 400, 40)
	gm.scene.AddDrawable(disconnect.Attach(ui.Center, ui.Middle))
	gm.scene.AddDrawable(txt)
	disconnect.ClickFunc = func() { closeWithError(errManualDisconnect) }

	rtg, txt := newButtonText("Return to game", 0, -50, 400, 40)
	gm.scene.AddDrawable(rtg.Attach(ui.Center, ui.Middle))
	gm.scene.AddDrawable(txt)
	rtg.ClickFunc = func() { setScreen(nil) }

	option, txt := newButtonText("Options", 0, 0, 400, 40)
	gm.scene.AddDrawable(option.Attach(ui.Center, ui.Middle))
	gm.scene.AddDrawable(txt)

	gm.scene.AddDrawable(
		ui.NewText("Steven - "+resource.ResourcesVersion, 5, 5, 255, 255, 255).Attach(ui.Bottom, ui.Left),
	)
	return gm
}

func (gm *gameMenu) hover(x, y float64, w, h int) {
	ui.Hover(x, y, w, h)
}
func (gm *gameMenu) click(x, y float64, w, h int) {
	ui.Click(x, y, w, h)
}
func (gm *gameMenu) tick(delta float64) {
	width, height := window.GetFramebufferSize()
	gm.background.W = float64(width)
	gm.background.H = float64(height)
}

func (m *gameMenu) handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Release {
		setScreen(nil)
	}
}

func (gm *gameMenu) remove() {
	gm.scene.Hide()
	window.SetKeyCallback(onKey)
}
