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
	"github.com/thinkofdeath/steven/render/ui"
	"github.com/thinkofdeath/steven/render/ui/scene"
)

type editServer struct {
	baseUI
	scene *scene.Type
	logo  uiLogo

	name    *textBox
	address *textBox
	focused *textBox

	index int
}

func newEditServer(index int) *editServer {
	se := &editServer{
		scene: scene.New(true),
		index: index,
	}

	// For the text boxes
	window.SetKeyCallback(se.handleKey)
	window.SetCharCallback(se.handleChar)
	se.logo.init(se.scene)

	uiFooter(se.scene)

	done, txt := newButtonText("Done", 110, 100, 200, 40)
	se.scene.AddDrawable(done.Attach(ui.Center, ui.Middle))
	se.scene.AddDrawable(txt)
	done.AddClick(func() {
		se.save()
	})

	cancel, txt := newButtonText("Cancel", -110, 100, 200, 40)
	se.scene.AddDrawable(cancel.Attach(ui.Center, ui.Middle))
	se.scene.AddDrawable(txt)
	cancel.AddClick(func() {
		setScreen(newServerList())
	})

	se.name = newTextBox(0, -20, 400, 40)
	se.name.back.Attach(ui.Middle, ui.Center)
	se.name.add(se.scene)
	label := ui.NewText("Name:", 0, -18, 255, 255, 255).Attach(ui.Top, ui.Left)
	label.AttachTo(se.name.back)
	se.scene.AddDrawable(label)
	se.name.back.AddClick(func() {
		if se.focused != nil {
			se.focused.Focused = false
		}
		se.name.Focused = true
		se.focused = se.name
	})

	se.address = newTextBox(0, 40, 400, 40)
	se.address.back.Attach(ui.Middle, ui.Center)
	se.address.add(se.scene)
	label = ui.NewText("Address:", 0, -18, 255, 255, 255).Attach(ui.Top, ui.Left)
	label.AttachTo(se.address.back)
	se.scene.AddDrawable(label)
	se.address.back.AddClick(func() {
		if se.focused != nil {
			se.focused.Focused = false
		}
		se.address.Focused = true
		se.focused = se.address
	})

	if index != -1 {
		server := Config.Servers[index]
		se.name.input = server.Name
		se.name.text.Update(se.name.input)
		se.address.input = server.Address
		se.address.text.Update(se.address.input)
	}

	return se
}

func (se *editServer) save() {
	if se.index == -1 {
		Config.Servers = append(Config.Servers, ConfigServer{
			Name:    se.name.input,
			Address: se.address.input,
		})
	} else {
		Config.Servers[se.index].Name = se.name.input
		Config.Servers[se.index].Address = se.address.input
	}
	saveServers()
	setScreen(newServerList())
}

func (se *editServer) tick(delta float64) {
	se.logo.tick(delta)
	se.name.tick(delta)
	se.address.tick(delta)
}

func (se *editServer) handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if se.focused == nil {
		return
	}

	if (key == glfw.KeyEnter || key == glfw.KeyTab) && action == glfw.Release {
		if se.focused == se.name {
			se.name.Focused = false
			se.focused = se.address
			se.address.Focused = true
		} else if se.focused == se.address {
			se.address.Focused = false
			se.focused = nil
			se.save()
		}
		return
	}

	if key == glfw.KeyEscape && action == glfw.Release {
		se.focused.Focused = false
		se.focused = nil
	}

	se.focused.handleKey(w, key, scancode, action, mods)
}

func (se *editServer) handleChar(w *glfw.Window, char rune) {
	if se.focused == nil {
		return
	}
	se.focused.handleChar(w, char)
}

func (se *editServer) remove() {
	se.scene.Hide()
	window.SetKeyCallback(onKey)
	window.SetCharCallback(nil)
}
