package main

import (
	"encoding/binary"

	"github.com/thinkofdeath/steven/nibble"
	"github.com/thinkofdeath/steven/render"
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

	Buffer   *render.ChunkBuffer
	renderID int
}

func (cs *chunkSection) block(x, y, z int) uint16 {
	return cs.Blocks[(y<<8)|(z<<4)|x]
}

func newChunkSection(c *chunk, x, y, z int) *chunkSection {
	return &chunkSection{
		chunk:      c,
		Y:          y,
		BlockLight: nibble.New(16 * 16 * 16),
		SkyLight:   nibble.New(16 * 16 * 16),
		Buffer:     render.AllocateChunkBuffer(x, y, z),
	}
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
		c.Sections[i] = newChunkSection(c, c.X, i, c.Z)
	}
	offset := 0
	for _, section := range c.Sections {
		if section == nil {
			continue
		}

		for i := 0; i < 16*16*16; i++ {
			section.Blocks[i] = binary.LittleEndian.Uint16(data[offset:])
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

	chunkMap[c.chunkPosition] = c

	for _, section := range c.Sections {
		if section == nil {
			continue
		}
		section.build()
	}

	return offset
}
