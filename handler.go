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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/render"
)

type handler map[reflect.Type]reflect.Value

var defaultHandler = handler{}

func init() {
	defaultHandler.Init()
}

func (h handler) Init() {
	v := reflect.ValueOf(h)

	packet := reflect.TypeOf((*protocol.Packet)(nil)).Elem()
	pm := reflect.TypeOf((*pluginMessage)(nil)).Elem()

	for i := 0; i < v.NumMethod(); i++ {
		m := v.Method(i)
		t := m.Type()
		if t.NumIn() != 1 && t.Name() != "Handle" {
			continue
		}
		in := t.In(0)
		if in.AssignableTo(packet) || in.AssignableTo(pm) {
			h[in] = m
		}
	}
}

func (h handler) Handle(packet interface{}) {
	m, ok := h[reflect.TypeOf(packet)]
	if ok {
		m.Call([]reflect.Value{reflect.ValueOf(packet)})
	}
}

func (handler) ServerMessage(msg *protocol.ServerMessage) {
	fmt.Printf("MSG(%d): %s\n", msg.Type, msg.Message.Value)
	Client.chat.Add(msg.Message)
}

func (handler) JoinGame(j *protocol.JoinGame) {
	clearChunks()
	sendPluginMessage(&pmMinecraftBrand{
		Brand: "Steven",
	})
	Client.GameMode = gameMode(j.Gamemode & 0x7)
	Client.HardCore = j.Gamemode&0x8 != 0
}

func (handler) Respawn(r *protocol.Respawn) {
	clearChunks()
	Client.GameMode = gameMode(r.Gamemode & 0x7)
	Client.HardCore = r.Gamemode&0x8 != 0
}

func (handler) Disconnect(d *protocol.Disconnect) {
	disconnectReason = d.Reason
	fmt.Println("Disconnect: ", disconnectReason)
	closeWithError(errManualDisconnect)
}

func (handler) UpdateHealth(u *protocol.UpdateHealth) {
	Client.UpdateHealth(float64(u.Health))
	Client.UpdateHunger(float64(u.Food))
}

func (handler) ChangeGameState(c *protocol.ChangeGameState) {
	switch c.Reason {
	case 3: // Change game mode
		Client.GameMode = gameMode(c.Value)
	}
}

func (handler) ChangeHotbarSlot(s *protocol.SetCurrentHotbarSlot) {
	Client.currentHotbarSlot = int(s.Slot)
}

func (handler) Teleport(t *protocol.TeleportPlayer) {
	Client.X = calculateTeleport(teleportRelX, t.Flags, Client.X, t.X)
	Client.Y = calculateTeleport(teleportRelY, t.Flags, Client.Y, t.Y)
	Client.Z = calculateTeleport(teleportRelZ, t.Flags, Client.Z, t.Z)
	Client.Yaw = calculateTeleport(teleportRelYaw, t.Flags, Client.Yaw, float64(-t.Yaw)*(math.Pi/180))
	Client.Pitch = calculateTeleport(teleportRelPitch, t.Flags, Client.Pitch, -float64(t.Pitch)*(math.Pi/180)+math.Pi)
	Client.checkGround()
	writeChan <- &protocol.PlayerPositionLook{
		X:        t.X,
		Y:        t.Y,
		Z:        t.Z,
		Yaw:      t.Yaw,
		Pitch:    t.Pitch,
		OnGround: Client.OnGround,
	}
	Client.copyToCamera()
	ready = true
}

func (handler) ChunkData(c *protocol.ChunkData) {
	if c.BitMask == 0 && c.New {
		pos := chunkPosition{int(c.ChunkX), int(c.ChunkZ)}
		c, ok := chunkMap[pos]
		if ok {
			c.free()
			delete(chunkMap, pos)
		}
		return
	}
	go loadChunk(int(c.ChunkX), int(c.ChunkZ), c.Data, c.BitMask, true, c.New)
}

func (handler) ChunkDataBulk(c *protocol.ChunkDataBulk) {
	go func() {
		offset := 0
		data := c.Data
		for _, meta := range c.Meta {
			offset += loadChunk(int(meta.ChunkX), int(meta.ChunkZ), data[offset:], meta.BitMask, c.SkyLight, true)
		}
	}()
}

func (handler) SetBlock(b *protocol.BlockChange) {
	block := GetBlockByCombinedID(uint16(b.BlockID))
	chunkMap.SetBlock(block, b.Location.X(), b.Location.Y(), b.Location.Z())
	chunkMap.UpdateBlock(b.Location.X(), b.Location.Y(), b.Location.Z())
}

func (handler) SetBlockBatch(b *protocol.MultiBlockChange) {
	chunk := chunkMap[chunkPosition{int(b.ChunkX), int(b.ChunkZ)}]
	if chunk == nil {
		return
	}
	for _, r := range b.Records {
		block := GetBlockByCombinedID(uint16(r.BlockID))
		x, y, z := int(r.XZ>>4), int(r.Y), int(r.XZ&0xF)
		chunk.setBlock(block, x, y, z)
		chunkMap.UpdateBlock((chunk.X<<4)+x, y, (chunk.Z<<4)+z)
	}
}

func (handler) SpawnPlayer(s *protocol.SpawnPlayer) {
	e := newPlayer()
	if p, ok := e.(PositionComponent); ok {
		p.SetPosition(
			float64(s.X)/32,
			float64(s.Y)/32,
			float64(s.Z)/32,
		)
	}
	if p, ok := e.(TargetPositionComponent); ok {
		p.SetTargetPosition(
			float64(s.X)/32,
			float64(s.Y)/32,
			float64(s.Z)/32,
		)
	}
	Client.entities.add(int(s.EntityID), e)
}

func (handler) SpawnMob(s *protocol.SpawnMob) {
	et, ok := entityTypes[int(s.Type)]
	if !ok {
		return
	}
	e := et()
	if p, ok := e.(PositionComponent); ok {
		p.SetPosition(
			float64(s.X)/32,
			float64(s.Y)/32,
			float64(s.Z)/32,
		)
	}
	if p, ok := e.(TargetPositionComponent); ok {
		p.SetTargetPosition(
			float64(s.X)/32,
			float64(s.Y)/32,
			float64(s.Z)/32,
		)
	}
	Client.entities.add(int(s.EntityID), e)
}

func (handler) EntityTeleport(t *protocol.EntityTeleport) {
	e, ok := Client.entities.entities[int(t.EntityID)]
	if !ok {
		return
	}
	if p, ok := e.(PositionComponent); ok {
		p.SetPosition(
			float64(t.X)/32,
			float64(t.Y)/32,
			float64(t.Z)/32,
		)
	}
	if p, ok := e.(TargetPositionComponent); ok {
		p.SetTargetPosition(
			float64(t.X)/32,
			float64(t.Y)/32,
			float64(t.Z)/32,
		)
	}
}

func (handler) EntityMove(m *protocol.EntityMove) {
	e, ok := Client.entities.entities[int(m.EntityID)]
	if !ok {
		return
	}
	dx, dy, dz := float64(m.DeltaX)/32, float64(m.DeltaY)/32, float64(m.DeltaZ)/32
	relMove(e, dx, dy, dz)
}

func (handler) EntityMoveLook(m *protocol.EntityLookAndMove) {
	e, ok := Client.entities.entities[int(m.EntityID)]
	if !ok {
		return
	}
	dx, dy, dz := float64(m.DeltaX)/32, float64(m.DeltaY)/32, float64(m.DeltaZ)/32
	relMove(e, dx, dy, dz)
}

func relMove(e Entity, dx, dy, dz float64) {
	if p, ok := e.(TargetPositionComponent); ok {
		x, y, z := p.TargetPosition()
		p.SetTargetPosition(
			x+dx,
			y+dy,
			z+dz,
		)
		return
	}
	if p, ok := e.(PositionComponent); ok {
		x, y, z := p.Position()
		p.SetPosition(
			x+dx,
			y+dy,
			z+dz,
		)
	}
}

func (handler) DestoryEntities(e *protocol.EntityDestroy) {
	for _, id := range e.EntityIDs {
		Client.entities.remove(int(id))
	}
}

func (handler) PlayerListInfo(p *protocol.PlayerInfo) {
	playerList := Client.playerList.info
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
					if !prop.IsSigned {
						closeWithError(errors.New("Missing signature from textures"))
						return
					}
					data, err := base64.StdEncoding.DecodeString(prop.Value)
					if err != nil {
						closeWithError(err)
						continue
					}

					sig, err := base64.StdEncoding.DecodeString(prop.Signature)
					if err != nil {
						closeWithError(err)
						continue
					}

					if err := verifySkinSignature([]byte(prop.Value), sig); err != nil {
						closeWithError(err)
						return
					}

					var blob skinBlob
					err = json.Unmarshal(data, &blob)
					if err != nil {
						closeWithError(err)
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

func (h handler) PluginMessage(p *protocol.PluginMessageClientbound) {
	h.handlePluginMessage(p.Channel, bytes.NewReader(p.Data), false)
}

func (h handler) ServerBrand(b *pmMinecraftBrand) {
	fmt.Printf("The server is running: %s\n", b.Brand)
}
