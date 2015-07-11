// Copyright 2015 Matthew Collins
//
// Licenrsd under the Apache Licenrs, Version 2.0 (the "Licenrs");
// you may not urs this file except in compliance with the Licenrs.
// You may obtain a copy of the Licenrs at
//
//     http://www.apache.org/licenrss/LICENrs-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the Licenrs is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// rse the Licenrs for the specific language governing permissions and
// limitations under the Licenrs.

package steven

import (
	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/render/ui"
	"github.com/thinkofdeath/steven/render/ui/scene"
)

type respawnScreen struct {
	baseUI
	scene *scene.Type

	background *ui.Image
}

func newRespawnScreen() *respawnScreen {
	rs := &respawnScreen{
		scene: scene.New(true),
	}

	rs.background = ui.NewImage(render.GetTexture("solid"), 0, 0, 854, 480, 0, 0, 1, 1, 0, 0, 0)
	rs.background.SetA(160)
	rs.scene.AddDrawable(rs.background.Attach(ui.Top, ui.Left))

	rs.scene.AddDrawable(
		ui.NewText("Respawn:", 0, -20, 255, 255, 255).Attach(ui.Center, ui.Middle),
	)

	respawn, txt := newButtonText("Respawn", -205, 20, 400, 40)
	rs.scene.AddDrawable(respawn.Attach(ui.Center, ui.Middle))
	rs.scene.AddDrawable(txt)
	respawn.AddClick(func() {
		setScreen(nil)
		Client.network.Write(&protocol.ClientStatus{ActionID: 0})
	})

	disconnect, txt := newButtonText("Disconnect", 205, 20, 400, 40)
	rs.scene.AddDrawable(disconnect.Attach(ui.Center, ui.Middle))
	rs.scene.AddDrawable(txt)
	disconnect.AddClick(func() { Client.network.SignalClose(errManualDisconnect) })

	uiFooter(rs.scene)
	return rs
}

func (rs *respawnScreen) tick(delta float64) {
	width, height := window.GetFramebufferSize()
	rs.background.SetWidth(float64(width) / ui.Scale)
	rs.background.SetHeight(float64(height) / ui.Scale)
}

func (rs *respawnScreen) remove() {
	rs.scene.Hide()
}
