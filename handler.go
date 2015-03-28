package main

import (
	"fmt"
	"math"
	"reflect"

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

	for i := 0; i < v.NumMethod(); i++ {
		m := v.Method(i)
		t := m.Type()
		if t.NumIn() != 1 && t.Name() != "Handle" {
			continue
		}
		in := t.In(0)
		if in.AssignableTo(packet) {
			h[in] = m
		}
	}
}

func (h handler) Handle(packet protocol.Packet) {
	m, ok := h[reflect.TypeOf(packet)]
	if ok {
		m.Call([]reflect.Value{reflect.ValueOf(packet)})
	}
}

func (handler) ServerMessage(msg *protocol.ServerMessage) {
	fmt.Printf("MSG(%d): %s\n", msg.Type, msg.Message.Value)
}

func (handler) Respawn(c *protocol.Respawn) {
	for _, c := range chunkMap {
		c.free()
	}
	chunkMap = map[chunkPosition]*chunk{}
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

func (handler) Teleport(t *protocol.TeleportPlayer) {
	render.Camera.X = t.X
	render.Camera.Y = t.Y
	render.Camera.Z = t.Z
	render.Camera.Yaw = float64(-t.Yaw) * (math.Pi / 180)
	render.Camera.Pitch = -float64(t.Pitch)*(math.Pi/180) + math.Pi
	writeChan <- &protocol.PlayerPositionLook{
		X:     t.X,
		Y:     t.Y,
		Z:     t.Z,
		Yaw:   t.Yaw,
		Pitch: t.Pitch,
	}
	ready = true
}
