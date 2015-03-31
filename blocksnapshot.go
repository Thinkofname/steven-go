package main

import (
	"math"

	"github.com/thinkofdeath/steven/type/nibble"
	"github.com/thinkofdeath/steven/world/biome"
)

type blocksSnapshot struct {
	Blocks     []Block
	BlockLight nibble.Array
	SkyLight   nibble.Array
	Biome      []*biome.Type

	x, y, z int
	w, h, d int
}

func getSnapshot(x, y, z, w, h, d int) *blocksSnapshot {
	bs := &blocksSnapshot{
		Blocks:     make([]Block, w*h*d),
		BlockLight: nibble.New(w * h * d),
		SkyLight:   nibble.New(w * h * d),
		Biome:      make([]*biome.Type, w*d),
		x:          x,
		y:          y,
		z:          z,
		w:          w,
		h:          h,
		d:          d,
	}
	for i := range bs.Blocks {
		bs.Blocks[i] = BlockAir.Blocks[0]
		bs.SkyLight.Set(i, 15)
	}
	for i := range bs.Biome {
		bs.Biome[i] = biome.Invalid
	}

	cx1 := int(math.Floor(float64(x) / 16.0))
	cx2 := int(math.Ceil(float64(x+w) / 16.0))
	cy1 := int(math.Floor(float64(y) / 16.0))
	cy2 := int(math.Ceil(float64(y+h) / 16.0))
	cz1 := int(math.Floor(float64(z) / 16.0))
	cz2 := int(math.Ceil(float64(z+d) / 16.0))

	for cx := cx1; cx < cx2; cx++ {
		for cz := cz1; cz < cz2; cz++ {
			chunk := chunkMap[chunkPosition{cx, cz}]
			if chunk == nil {
				continue
			}
			for cy := cy1; cy < cy2; cy++ {
				if cy < 0 || cy > 15 {
					continue
				}
				cs := chunk.Sections[cy]
				if cs == nil {
					continue
				}
				x1 := x - cx<<4
				x2 := x + w - cx<<4
				y1 := y - cy<<4
				y2 := y + h - cy<<4
				z1 := z - cz<<4
				z2 := z + d - cz<<4

				if x1 < 0 {
					x1 = 0
				}
				if x2 > 16 {
					x2 = 16
				}
				if y1 < 0 {
					y1 = 0
				}
				if y2 > 16 {
					y2 = 16
				}
				if z1 < 0 {
					z1 = 0
				}
				if z2 > 16 {
					z2 = 16
				}

				for yy := y1; yy < y2; yy++ {
					for zz := z1; zz < z2; zz++ {
						for xx := x1; xx < x2; xx++ {
							bl := cs.block(xx, yy, zz)
							ox, oy, oz := xx+(cx<<4), yy+(cy<<4), zz+(cz<<4)
							bs.setBlock(ox, oy, oz, bl)
							bs.setBlockLight(ox, oy, oz, cs.blockLight(xx, yy, zz))
							bs.setSkyLight(ox, oy, oz, cs.skyLight(xx, yy, zz))

							bs.setBiome(ox, oz, chunk.biome(xx, zz))
						}
					}
				}
			}
		}
	}

	return bs
}

func (bs *blocksSnapshot) block(x, y, z int) Block {
	return bs.Blocks[bs.index(x, y, z)]
}

func (bs *blocksSnapshot) setBlock(x, y, z int, b Block) {
	bs.Blocks[bs.index(x, y, z)] = b
}

func (bs *blocksSnapshot) blockLight(x, y, z int) byte {
	return bs.BlockLight.Get(bs.index(x, y, z))
}

func (bs *blocksSnapshot) setBlockLight(x, y, z int, b byte) {
	bs.BlockLight.Set(bs.index(x, y, z), b)
}

func (bs *blocksSnapshot) skyLight(x, y, z int) byte {
	return bs.SkyLight.Get(bs.index(x, y, z))
}

func (bs *blocksSnapshot) setSkyLight(x, y, z int, b byte) {
	bs.SkyLight.Set(bs.index(x, y, z), b)
}

func (bs *blocksSnapshot) biome(x, z int) *biome.Type {
	x -= bs.x
	z -= bs.z
	return bs.Biome[(z*bs.w)|x]
}

func (bs *blocksSnapshot) setBiome(x, z int, b *biome.Type) {
	x -= bs.x
	z -= bs.z
	bs.Biome[(z*bs.w)|x] = b
}

func (bs *blocksSnapshot) index(x, y, z int) int {
	x -= bs.x
	y -= bs.y
	z -= bs.z
	return x + z*bs.w + y*bs.w*bs.d
}
