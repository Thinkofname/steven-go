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

package main

import (
	"encoding/binary"
	"sort"

	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/type/direction"
	"github.com/thinkofdeath/steven/type/nibble"
	"github.com/thinkofdeath/steven/world/biome"
)

var chunkMap world = map[chunkPosition]*chunk{}

type world map[chunkPosition]*chunk

func (w world) Block(x, y, z int) Block {
	cx := x >> 4
	cz := z >> 4
	chunk := w[chunkPosition{cx, cz}]
	if chunk == nil {
		return BlockBedrock.Base
	}
	return chunk.block(x&0xF, y, z&0xF)
}

func (w world) SetBlock(b Block, x, y, z int) {
	cx := x >> 4
	cz := z >> 4
	chunk := w[chunkPosition{cx, cz}]
	if chunk == nil {
		return
	}
	chunk.setBlock(b, x&0xF, y, z&0xF)
	for _, d := range direction.Values {
		ox, oy, oz := d.Offset()
		w.dirty(x+ox, y+oy, z+oz)
	}
}

func (w world) dirty(x, y, z int) {
	cx := x >> 4
	cz := z >> 4
	chunk := w[chunkPosition{cx, cz}]
	if chunk == nil || y < 0 || y > 255 {
		return
	}
	cs := chunk.Sections[y>>4]
	if cs == nil {
		return
	}
	cs.dirty = true
}

func (w world) UpdateBlock(x, y, z int) {
	for yy := -1; yy <= 1; yy++ {
		for zz := -1; zz <= 1; zz++ {
			for xx := -1; xx <= 1; xx++ {
				bx, by, bz := x+xx, y+yy, z+zz
				b := w.Block(bx, by, bz)
				nb := b.UpdateState(bx, by, bz)
				if b != nb {
					w.SetBlock(nb, bx, by, bz)
				}
			}
		}
	}
}

type chunkPosition struct {
	X, Z int
}

type chunk struct {
	chunkPosition

	Sections [16]*chunkSection
	Biomes   [16 * 16]byte
}

func (c *chunk) block(x, y, z int) Block {
	s := y >> 4
	if s < 0 || s > 15 {
		return BlockAir.Base
	}
	sec := c.Sections[s]
	if sec == nil {
		return BlockAir.Base
	}
	return sec.block(x, y&0xF, z)
}

func (c *chunk) setBlock(b Block, x, y, z int) {
	s := y >> 4
	if s < 0 || s > 15 {
		return
	}
	sec := c.Sections[s]
	if sec == nil {
		sec = newChunkSection(c, s)
		sec.Buffer = render.AllocateChunkBuffer(c.X, s, c.Z)
	}
	sec.setBlock(b, x, y&0xF, z)
	var maxB, maxS int8
	for _, d := range direction.Values {
		ox, oy, oz := d.Offset()
		l := int8(c.relLight(x+ox, y+oy, z+oz, (*chunkSection).blockLight, false)) - 1
		if l > maxB {
			maxB = l
		}
		l = int8(c.relLight(x+ox, y+oy, z+oz, (*chunkSection).skyLight, true))
		if !(l == 15 && d == direction.Up) {
			l--
		}
		if l > maxS {
			maxS = l
		}
	}
	updateLight(c, specialLight, maxB, x, y, z, (*chunkSection).blockLight, (*chunkSection).setBlockLight, false)
	updateLight(c, specialLight, maxS, x, y, z, (*chunkSection).skyLight, (*chunkSection).setSkyLight, true)
}

const specialLight int8 = -55

type getLight func(cs *chunkSection, x, y, z int) byte
type setLight func(cs *chunkSection, l byte, x, y, z int)

func clampInt8(x, min, max int8) int8 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

type lightState struct {
	chunk      *chunk
	exLight, l int8
	x, y, z    int
}

func updateLight(c *chunk, exLight, l int8, x, y, z int, get getLight, set setLight, sky bool) {
	queue := []lightState{
		{c, exLight, l, x, y, z},
	}
itQueue:
	for len(queue) > 0 {
		// Take the first item from the queue
		state := queue[0]
		queue = queue[1:]
		c := state.chunk
		exLight, l, x, y, z = state.exLight, state.l, state.x, state.y, state.z
		// Handle neighbor chunks
		if x < 0 || x > 15 || z < 0 || z > 15 {
			ch := chunkMap[chunkPosition{c.X + (x >> 4), c.Z + (z >> 4)}]
			if ch == nil {
				continue itQueue
			}
			x &= 0xF
			z &= 0xF
			// ch.updateLight(exLight, l, x, y, z, get, set, sky)
			queue = append(queue, lightState{ch, exLight, l, x, y, z})
			continue itQueue
		}
		s := y >> 4
		sec := c.Sections[s]
		if sec == nil {
			continue itQueue
		}
		// Needs a redraw after changing the lighting
		sec.dirty = true
		y &= 0xF
		b := sec.block(x, y, z)
		curL := int8(get(sec, x, y, z))
		l -= int8(b.LightReduction())
		if !sky {
			l += int8(b.LightEmitted())
		}
		l = clampInt8(l, 0, 15)
		ex := exLight - int8(b.LightReduction())
		if !sky {
			ex += int8(b.LightEmitted())
		}
		ex = clampInt8(ex, 0, 15)
		// If the light isn't what we expect it to be or its already
		// at the value we want to change it too then don't update
		// this position.
		if (exLight != specialLight && ex != curL) || curL == l {
			continue itQueue
		}
		set(sec, byte(l), x, y, z)
		// Update the surrounding blocks
		for _, d := range direction.Values {
			ox, oy, oz := d.Offset()
			nl := l
			ex := curL
			if !(sky && d == direction.Down && nl == 15) {
				nl--
				if nl < 0 {
					nl = 0
				}
			}
			if !(sky && d == direction.Down && ex == 15) {
				ex--
				if ex < 0 {
					ex = 0
				}
			}
			// c.updateLight(ex, nl, x+ox, (sec.Y<<4)+y+oy, z+oz, get, set, sky)
			queue = append(queue, lightState{c, ex, nl, x + ox, (sec.Y << 4) + y + oy, z + oz})
		}
	}
}

func (c *chunk) relLight(x, y, z int, f getLight, sky bool) byte {
	ch := c
	if x < 0 || x > 15 || z < 0 || z > 15 {
		ch = chunkMap[chunkPosition{c.X + (x >> 4), c.Z + (z >> 4)}]
		x &= 0xF
		z &= 0xF
	}
	if ch == nil || y < 0 || y > 255 {
		return 0
	}
	s := y >> 4
	sec := ch.Sections[s]
	if sec == nil {
		if sky {
			return 15
		}
		return 0
	}
	return f(sec, x&0xF, y&0xF, z&0xF)
}

func (c *chunk) biome(x, z int) *biome.Type {
	return biome.ById(c.Biomes[z<<4|x])
}

func (c *chunk) free() {
	for _, s := range c.Sections {
		if s != nil {
			s.Buffer.Free()
		}
	}
}

type chunkSection struct {
	chunk *chunk
	Y     int

	Blocks     [16 * 16 * 16]uint16
	BlockLight nibble.Array
	SkyLight   nibble.Array

	Buffer *render.ChunkBuffer

	dirty    bool
	building bool
}

func (cs *chunkSection) block(x, y, z int) Block {
	return allBlocks[cs.Blocks[(y<<8)|(z<<4)|x]]
}

func (cs *chunkSection) setBlock(b Block, x, y, z int) {
	cs.Blocks[(y<<8)|(z<<4)|x] = b.SID()
	cs.dirty = true
}

func (cs *chunkSection) blockLight(x, y, z int) byte {
	return cs.BlockLight.Get((y << 8) | (z << 4) | x)
}

func (cs *chunkSection) setBlockLight(l byte, x, y, z int) {
	cs.BlockLight.Set((y<<8)|(z<<4)|x, l)
}

func (cs *chunkSection) skyLight(x, y, z int) byte {
	return cs.SkyLight.Get((y << 8) | (z << 4) | x)
}
func (cs *chunkSection) setSkyLight(l byte, x, y, z int) {
	cs.SkyLight.Set((y<<8)|(z<<4)|x, l)
}

func newChunkSection(c *chunk, y int) *chunkSection {
	cs := &chunkSection{
		chunk:      c,
		Y:          y,
		BlockLight: nibble.New(16 * 16 * 16),
		SkyLight:   nibble.New(16 * 16 * 16),
	}
	for i := range cs.Blocks {
		cs.Blocks[i] = BlockAir.Blocks[0].SID()
	}
	return cs
}

func loadChunk(x, z int, data []byte, mask uint16, sky, isNew bool) int {
	var c *chunk
	if isNew {
		c = &chunk{
			chunkPosition: chunkPosition{
				X: x, Z: z,
			},
		}
	} else {
		c = chunkMap[chunkPosition{
			X: x, Z: z,
		}]
		if c == nil {
			return 0
		}
	}

	for i := 0; i < 16; i++ {
		if mask&(1<<uint(i)) == 0 {
			continue
		}
		if c.Sections[i] == nil {
			c.Sections[i] = newChunkSection(c, i)
		}
	}
	offset := 0
	for i, section := range c.Sections {
		if section == nil || mask&(1<<uint(i)) == 0 {
			continue
		}

		for i := 0; i < 16*16*16; i++ {
			section.Blocks[i] = GetBlockByCombinedID(binary.LittleEndian.Uint16(data[offset:])).SID()
			offset += 2
		}
	}
	for i, section := range c.Sections {
		if section == nil || mask&(1<<uint(i)) == 0 {
			continue
		}
		copy(section.BlockLight, data[offset:])
		offset += len(section.BlockLight)
	}
	if sky {
		for i, section := range c.Sections {
			if section == nil || mask&(1<<uint(i)) == 0 {
				continue
			}
			copy(section.SkyLight, data[offset:])
			offset += len(section.BlockLight)
		}
	}

	if isNew {
		copy(c.Biomes[:], data[offset:])
		offset += len(c.Biomes)
	}

	syncChan <- func() {
		// Allocate the render buffers sync
		for y, section := range c.Sections {
			if section != nil && section.Buffer == nil {
				section.Buffer = render.AllocateChunkBuffer(c.X, y, c.Z)
			}
		}

		chunkMap[c.chunkPosition] = c
		for _, section := range c.Sections {
			if section == nil {
				continue
			}
			cx := c.X << 4
			cy := section.Y << 4
			cz := c.Z << 4
			for y := 0; y < 16; y++ {
				for z := 0; z < 16; z++ {
					for x := 0; x < 16; x++ {
						section.setBlock(
							section.block(x, y, z).UpdateState(cx+x, cy+y, cz+z),
							x, y, z,
						)
					}
				}
			}
		}

		for xx := -1; xx <= 1; xx++ {
			for zz := -1; zz <= 1; zz++ {
				c := chunkMap[chunkPosition{x + xx, z + zz}]
				if c != nil {
					for _, section := range c.Sections {
						if section == nil {
							continue
						}
						cx, cy, cz := c.X<<4, section.Y<<4, c.Z<<4
						for y := 0; y < 16; y++ {
							if !(xx != 0 && zz != 0) {
								// Row/Col
								for i := 0; i < 16; i++ {
									var bx, bz int
									if xx != 0 {
										bz = i
										if xx == -1 {
											bx = 15
										}
									} else {
										bx = i
										if zz == -1 {
											bz = 15
										}
									}
									section.setBlock(
										section.block(bx, y, bz).UpdateState(cx+bx, cy+y, cz+bz),
										bx, y, bz,
									)
								}
							} else {
								// Just the corner
								var bx, bz int
								if xx == -1 {
									bx = 15
								}
								if zz == -1 {
									bz = 15
								}
								section.setBlock(
									section.block(bx, y, bz).UpdateState(cx+bx, cy+y, cz+bz),
									bx, y, bz,
								)
							}
						}
						section.dirty = true
					}
				}
			}
		}
	}

	return offset
}

func sortedChunks() []*chunk {
	out := make([]*chunk, len(chunkMap))
	i := 0
	for _, c := range chunkMap {
		out[i] = c
		i++
	}
	sort.Sort(chunkSorter(out))
	return out
}

type chunkSorter []*chunk

func (cs chunkSorter) Len() int {
	return len(cs)
}

func (cs chunkSorter) Less(a, b int) bool {
	ac := cs[a]
	bc := cs[b]
	xx := float64(ac.X<<4+8) - Client.X
	zz := float64(ac.Z<<4+8) - Client.Z
	adist := xx*xx + zz*zz
	xx = float64(bc.X<<4+8) - Client.X
	zz = float64(bc.Z<<4+8) - Client.Z
	bdist := xx*xx + zz*zz
	return adist < bdist
}

func (cs chunkSorter) Swap(a, b int) {
	cs[a], cs[b] = cs[b], cs[a]
}
