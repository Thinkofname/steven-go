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
	"encoding/base64"
	"encoding/json"
	"sort"
	"strings"

	"github.com/thinkofdeath/steven/chat"
	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/render/ui"
)

const playerListWidth = 150

var playerList = map[protocol.UUID]*playerInfo{}

type playerInfo struct {
	name        string
	uuid        protocol.UUID
	displayName chat.AnyComponent
	gameMode    gameMode
	ping        int

	skin     *render.TextureInfo
	skinHash string
}

type playerListUI struct {
	enabled bool

	background [4]*ui.Image
	entries    []*playerListUIEntry
}

type playerListUIEntry struct {
	text    *ui.Text
	icon    *ui.Image
	iconHat *ui.Image
	ping    *ui.Image
}

func (p playerListUIEntry) set(enabled bool) {
	p.text.Visible = enabled
	p.icon.Visible = enabled
	p.iconHat.Visible = enabled
	p.ping.Visible = enabled
}

func (p *playerListUI) init() {
	for i := range p.background {
		p.background[i] = ui.NewImage(render.GetTexture("solid"), 0, 16, playerListWidth+48, 16, 0, 0, 1, 1, 0, 0, 0)
		p.background[i].A = 120
		p.background[i].Visible = false
		ui.AddDrawable(p.background[i].Attach(ui.Top, ui.Center))
	}
}

func (p *playerListUI) set(enabled bool) {
	p.enabled = enabled
	for _, b := range p.background {
		b.Visible = enabled
	}
	for _, e := range p.entries {
		e.set(enabled)
	}
}

func (p *playerListUI) render(delta float64) {
	if !p.enabled {
		return
	}
	for _, b := range p.background {
		b.Visible = false
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
				if e.icon.Visible {
					e.icon.X = -playerListWidth/2 - 12
					e.iconHat.X = -playerListWidth/2 - 12
					e.ping.X = playerListWidth/2 + 12
				}
			}
			p.background[bTab].H = float64(count * 18)
			count = 0
			bTab++
			if bTab >= len(p.background) {
				break
			}
		}
		background := p.background[bTab]
		background.Visible = true
		if offset >= len(p.entries) {
			text := ui.NewText("", 24, 0, 255, 255, 255).
				Attach(ui.Top, ui.Left)
			ui.AddDrawable(text)
			icon := ui.NewImage(pl.skin, 0, 0, 16, 16, 8/64.0, 8/64.0, 8/64.0, 8/64.0, 255, 255, 255).
				Attach(ui.Top, ui.Center)
			ui.AddDrawable(icon)
			iconHat := ui.NewImage(pl.skin, 0, 0, 16, 16, 40/64.0, 8/64.0, 8/64.0, 8/64.0, 255, 255, 255).
				Attach(ui.Top, ui.Center)
			ui.AddDrawable(iconHat)
			ping := ui.NewImage(render.GetTexture("gui/icons"), 0, 0, 20, 16, 0, 16/256.0, 10/256.0, 8/256.0, 255, 255, 255).
				Attach(ui.Top, ui.Center)
			ui.AddDrawable(ping)

			text.Parent = background
			icon.Parent = background
			iconHat.Parent = background
			ping.Parent = background

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
		e.text.Y = 1 + 18*float64(count)
		e.text.Update(pl.name)
		e.icon.Y = 1 + 18*float64(count)
		e.icon.Texture = pl.skin
		e.iconHat.Y = 1 + 18*float64(count)
		e.iconHat.Texture = pl.skin

		e.ping.Y = 1 + 18*float64(count)
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
		e.ping.TY = y
		count++
	}

	if bTab < len(p.background) {
		for _, e := range p.entries {
			if e.icon.Visible {
				e.icon.X = -playerListWidth/2 - 12
				e.iconHat.X = -playerListWidth/2 - 12
				e.ping.X = playerListWidth/2 + 12
			}
		}
		p.background[bTab].H = float64(count * 18)
	}

	switch bTab {
	case 0: // Single
		p.background[0].X = 0
	case 1: // Double
		p.background[0].X = -p.background[0].W / 2
		p.background[1].X = p.background[1].W / 2
	case 2: // Triple
		p.background[0].X = -(p.background[1].W / 2) - p.background[0].W/2
		p.background[1].X = 0
		p.background[2].X = (p.background[1].W / 2) + p.background[2].W/2
	default: // Quad
		p.background[0].X = -p.background[0].W/2 - p.background[1].W
		p.background[3].X = p.background[3].W/2 + p.background[2].W

		p.background[1].X = -p.background[1].W / 2
		p.background[2].X = p.background[2].W / 2
	}
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
		if _, ok := playerList[pl.UUID]; (!ok && p.Action != 0) || (ok && p.Action == 0) {
			continue
		}
		switch p.Action {
		case 0: // Add
			i := &playerInfo{
				name:        pl.Name,
				uuid:        pl.UUID,
				displayName: pl.DisplayName,
				gameMode:    gameMode(pl.GameMode),
				ping:        int(pl.Ping),
			}
			for _, prop := range pl.Properties {
				if prop.Name == "textures" {
					data, err := base64.URLEncoding.DecodeString(prop.Value)
					if err != nil {
						continue
					}
					var blob skinBlob
					err = json.Unmarshal(data, &blob)
					if err != nil {
						continue
					}
					url := blob.Textures.Skin.Url
					if strings.HasPrefix(url, "http://textures.minecraft.net/texture/") {
						i.skinHash = url[len("http://textures.minecraft.net/texture/"):]
						render.RefSkin(i.skinHash)
						i.skin = render.Skin(i.skinHash)
					}
				}
			}
			if i.skin == nil {
				i.skin = render.GetTexture("entity/steve")
			}
			playerList[pl.UUID] = i
		case 1: // Update gamemode
			playerList[pl.UUID].gameMode = gameMode(pl.GameMode)
		case 2: // Update ping
			playerList[pl.UUID].ping = int(pl.Ping)
		case 3: // Update display name
			playerList[pl.UUID].displayName = pl.DisplayName
		case 4: // Remove
			i := playerList[pl.UUID]
			if i.skinHash != "" {
				render.FreeSkin(i.skinHash)
			}
			delete(playerList, pl.UUID)
		}
	}
}

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
