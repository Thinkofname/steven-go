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
	"sort"

	"github.com/thinkofdeath/steven/chat"
	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/ui"
	"github.com/thinkofdeath/steven/ui/scene"
)

const playerListWidth = 150

type playerInfo struct {
	name        string
	uuid        protocol.UUID
	displayName chat.AnyComponent
	gameMode    gameMode
	ping        int

	skin     render.TextureInfo
	skinHash string
}

type playerListUI struct {
	info       map[protocol.UUID]*playerInfo
	background [4]*ui.Image
	entries    []*playerListUIEntry
	scene      *scene.Type
}

type playerListUIEntry struct {
	text    *ui.Text
	icon    *ui.Image
	iconHat *ui.Image
	ping    *ui.Image
}

func (p playerListUIEntry) set(enabled bool) {
	p.text.SetDraw(enabled)
	p.icon.SetDraw(enabled)
	p.iconHat.SetDraw(enabled)
	p.ping.SetDraw(enabled)
}

func (p *playerListUI) init() {
	p.info = map[protocol.UUID]*playerInfo{}
	p.scene = scene.New(false)
	for i := range p.background {
		p.background[i] = ui.NewImage(render.GetTexture("solid"), 0, 16, playerListWidth+48, 16, 0, 0, 1, 1, 0, 0, 0)
		p.background[i].SetA(120)
		p.background[i].SetDraw(false)
		p.scene.AddDrawable(p.background[i].Attach(ui.Top, ui.Center))
	}
}

func (p *playerListUI) free() {
	for _, pl := range p.info {
		if pl.skinHash != "" {
			render.FreeSkin(pl.skinHash)
		}
	}
}
func (p *playerListUI) set(enabled bool) {
	if enabled {
		p.scene.Show()
	} else {
		p.scene.Hide()
	}
}

func (p *playerListUI) render(delta float64) {
	if !p.scene.IsVisible() {
		return
	}
	for _, b := range p.background {
		b.SetDraw(false)
	}
	for _, e := range p.entries {
		e.set(false)
	}
	offset := 0
	count := 0
	bTab := 0
	lastEntry := 0
	for _, pl := range p.players() {
		if count >= 20 {
			entries := p.entries[lastEntry:offset]
			lastEntry = offset
			for _, e := range entries {
				if e.icon.ShouldDraw() {
					e.icon.SetX(-playerListWidth/2 - 12)
					e.iconHat.SetX(-playerListWidth/2 - 12)
					e.ping.SetX(playerListWidth/2 + 12)
				}
			}
			p.background[bTab].SetHeight(float64(count * 18))
			count = 0
			bTab++
			if bTab >= len(p.background) {
				break
			}
		}
		background := p.background[bTab]
		background.SetDraw(true)
		if offset >= len(p.entries) {
			text := ui.NewText("", 24, 0, 255, 255, 255).
				Attach(ui.Top, ui.Left)
			p.scene.AddDrawable(text)
			icon := ui.NewImage(pl.skin, 0, 0, 16, 16, 8/64.0, 8/64.0, 8/64.0, 8/64.0, 255, 255, 255).
				Attach(ui.Top, ui.Center)
			p.scene.AddDrawable(icon)
			iconHat := ui.NewImage(pl.skin, 0, 0, 16, 16, 40/64.0, 8/64.0, 8/64.0, 8/64.0, 255, 255, 255).
				Attach(ui.Top, ui.Center)
			p.scene.AddDrawable(iconHat)
			ping := ui.NewImage(render.GetTexture("gui/icons"), 0, 0, 20, 16, 0, 16/256.0, 10/256.0, 8/256.0, 255, 255, 255).
				Attach(ui.Top, ui.Center)
			p.scene.AddDrawable(ping)

			text.AttachTo(background)
			icon.AttachTo(background)
			iconHat.AttachTo(background)
			ping.AttachTo(background)

			p.entries = append(p.entries, &playerListUIEntry{
				text:    text,
				icon:    icon,
				iconHat: iconHat,
				ping:    ping,
			})
		}
		e := p.entries[offset]
		e.set(true)
		offset++
		e.text.SetY(1 + 18*float64(count))
		e.text.Update(pl.name)
		e.icon.SetY(1 + 18*float64(count))
		e.icon.SetTexture(pl.skin)
		e.iconHat.SetY(1 + 18*float64(count))
		e.iconHat.SetTexture(pl.skin)

		e.ping.SetY(1 + 18*float64(count))
		y := 0.0
		switch {
		case pl.ping <= 75:
			y = 16 / 256.0
		case pl.ping <= 150:
			y = 24 / 256.0
		case pl.ping <= 225:
			y = 32 / 256.0
		case pl.ping <= 350:
			y = 40 / 256.0
		case pl.ping < 999:
			y = 48 / 256.0
		default:
			y = 56 / 256.0
		}
		e.ping.SetTextureY(y)
		count++
	}

	if bTab < len(p.background) {
		for _, e := range p.entries {
			if e.icon.ShouldDraw() {
				e.icon.SetX(-playerListWidth/2 - 12)
				e.iconHat.SetX(-playerListWidth/2 - 12)
				e.ping.SetX(playerListWidth/2 + 12)
			}
		}
		p.background[bTab].SetHeight(float64(count * 18))
	}

	switch bTab {
	case 0: // Single
		p.background[0].SetX(0)
	case 1: // Double
		p.background[0].SetX(-p.background[0].Width() / 2)
		p.background[1].SetX(p.background[1].Width() / 2)
	case 2: // Triple
		p.background[0].SetX(-(p.background[1].Width() / 2) - p.background[0].Width()/2)
		p.background[1].SetX(0)
		p.background[2].SetX((p.background[1].Width() / 2) + p.background[2].Width()/2)
	default: // Quad
		p.background[0].SetX(-p.background[0].Width()/2 - p.background[1].Width())
		p.background[3].SetX(p.background[3].Width()/2 + p.background[2].Width())

		p.background[1].SetX(-p.background[1].Width() / 2)
		p.background[2].SetX(p.background[2].Width() / 2)
	}
}

func (p *playerListUI) players() (out []*playerInfo) {
	for _, pl := range p.info {
		out = append(out, pl)
	}
	sort.Sort(sortedPlayerList(out))
	return
}

type sortedPlayerList []*playerInfo

func (s sortedPlayerList) Len() int { return len(s) }
func (s sortedPlayerList) Less(a, b int) bool {
	if s[a].name < s[b].name {
		return true
	}
	return false
}
func (s sortedPlayerList) Swap(a, b int) { s[a], s[b] = s[b], s[a] }

type skinBlob struct {
	Timestamp     int64
	ProfileID     string
	ProfileString string
	IsPublic      bool
	Textures      struct {
		Skin skinPath
		Cape skinPath
	}
}

type skinPath struct {
	Url string
}
