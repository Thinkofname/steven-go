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
}

func (w world) UpdateBlock(x, y, z int) {
	for yy := -1; yy <= 1; yy++ {
		for zz := -1; zz <= 1; zz++ {
			for xx := -1; xx <= 1; xx++ {
				bx, by, bz := x+xx, y+yy, z+zz
				w.SetBlock(w.Block(bx, by, bz).UpdateState(bx, by, bz), bx, by, bz)
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
		return
	}
	sec.setBlock(b, x, y&0xF, z)
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

	Blocks     [16 * 16 * 16]Block
	BlockLight nibble.Array
	SkyLight   nibble.Array

	Buffer *render.ChunkBuffer

	dirty    bool
	building bool
}

func (cs *chunkSection) block(x, y, z int) Block {
	return cs.Blocks[(y<<8)|(z<<4)|x]
}

func (cs *chunkSection) setBlock(b Block, x, y, z int) {
	cs.Blocks[(y<<8)|(z<<4)|x] = b
	cs.dirty = true
}

func (cs *chunkSection) blockLight(x, y, z int) byte {
	return cs.BlockLight.Get((y << 8) | (z << 4) | x)
}

func (cs *chunkSection) skyLight(x, y, z int) byte {
	return cs.SkyLight.Get((y << 8) | (z << 4) | x)
}

func newChunkSection(c *chunk, y int) *chunkSection {
	cs := &chunkSection{
		chunk:      c,
		Y:          y,
		BlockLight: nibble.New(16 * 16 * 16),
		SkyLight:   nibble.New(16 * 16 * 16),
	}
	for i := range cs.Blocks {
		cs.Blocks[i] = BlockAir.Blocks[0]
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
			section.Blocks[i] = GetBlockByCombinedID(binary.LittleEndian.Uint16(data[offset:]))
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
	xx := float64(ac.X<<4+8) - render.Camera.X
	zz := float64(ac.Z<<4+8) - render.Camera.Z
	adist := xx*xx + zz*zz
	xx = float64(bc.X<<4+8) - render.Camera.X
	zz = float64(bc.Z<<4+8) - render.Camera.Z
	bdist := xx*xx + zz*zz
	return adist < bdist
}

func (cs chunkSorter) Swap(a, b int) {
	cs[a], cs[b] = cs[b], cs[a]
}
