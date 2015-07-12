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
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"math"
	"strings"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/format"
	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/render/ui"
	"github.com/thinkofdeath/steven/render/ui/scene"
)

var (
	disconnectReason    format.AnyComponent
	errManualDisconnect = errors.New("manual disconnect")
)

type serverList struct {
	baseUI
	scene *scene.Type
	logo  uiLogo

	servers []*serverListItem
}

type serverListItem struct {
	*scene.Type

	X, Y      float64
	container *ui.Container
	offset    float64
	id        string
}

func newServerList() screen {
	sl := &serverList{
		scene: scene.New(true),
	}
	sl.logo.init(sl.scene)

	uiFooter(sl.scene)

	sl.redraw()

	refresh, txt := newButtonText("Refresh", 300, -50-15, 100, 30)
	sl.scene.AddDrawable(refresh.Attach(ui.Center, ui.Middle))
	sl.scene.AddDrawable(txt)
	refresh.AddClick(sl.redraw)

	add, txt := newButtonText("Add", 200, -50-15, 100, 30)
	sl.scene.AddDrawable(add.Attach(ui.Center, ui.Middle))
	sl.scene.AddDrawable(txt)
	add.AddClick(func() {
		setScreen(newEditServer(-1))
	})

	options := ui.NewButton(5, 25, 40, 40)
	sl.scene.AddDrawable(options.Attach(ui.Bottom, ui.Right))
	cog := ui.NewImage(render.GetTexture("steven:gui/cog"), 0, 0, 40, 40, 0, 0, 1, 1, 255, 255, 255)
	cog.AttachTo(options)
	sl.scene.AddDrawable(cog.Attach(ui.Center, ui.Middle))
	options.AddClick(func() {
		setScreen(newOptionMenu(newServerList))
	})

	if disconnectReason.Value != nil {
		disMsg := ui.NewText("Disconnected", 0, 32, 255, 0, 0).Attach(ui.Top, ui.Center)
		dis := ui.NewFormattedWidth(disconnectReason, 0, 48, 600)
		disB := ui.NewImage(render.GetTexture("solid"), 0, 30, math.Max(dis.Width, disMsg.Width)+4, dis.Height+4+16, 0, 0, 1, 1, 0, 0, 0)
		disB.SetA(100)
		sl.scene.AddDrawable(disB.Attach(ui.Top, ui.Center))
		sl.scene.AddDrawable(dis.Attach(ui.Top, ui.Center))
		sl.scene.AddDrawable(disMsg)
	}

	return sl
}

func (sl *serverList) init() {
	window.SetScrollCallback(sl.onScroll)
	sl.playMusic()
}

func (sl *serverList) playMusic() {
	if currentScreen != sl {
		return
	}
	StopAllMusic()
	PlaySoundCallback("music.menu", sl.playMusic)
}

func (sl *serverList) onScroll(w *glfw.Window, xoff float64, yoff float64) {
	if len(sl.servers) == 0 {
		return
	}
	diff := yoff / 10
	if s := sl.servers[len(sl.servers)-1]; s.offset+diff <= 2 {
		diff = 2 - s.offset
	}
	if s := sl.servers[0]; s.offset+diff >= 0 {
		diff = -s.offset
	}
	for _, s := range sl.servers {
		s.offset += diff
		s.updatePosition()
	}
}

func (si *serverListItem) updatePosition() {
	if si.offset < 0 {
		si.Y = si.offset * 200
	} else if si.offset >= 2 {
		si.Y = si.offset * 100
	} else {
		si.Y = si.offset * 100
	}
}

func (sl *serverList) redraw() {
	for _, s := range sl.servers {
		s.Hide()
		render.FreeIcon(s.id)
	}
	sl.servers = sl.servers[:0]
	for i, s := range Config.Servers {
		sc := scene.New(true)
		container := ui.NewContainer(0, float64(i)*100, 700, 100).
			Attach(ui.Center, ui.Middle)
		r := make([]byte, 20)
		rand.Read(r)
		si := &serverListItem{
			Type:      sc,
			container: container,
			offset:    float64(i),
			id:        "servericon:" + string(r),
		}
		si.updatePosition()
		sl.servers = append(sl.servers, si)

		bck := ui.NewImage(render.GetTexture("solid"), 0, 0, 700, 100, 0, 0, 1, 1, 0, 0, 0).Attach(ui.Top, ui.Left)
		bck.SetA(100)
		bck.AttachTo(container)
		sc.AddDrawable(bck)
		txt := ui.NewText(s.Name, 90+10, 5, 255, 255, 255).Attach(ui.Top, ui.Left)
		txt.AttachTo(container)
		sc.AddDrawable(txt)

		icon := ui.NewImage(render.GetTexture("misc/unknown_server"), 5, 5, 90, 90, 0, 0, 1, 1, 255, 255, 255).
			Attach(ui.Top, ui.Left)
		icon.AttachTo(container)
		sc.AddDrawable(icon)

		ping := ui.NewImage(render.GetTexture("gui/icons"), 5, 5, 20, 16, 0, 56/256.0, 10/256.0, 8/256.0, 255, 255, 255).
			Attach(ui.Top, ui.Right)
		ping.AttachTo(container)
		sc.AddDrawable(ping)

		players := ui.NewText("???", 30, 5, 255, 255, 255).
			Attach(ui.Top, ui.Right)
		players.AttachTo(container)
		sc.AddDrawable(players)

		msg := &format.TextComponent{Text: "Connecting..."}
		motd := ui.NewFormattedWidth(format.Wrap(msg), 90+10, 5+18, 700-(90+10+5)).Attach(ui.Top, ui.Left)
		motd.AttachTo(container)
		sc.AddDrawable(motd)
		s := s
		go sl.pingServer(s.Address, motd, icon, si.id, ping, players)
		container.ClickFunc = func() {
			PlaySound("random.click")
			sl.connect(s.Address)
		}
		container.HoverFunc = func(over bool) {
			if over {
				bck.SetA(200)
			} else {
				bck.SetA(100)
			}
		}

		sc.AddDrawable(container)

		index := i
		del, txt := newButtonText("X", 0, 0, 25, 25)
		del.AttachTo(container)
		sc.AddDrawable(del.Attach(ui.Bottom, ui.Right))
		sc.AddDrawable(txt)
		del.AddClick(func() {
			Config.Servers = append(Config.Servers[:index], Config.Servers[index+1:]...)
			saveServers()
			sl.redraw()
		})
		edit, txt := newButtonText("E", 25, 0, 25, 25)
		edit.AttachTo(container)
		sc.AddDrawable(edit.Attach(ui.Bottom, ui.Right))
		sc.AddDrawable(txt)
		edit.AddClick(func() {
			setScreen(newEditServer(index))
		})
	}
}

func (sl *serverList) pingServer(addr string, motd *ui.Formatted,
	icon *ui.Image, id string, ping *ui.Image, players *ui.Text) {
	conn, err := protocol.Dial(addr)
	if err != nil {
		syncChan <- func() {
			msg := &format.TextComponent{Text: err.Error()}
			msg.Color = format.Red
			motd.Update(format.Wrap(msg))
		}
		return
	}
	defer conn.Close()
	resp, pingTime, err := conn.RequestStatus()
	syncChan <- func() {
		if err != nil {
			msg := &format.TextComponent{Text: err.Error()}
			msg.Color = format.Red
			motd.Update(format.Wrap(msg))
			return
		}
		y := 0.0
		pt := pingTime.Seconds() / 1000
		switch {
		case pt <= 75:
			y = 16 / 256.0
		case pt <= 150:
			y = 24 / 256.0
		case pt <= 225:
			y = 32 / 256.0
		case pt <= 350:
			y = 40 / 256.0
		case pt < 999:
			y = 48 / 256.0
		default:
			y = 56 / 256.0
		}
		ping.SetTextureY(y)

		players.Update(fmt.Sprintf("%d/%d", resp.Players.Online, resp.Players.Max))

		desc := resp.Description
		format.ConvertLegacy(desc)
		motd.Update(desc)

		if strings.HasPrefix(resp.Favicon, "data:image/png;base64,") {
			favicon := resp.Favicon[len("data:image/png;base64,"):]
			data, err := base64.StdEncoding.DecodeString(favicon)
			if err != nil {
				fmt.Printf("error base64 decoding favicon: %s\n", err)
				return
			}
			img, err := png.Decode(bytes.NewReader(data))
			if err != nil {
				fmt.Printf("error decoding favicon: %s\n", err)
				return
			}
			render.AddIcon(id, img)
			icon.SetTexture(render.Icon(id))
		}
	}
}

func (sl *serverList) connect(s string) {
	connect(s)
}

func (sl *serverList) tick(delta float64) {
	sl.logo.tick(delta)
	for _, s := range sl.servers {
		dx := s.X - s.container.X()
		dy := s.Y - s.container.Y()
		if dx*dx > 1 {
			s.container.SetX(s.container.X() + delta*dx*0.1)
		} else {
			s.container.SetX(s.X)
		}
		if dy*dy > 1 {
			s.container.SetY(s.container.Y() + delta*dy*0.1)
		} else {
			s.container.SetY(s.Y)
		}
	}
}

func (sl *serverList) remove() {
	window.SetScrollCallback(onScroll)
	sl.scene.Hide()
	for _, s := range sl.servers {
		s.Hide()
		render.FreeIcon(s.id)
	}
}
