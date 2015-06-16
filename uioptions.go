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
	"github.com/thinkofdeath/steven/render/ui"
	"github.com/thinkofdeath/steven/render/ui/scene"
)

type optionMenu struct {
	baseUI
	scene *scene.Type

	background *ui.Image
	fov        *slider
	mouseS     *slider

	ret func() screen
}

func newOptionMenu(ret func() screen) *optionMenu {
	om := &optionMenu{
		scene: scene.New(true),
		ret:   ret,
	}

	om.background = ui.NewImage(render.GetTexture("solid"), 0, 0, 854, 480, 0, 0, 1, 1, 0, 0, 0)
	om.background.SetA(160)
	om.scene.AddDrawable(om.background.Attach(ui.Top, ui.Left))

	done, txt := newButtonText("Done", 0, 50, 400, 40)
	om.scene.AddDrawable(done.Attach(ui.Bottom, ui.Middle))
	om.scene.AddDrawable(txt)
	done.ClickFunc = func() { setScreen(om.ret()) }

	rp, txt := newButtonText("Resource packs", -160, 150, 300, 40)
	om.scene.AddDrawable(rp.Attach(ui.Bottom, ui.Middle))
	om.scene.AddDrawable(txt)
	rp.ClickFunc = func() { setScreen(newResourceList(om.ret)) }

	fov := newSlider(160, -100, 300, 40)
	fov.back.Attach(ui.Center, ui.Middle)
	fov.add(om.scene)
	om.fov = fov
	ftxt := ui.NewText("", 0, 0, 255, 255, 255).Attach(ui.Center, ui.Middle)
	ftxt.AttachTo(fov.back)
	om.scene.AddDrawable(ftxt)
	fov.UpdateFunc = func() {
		render.FOV.SetValue(60 + round(119*fov.Value))
		ftxt.Update(fmt.Sprintf("FOV: %d", render.FOV.Value()))
	}
	fov.Value = (float64(render.FOV.Value()) - 60) / 119.0
	fov.update()

	vsync, vtxt := newButtonText("", -160, -50, 300, 40)
	om.scene.AddDrawable(vsync.Attach(ui.Center, ui.Middle))
	om.scene.AddDrawable(vtxt)
	vsync.ClickFunc = func() {
		renderVSync.SetValue(!renderVSync.Value())
		if renderVSync.Value() {
			vtxt.Update("VSync: Enabled")
		} else {
			vtxt.Update("VSync: Disabled")
		}
	}
	renderVSync.SetValue(!renderVSync.Value())
	vsync.ClickFunc()

	mouseS := newSlider(160, -50, 300, 40)
	mouseS.back.Attach(ui.Center, ui.Middle)
	mouseS.add(om.scene)
	om.mouseS = mouseS
	mtxt := ui.NewText("", 0, 0, 255, 255, 255).Attach(ui.Center, ui.Middle)
	mtxt.AttachTo(mouseS.back)
	om.scene.AddDrawable(mtxt)
	mouseS.UpdateFunc = func() {
		mouseSensitivity.SetValue(500 + round(9500.0*mouseS.Value))
		mtxt.Update(fmt.Sprintf("Mouse Speed: %d", mouseSensitivity.Value()))
	}
	mouseS.Value = (float64(mouseSensitivity.Value()) - 500) / 9500.0
	mouseS.update()

	om.scene.AddDrawable(
		ui.NewText("* Requires a client restart to take effect", 0, 100, 255, 200, 200).Attach(ui.Bottom, ui.Middle),
	)

	uiFooter(om.scene)

	scales := []string{
		uiAuto, uiSmall, uiMedium, uiLarge,
	}
	curScale := func() int {
		for i, s := range scales {
			if s == uiScale.Value() {
				return i
			}
		}
		return 0
	}

	uiS, utxt := newButtonText("", -160, 0, 300, 40)
	om.scene.AddDrawable(uiS.Attach(ui.Center, ui.Middle))
	om.scene.AddDrawable(utxt)
	uiS.ClickFunc = func() {
		uiScale.SetValue(scales[(curScale()+1)%len(scales)])
		utxt.Update(fmt.Sprintf("UI Scale: %s", uiScale.Value()))
	}
	uiScale.SetValue(scales[(len(scales)+curScale()-1)%len(scales)])
	uiS.ClickFunc()

	return om
}

func (om *optionMenu) init() {
	window.SetKeyCallback(om.handleKey)
}

func (om *optionMenu) hover(x, y float64, w, h int) {
	om.fov.hover(x, y, w, h)
	om.mouseS.hover(x, y, w, h)
	ui.Hover(x, y, w, h)
}
func (om *optionMenu) click(down bool, x, y float64, w, h int) {
	om.fov.click(down, x, y, w, h)
	om.mouseS.click(down, x, y, w, h)
	if down {
		return
	}
	ui.Click(x, y, w, h)
}
func (om *optionMenu) tick(delta float64) {
	width, height := window.GetFramebufferSize()
	om.background.SetWidth(float64(width) / ui.Scale)
	om.background.SetHeight(float64(height) / ui.Scale)
}

func (om *optionMenu) handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Release {
		setScreen(om.ret())
	}
}

func (om *optionMenu) remove() {
	om.scene.Hide()
	window.SetKeyCallback(onKey)
}
