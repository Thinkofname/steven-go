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
	"github.com/thinkofdeath/steven/ui"
	"github.com/thinkofdeath/steven/ui/scene"
)

type editServer struct {
	baseUI
	scene *scene.Type
	logo  uiLogo

	name    *ui.TextBox
	address *ui.TextBox

	index int
}

func newEditServer(index int) *editServer {
	se := &editServer{
		scene: scene.New(true),
		index: index,
	}

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

	se.name = ui.NewTextBox(0, -20, 400, 40)
	se.scene.AddDrawable(se.name.Attach(ui.Middle, ui.Center))
	label := ui.NewText("Name:", 0, -18, 255, 255, 255).Attach(ui.Top, ui.Left)
	label.AttachTo(se.name)
	se.scene.AddDrawable(label)

	se.address = ui.NewTextBox(0, 40, 400, 40)
	se.scene.AddDrawable(se.address.Attach(ui.Middle, ui.Center))
	label = ui.NewText("Address:", 0, -18, 255, 255, 255).Attach(ui.Top, ui.Left)
	label.AttachTo(se.address)
	se.scene.AddDrawable(label)

	if index != -1 {
		server := Config.Servers[index]
		se.name.Update(server.Name)
		se.address.Update(server.Address)
	}

	return se
}

func (se *editServer) save() {
	if se.index == -1 {
		Config.Servers = append(Config.Servers, ConfigServer{
			Name:    se.name.Value(),
			Address: se.address.Value(),
		})
	} else {
		Config.Servers[se.index].Name = se.name.Value()
		Config.Servers[se.index].Address = se.address.Value()
	}
	saveServers()
	setScreen(newServerList())
}

func (se *editServer) tick(delta float64) {
	se.logo.tick(delta)
}

func (se *editServer) remove() {
	se.scene.Hide()
}
