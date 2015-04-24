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
	"github.com/thinkofdeath/steven/render/ui"
)

var playerList = map[protocol.UUID]*playerInfo{}

type playerInfo struct {
	name        string
	uuid        protocol.UUID
	displayName chat.AnyComponent
	gameMode    gameMode
	ping        int
}

type playerListUI struct {
	enabled bool

	background *ui.Image
	elements   []*ui.Text
	icons      []*ui.Image
}

func (p *playerListUI) init() {
	p.background = ui.NewImage(render.GetTexture("solid"), 0, 16, 16, 16, 0, 0, 1, 1, 0, 0, 0)
	p.background.A = 120
	p.background.Visible = false
	ui.AddDrawable(p.background, ui.Top, ui.Center)
}

func (p *playerListUI) set(enabled bool) {
	p.enabled = enabled
	p.background.Visible = enabled
	for _, e := range p.elements {
		e.Visible = enabled
	}
	for _, i := range p.icons {
		i.Visible = false
	}
}

func (p *playerListUI) render(delta float64) {
	if !p.enabled {
		return
	}
	for _, e := range p.elements {
		e.Visible = false
	}
	for _, i := range p.icons {
		i.Visible = false
	}
	offset := 0
	count := 0
	width := 0.0
	for i, pl := range p.players() {
		if offset >= len(p.elements) {
			text := ui.NewText("", 8, 0, 255, 255, 255)
			p.elements = append(p.elements, text)
			ui.AddDrawable(text, ui.Top, ui.Center)
			icon := ui.NewImage(render.GetTexture("entity/steve"), 0, 0, 16, 16, 8/64.0, 8/64.0, 8/64.0, 8/64.0, 255, 255, 255)
			p.icons = append(p.icons, icon)
			ui.AddDrawable(icon, ui.Top, ui.Center)
			text.Parent = p.background
			icon.Parent = p.background
		}
		text := p.elements[offset]
		icon := p.icons[offset]
		offset++
		text.Visible = true
		text.Y = 1 + 18*float64(i)
		text.Update(pl.name)
		count++
		if text.Width > width {
			width = text.Width
		}
		icon.Visible = true
		icon.Y = 1 + 18*float64(i)
	}
	for _, i := range p.icons {
		if i.Visible {
			i.X = -width/2 - 4
		}
	}

	p.background.W = width + 32
	p.background.H = float64(count * 18)
}

func (p *playerListUI) players() (out []*playerInfo) {
	for _, pl := range playerList {
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

func (handler) PlayerListInfo(p *protocol.PlayerInfo) {
	for _, pl := range p.Players {
		if _, ok := playerList[pl.UUID]; !ok && p.Action != 0 {
			continue
		}
		switch p.Action {
		case 0: // Add
			playerList[pl.UUID] = &playerInfo{
				name:        pl.Name,
				uuid:        pl.UUID,
				displayName: pl.DisplayName,
				gameMode:    gameMode(pl.GameMode),
				ping:        int(pl.Ping),
			}
		case 1: // Update gamemode
			playerList[pl.UUID].gameMode = gameMode(pl.GameMode)
		case 2: // Update ping
			playerList[pl.UUID].ping = int(pl.Ping)
		case 3: // Update display name
			playerList[pl.UUID].displayName = pl.DisplayName
		case 4: // Remove
			delete(playerList, pl.UUID)
		}
	}
}
