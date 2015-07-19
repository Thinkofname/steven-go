// Copyright 2015 Matthew Collins
//
// Licenvmd under the Apache Licenvm, Version 2.0 (the "Licenvm");
// you may not uvm this file except in cvmpliance with the Licenvm.
// You may obtain a copy of the Licenvm at
//
//     http://www.apache.org/licenvms/LICENvm-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the Licenvm is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// vme the Licenvm for the specific language governing permissions and
// limitations under the Licenvm.

package steven

import (
	"fmt"
	"strings"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/ui"
	"github.com/thinkofdeath/steven/ui/scene"
)

type volumeMenu struct {
	baseUI
	scene *scene.Type

	background *ui.Image
	sliders    []*slider

	ret func() screen
}

func newVolumeMenu(ret func() screen) *volumeMenu {
	vm := &volumeMenu{
		scene: scene.New(true),
		ret:   ret,
	}

	vm.background = ui.NewImage(render.GetTexture("solid"), 0, 0, 854, 480, 0, 0, 1, 1, 0, 0, 0)
	vm.background.SetA(160)
	vm.scene.AddDrawable(vm.background.Attach(ui.Top, ui.Left))

	done, txt := newButtonText("Done", 0, 50, 400, 40)
	vm.scene.AddDrawable(done.Attach(ui.Bottom, ui.Middle))
	vm.scene.AddDrawable(txt)
	done.AddClick(func() { setScreen(newOptionMenu(vm.ret)) })

	master := newSlider(0, -100, 620, 40)
	master.back.Attach(ui.Center, ui.Middle)
	master.add(vm.scene)
	mtxt := ui.NewText("", 0, 0, 255, 255, 255).Attach(ui.Center, ui.Middle)
	mtxt.AttachTo(master.back)
	vm.scene.AddDrawable(mtxt)
	master.UpdateFunc = func() {
		muVolMaster.SetValue(round(master.Value * 100))
		if muVolMaster.Value() == 0 {
			mtxt.Update("Master: OFF")
			return
		}
		mtxt.Update(fmt.Sprintf("Master: %d%%", muVolMaster.Value()))
	}
	master.Value = float64(muVolMaster.Value()) / 100.0
	master.update()
	vm.sliders = append(vm.sliders, master)

	for i, cat := range soundCategories {
		cat := cat
		x := 160.0
		if i&1 == 0 {
			x = -x
		}
		y := 50 * float64(i/2)
		snd := newSlider(x, -50+y, 300, 40)
		snd.back.Attach(ui.Center, ui.Middle)
		snd.add(vm.scene)
		stxt := ui.NewText("", 0, 0, 255, 255, 255).Attach(ui.Center, ui.Middle)
		stxt.AttachTo(snd.back)
		vm.scene.AddDrawable(stxt)

		v := volVars[cat]

		snd.UpdateFunc = func() {
			v.SetValue(round(snd.Value * 100))
			if val := v.Value(); val != 0 {
				stxt.Update(fmt.Sprintf("%s: %d%%", strings.Title(string(cat)), val))
				return
			}
			stxt.Update(fmt.Sprintf("%s: OFF", strings.Title(string(cat))))
		}
		snd.Value = float64(v.Value()) / 100.0
		snd.update()
		vm.sliders = append(vm.sliders, snd)
	}

	return vm
}

func (vm *volumeMenu) init() {
	window.SetKeyCallback(vm.handleKey)
}

func (vm *volumeMenu) hover(x, y float64, w, h int) {
	for _, s := range vm.sliders {
		s.hover(x, y, w, h)
	}
	ui.Hover(x, y, w, h)
}
func (vm *volumeMenu) click(down bool, x, y float64, w, h int) {
	for _, s := range vm.sliders {
		s.click(down, x, y, w, h)
	}
	if down {
		return
	}
	ui.Click(x, y, w, h)
}
func (vm *volumeMenu) tick(delta float64) {
	width, height := window.GetFramebufferSize()
	vm.background.SetWidth(float64(width) / ui.Scale)
	vm.background.SetHeight(float64(height) / ui.Scale)
}

func (vm *volumeMenu) handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Release {
		setScreen(newOptionMenu(vm.ret))
	}
}

func (vm *volumeMenu) remove() {
	vm.scene.Hide()
	window.SetKeyCallback(onKey)
}
