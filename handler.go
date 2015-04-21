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
	"fmt"
	"math"
	"reflect"

	"github.com/thinkofdeath/steven/protocol"
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
	sendPluginMessage(&pmMinecraftBrand{
		Brand: "Steven",
	})
	Client.GameMode = gameMode(j.Gamemode & 0x7)
	Client.HardCore = j.Gamemode&0x8 != 0
}

func (handler) Respawn(r *protocol.Respawn) {
	for _, c := range chunkMap {
		c.free()
	}
	chunkMap = map[chunkPosition]*chunk{}
	Client.GameMode = gameMode(r.Gamemode & 0x7)
	Client.HardCore = r.Gamemode&0x8 != 0
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
func (h handler) PluginMessage(p *protocol.PluginMessageClientbound) {
	h.handlePluginMessage(p.Channel, bytes.NewReader(p.Data), false)
}

func (h handler) ServerBrand(b *pmMinecraftBrand) {
	fmt.Printf("The server is running: %s\n", b.Brand)
}
