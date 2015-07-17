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
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/console"
	"github.com/thinkofdeath/steven/format"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/ui"
	"github.com/thinkofdeath/steven/ui/scene"
)

var con consoleScreen

type consoleScreen struct {
	scene *scene.Type

	container  *ui.Container
	background *ui.Image
	inputText  *ui.Text

	lastLine  format.AnyComponent
	lastWidth float64

	lines []*ui.Formatted

	pos     float64
	visible bool

	prevKey  glfw.KeyCallback
	prevChar glfw.CharCallback

	input      string
	cursorTick float64
}

func (cs *consoleScreen) init() {
	cs.scene = scene.New(true)
	cs.container = ui.NewContainer(0, -220, 854, 220)
	cs.container.SetLayer(-200)
	cs.background = ui.NewImage(render.GetTexture("solid"), 0, 0, 854, 220, 0, 0, 1, 1, 0, 0, 0)
	cs.background.SetA(180)
	cs.background.AttachTo(cs.container)
	cs.scene.AddDrawable(cs.background.Attach(ui.Top, ui.Left))

	cs.inputText = ui.NewText("", 5, 200, 255, 255, 255)
	cs.inputText.AttachTo(cs.container)
	cs.scene.AddDrawable(cs.inputText.Attach(ui.Top, ui.Left))
	cs.pos = -220
	cs.visible = false
}

func (cs *consoleScreen) focus() {
	cs.visible = true

	cs.prevKey = window.SetKeyCallback(cs.onKey)
	cs.prevChar = window.SetCharCallback(cs.onChar)
}

func (cs *consoleScreen) onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyGraveAccent && action == glfw.Release {
		cs.visible = false
		window.SetKeyCallback(cs.prevKey)
		window.SetCharCallback(cs.prevChar)
		cs.input = cs.input[:len(cs.input)-1]
	}
	if key == glfw.KeyBackspace && action != glfw.Release {
		if len(cs.input) > 0 {
			cs.input = cs.input[:len(cs.input)-1]
		}
	}
	if key == glfw.KeyEnter && action == glfw.Release {
		console.Component(format.
			Build("> ").
			Color(format.Yellow).
			Append(cs.input).
			Create(),
		)
		err := console.Execute(cs.input)
		cs.input = ""
		if err != nil {
			console.Component(format.
				Build("Error").
				Color(format.Red).
				Append(": ").
				Color(format.Yellow).
				Append(err.Error()).
				Color(format.Red).
				Create(),
			)
		}
	}
}

func (cs *consoleScreen) onChar(w *glfw.Window, char rune) {
	cs.input += string(char)
}

func (cs *consoleScreen) tick(delta float64) {
	if cs.visible {
		if cs.pos == -220 {
			cs.lastWidth = -1 // Force a redraw
		}
		if cs.pos < 0 {
			cs.pos += delta * 4
		} else {
			cs.pos = 0
		}
	} else {
		if cs.pos > -220 {
			cs.pos -= delta * 4
		} else {
			cs.pos = -220
		}
		if cs.pos == -220 {
			cs.scene.Hide()
			return
		}
	}
	cs.container.SetY(cs.pos)
	// Not the most efficent way but handles the alpha issue for
	// now
	cs.scene.Hide()
	cs.scene.Show()

	// Resize to fill the screen width
	width, _ := window.GetFramebufferSize()
	sw := 854 / float64(width)
	var w float64
	if ui.DrawMode == ui.Unscaled {
		sw = ui.Scale
		w = 854 / sw
	} else {
		w = float64(width)
	}
	cs.container.SetWidth(w)
	cs.background.SetWidth(w)

	if cs.lastLine != console.History(1)[0] || cs.lastWidth != w {
		for _, l := range cs.lines {
			l.Remove()
		}
		cs.lines = cs.lines[:0]
		hist := console.History(-1)
		offset := 200.0
		for i := 0; i < len(hist) && offset > 0; i++ {
			line := hist[i]
			if line.Value == nil {
				break
			}
			fmt := ui.NewFormattedWidth(line, 5, 0, w-10)
			offset -= fmt.Height
			fmt.SetY(offset)
			fmt.AttachTo(cs.container)
			cs.lines = append(cs.lines, fmt)
			cs.scene.AddDrawable(fmt.Attach(ui.Top, ui.Left))
		}
	}

	cs.lastWidth = w
	cs.lastLine = console.History(1)[0]

	cs.cursorTick += delta
	// Add on our cursor
	if int(cs.cursorTick/30)%2 == 0 {
		cs.inputText.Update(string(cs.input) + "|")
	} else {
		cs.inputText.Update(string(cs.input))
	}
}
