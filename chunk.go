package main

import (
	"encoding/binary"
	"sort"

	"github.com/thinkofdeath/steven/nibble"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/world/biome"
)

var chunkMap = map[chunkPosition]*chunk{}

type chunkPosition struct {
	X, Z int
}

type chunk struct {
	chunkPosition

	Sections [16]*chunkSection
	Biomes   [16 * 16]byte
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

func loadChunk(x, z int, data []byte, mask uint16, sky, hasBiome bool) int {
	c := &chunk{
		chunkPosition: chunkPosition{
			X: x, Z: z,
		},
	}

	for i := 0; i < 16; i++ {
		if mask&(1<<uint(i)) == 0 {
			continue
		}
		c.Sections[i] = newChunkSection(c, i)
	}
	offset := 0
	for _, section := range c.Sections {
		if section == nil {
			continue
		}

		for i := 0; i < 16*16*16; i++ {
			section.Blocks[i] = GetBlockByCombinedID(binary.LittleEndian.Uint16(data[offset:]))
			offset += 2
		}
	}
	for _, section := range c.Sections {
		if section == nil {
			continue
		}
		copy(section.BlockLight, data[offset:])
		offset += len(section.BlockLight)
	}
	if sky {
		for _, section := range c.Sections {
			if section == nil {
				continue
			}
			copy(section.SkyLight, data[offset:])
			offset += len(section.BlockLight)
		}
	}

	if hasBiome {
		copy(c.Biomes[:], data[offset:])
		offset += len(c.Biomes)
	}

	syncChan <- func() {
		// Allocate the render buffers sync
		for y, section := range c.Sections {
			if section != nil {
				section.Buffer = render.AllocateChunkBuffer(c.X, y, c.Z)
			}
		}

		chunkMap[c.chunkPosition] = c

		for xx := -1; xx <= 1; xx++ {
			for zz := -1; zz <= 1; zz++ {
				c := chunkMap[chunkPosition{x + xx, z + zz}]
				if c != nil {
					for _, section := range c.Sections {
						if section == nil {
							continue
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
