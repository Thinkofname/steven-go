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
	"math"
	"strings"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/console"
	"github.com/thinkofdeath/steven/render/ui"
	"github.com/thinkofdeath/steven/render/ui/scene"
)

const (
	uiAuto   = "auto"
	uiSmall  = "small"
	uiMedium = "medium"
	uiLarge  = "large"
)

var uiScale = console.NewStringVar("cl_ui_scale", "auto", console.Mutable, console.Serializable).
	Doc(`
cl_ui_scale sets the scaling used for the user interface. 
Valid values are:
- auto
- small
- medium
- large
`)

func init() {
	uiScale.Callback(func() {
		setUIScale()
	})
}

func setUIScale() {
	switch uiScale.Value() {
	case uiAuto:
		ui.DrawMode = ui.Scaled
		ui.Scale = 1.0
	case uiSmall:
		ui.DrawMode = ui.Unscaled
		ui.Scale = 0.4
	case uiMedium:
		ui.DrawMode = ui.Unscaled
		ui.Scale = 0.6
	case uiLarge:
		ui.DrawMode = ui.Unscaled
		ui.Scale = 1.0
	}
	ui.ForceDraw()
}

func uiFooter(scene *scene.Type) {
	scene.AddDrawable(
		ui.NewText("Steven - "+stevenVersion(), 5, 5, 255, 255, 255).Attach(ui.Bottom, ui.Left),
	)
	scene.AddDrawable(
		ui.NewText("Not affiliated with Mojang/Minecraft", 5, 5, 255, 200, 200).Attach(ui.Bottom, ui.Right),
	)
}

type baseUI struct{}

func (b *baseUI) init()                        {}
func (b *baseUI) hover(x, y float64, w, h int) { ui.Hover(x, y, w, h) }
func (b *baseUI) click(down bool, x, y float64, w, h int) {
	if down {
		return
	}
	ui.Click(x, y, w, h)
}

func newButtonText(str string, x, y, w, h float64) (*ui.Button, *ui.Text) {
	btn := ui.NewButton(x, y, w, h)
	text := ui.NewText(str, 0, 0, 255, 255, 255).Attach(ui.Middle, ui.Center)
	text.AttachTo(btn)
	btn.AddHover(func(over bool) {
		if over && !btn.Disabled() {
			text.SetB(160)
		} else {
			text.SetB(255)
		}
	})
	btn.AddClick(func() { PlaySound("random.click") })
	return btn, text
}

type textBox struct {
	back       *ui.Button
	text       *ui.Text
	Password   bool
	Focused    bool
	wasFocused bool
	input      string
	cursorTick float64
}

func newTextBox(x, y, w, h float64) *textBox {
	btn := ui.NewButton(x, y, w, h)
	btn.SetDisabled(true)
	text := ui.NewText("", 5, 0, 255, 255, 255).Attach(ui.Middle, ui.Left)
	text.AttachTo(btn)
	return &textBox{
		back: btn,
		text: text,
	}
}

func (t *textBox) value() string {
	if t.Password {
		return strings.Repeat("*", len(t.input))
	}
	return t.input
}

func (t *textBox) add(s *scene.Type) {
	s.AddDrawable(t.back)
	s.AddDrawable(t.text)
}

func (t *textBox) handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyBackspace && action != glfw.Release {
		if len(t.input) > 0 {
			t.input = t.input[:len(t.input)-1]
			t.text.Update(t.value())
		}
	}
}

func (t *textBox) handleChar(w *glfw.Window, char rune) {
	t.input += string(char)
	t.text.Update(t.value())
}

func (t *textBox) tick(delta float64) {
	if !t.Focused {
		if t.wasFocused {
			t.wasFocused = t.Focused
			t.text.Update(t.value())
		}
		return
	}
	t.wasFocused = true
	t.cursorTick += delta
	if int(t.cursorTick/30)%2 == 0 {
		t.text.Update(t.value() + "|")
	} else {
		t.text.Update(t.value())
	}
	// Lazy way of preventing rounding errors buiding up over time
	if t.cursorTick > 0xFFFFFF {
		t.cursorTick = 0
	}
}

type slider struct {
	back       *ui.Button
	slider     *ui.Button
	Value      float64
	sliding    bool
	UpdateFunc func()
}

func newSlider(x, y, w, h float64) *slider {
	btn := ui.NewButton(x, y, w, h)
	btn.SetDisabled(true)
	sl := ui.NewButton(0, 0, 20, h).Attach(ui.Left, ui.Top)
	sl.AttachTo(btn)
	return &slider{
		back:   btn,
		slider: sl,
	}
}

func (sl *slider) update() {
	ww, _ := sl.back.Size()
	sl.slider.SetX(sl.Value * (ww - 20))
	if sl.UpdateFunc != nil {
		sl.UpdateFunc()
	}
}

func (sl *slider) click(down bool, x, y float64, w, h int) {
	if !down {
		if sl.sliding {
			sl.sliding = false
		}
		return
	}
	_, _, ok := ui.Intersects(sl.slider, x, y, w, h)
	if ok {
		sl.sliding = true
	} else {
		ox, _, ok := ui.Intersects(sl.back, x, y, w, h)
		if ok {
			ww, _ := sl.back.Size()
			v := math.Min(ww-10, math.Max(10, ox)) - 10
			sl.Value = v / (ww - 20)
			sl.update()
		}
	}
}
func (sl *slider) hover(x, y float64, w, h int) {
	if sl.sliding {
		ox, _, ok := ui.Intersects(sl.back, x, y, w, h)
		if ok {
			ww, _ := sl.back.Size()
			v := math.Min(ww-10, math.Max(10, ox)) - 10
			sl.Value = v / (ww - 20)
			sl.update()
		}
	}
}

func (sl *slider) add(s *scene.Type) {
	s.AddDrawable(sl.back)
	s.AddDrawable(sl.slider)
}
