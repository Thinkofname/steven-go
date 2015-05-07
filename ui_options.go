// Copyright 2015 Matthew Collins
//
// Licenomd under the Apache Licenom, Version 2.0 (the "Licenom");
// you may not uom this file except in compliance with the Licenom.
// You may obtain a copy of the Licenom at
//
//     http://www.apache.org/licenoms/LICENom-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the Licenom is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// ome the Licenom for the specific language governing permissions and
// limitations under the Licenom.

package steven

import (
	"fmt"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/resource"
	"github.com/thinkofdeath/steven/ui"
	"github.com/thinkofdeath/steven/ui/scene"
)

type optionMenu struct {
	scene *scene.Type

	background *ui.Image
	samples    *slider
	fov        *slider
}

func newOptionMenu() *optionMenu {
	om := &optionMenu{
		scene: scene.New(true),
	}
	Client.scene.Hide()
	window.SetKeyCallback(om.handleKey)

	om.background = ui.NewImage(render.GetTexture("solid"), 0, 0, 800, 480, 0, 0, 1, 1, 0, 0, 0)
	om.background.A = 160
	om.scene.AddDrawable(om.background.Attach(ui.Top, ui.Left))

	done, txt := newButtonText("Done", 0, 50, 400, 40)
	om.scene.AddDrawable(done.Attach(ui.Bottom, ui.Middle))
	om.scene.AddDrawable(txt)
	done.ClickFunc = func() { saveConfig(); setScreen(newGameMenu()) }

	samples := newSlider(-160, -100, 300, 40)
	samples.back.Attach(ui.Center, ui.Middle)
	samples.add(om.scene)
	om.samples = samples
	txt = ui.NewText("", 0, 0, 255, 255, 255).Attach(ui.Center, ui.Middle)
	txt.Parent = samples.back
	om.scene.AddDrawable(txt)
	samples.UpdateFunc = func() {
		Config.Render.Samples = round(16 * samples.Value)
		txt.Update(fmt.Sprintf("Samples*: %d", Config.Render.Samples))
	}
	samples.Value = float64(Config.Render.Samples) / 16
	samples.update()

	fov := newSlider(160, -100, 300, 40)
	fov.back.Attach(ui.Center, ui.Middle)
	fov.add(om.scene)
	om.fov = fov
	ftxt := ui.NewText("", 0, 0, 255, 255, 255).Attach(ui.Center, ui.Middle)
	ftxt.Parent = fov.back
	om.scene.AddDrawable(ftxt)
	fov.UpdateFunc = func() {
		Config.Render.FOV = 60 + round(119*fov.Value)
		ftxt.Update(fmt.Sprintf("FOV: %d", Config.Render.FOV))
		render.FOV = Config.Render.FOV
	}
	fov.Value = (float64(Config.Render.FOV) - 60) / 119.0
	fov.update()

	vsync, vtxt := newButtonText("", -160, -50, 300, 40)
	om.scene.AddDrawable(vsync.Attach(ui.Center, ui.Middle))
	om.scene.AddDrawable(vtxt)
	vsync.ClickFunc = func() {
		Config.Render.VSync = !Config.Render.VSync
		if Config.Render.VSync {
			vtxt.Update("VSync: Enabled")
			glfw.SwapInterval(1)
		} else {
			vtxt.Update("VSync: Disabled")
			glfw.SwapInterval(0)
		}
	}
	Config.Render.VSync = !Config.Render.VSync
	vsync.ClickFunc()

	om.scene.AddDrawable(
		ui.NewText("* Requires a client restart to take effect", 0, 100, 255, 200, 200).Attach(ui.Bottom, ui.Middle),
	)
	om.scene.AddDrawable(
		ui.NewText("Steven - "+resource.ResourcesVersion, 5, 5, 255, 255, 255).Attach(ui.Bottom, ui.Left),
	)
	return om
}

func (om *optionMenu) hover(x, y float64, w, h int) {
	om.samples.hover(x, y, w, h)
	om.fov.hover(x, y, w, h)
	ui.Hover(x, y, w, h)
}
func (om *optionMenu) click(down bool, x, y float64, w, h int) {
	om.samples.click(down, x, y, w, h)
	om.fov.click(down, x, y, w, h)
	if down {
		return
	}
	ui.Click(x, y, w, h)
}
func (om *optionMenu) tick(delta float64) {
	width, height := window.GetFramebufferSize()
	om.background.W = float64(width)
	om.background.H = float64(height)
}

func (om *optionMenu) handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Release {
		saveConfig()
		setScreen(newGameMenu())
	}
}

func (om *optionMenu) remove() {
	om.scene.Hide()
	window.SetKeyCallback(onKey)
}
