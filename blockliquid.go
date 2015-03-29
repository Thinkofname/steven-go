package main

import "github.com/thinkofdeath/steven/render"

type blockLiquid struct {
	baseBlock
	Lava  bool
	Level int `state:"level,0-15"`
}

func initLiquid(name string, lava bool) *BlockSet {
	l := &blockLiquid{}
	l.init(name)
	l.Lava = lava
	l.cullAgainst = false
	set := alloc(l)
	return set
}

func (l *blockLiquid) String() string {
	return l.Parent.stringify(l)
}

func (l *blockLiquid) clone() Block {
	return &blockLiquid{
		baseBlock: *(l.baseBlock.clone().(*baseBlock)),
		Level:     l.Level,
		Lava:      l.Lava,
	}
}

func (l *blockLiquid) toData() int {
	return l.Level
}

func (l *blockLiquid) renderLiquid(bs *blocksSnapshot, x, y, z int) []chunkVertex {
	var out []chunkVertex
	var tex *render.TextureInfo
	var b1, b2 *BlockSet
	if l.Lava {
		b1 = BlockLava
		b2 = BlockFlowingLava
		tex = render.GetTexture("lava_still")
	} else {
		b1 = BlockWater
		b2 = BlockFlowingWater
		tex = render.GetTexture("water_still")
	}

	var tl, tr, bl, br int
	if b := bs.block(x, y+1, z); b.Is(b1) || b.Is(b2) {
		tl = 8
		tr = 8
		bl = 8
		br = 8
	} else {
		tl = l.averageLiquidLevel(bs, x, y, z)
		tr = l.averageLiquidLevel(bs, x+1, y, z)
		bl = l.averageLiquidLevel(bs, x, y, z+1)
		br = l.averageLiquidLevel(bs, x+1, y, z+1)
	}

	for f := range faceVertices {
		var ox, oy, oz int
		switch f {
		case 0: // Up
			oy = 1
		case 1: // Down
			oy = -1
		case 2: // North
			oz = -1
		case 3: // South
			oz = 1
		case 4: // West
			ox = -1
		case 5: // East
			ox = 1
		}
		if b := bs.block(x+ox, y+oy, z+oz); !b.Is(b1) && !b.Is(b2) {
			vert := faceVertices[f]

			var cr, cg, cb byte
			cr = 255
			cg = 255
			cb = 255

			// TODO: Needs fixing (maybe?)
			ux1 := int16(0)
			ux2 := int16(16 * tex.Width)
			uy1 := int16(0)
			uy2 := int16(16 * tex.Height)
			for v := range vert {
				vert[v].TX = uint16(tex.X)
				vert[v].TY = uint16(tex.Y + tex.Atlas*1024.0)
				vert[v].TW = uint16(tex.Width)
				vert[v].TH = uint16(tex.Height)
				vert[v].R = cr
				vert[v].G = cg
				vert[v].B = cb

				if vert[v].Y == 0 {
					vert[v].Y = int16(0 + y*256)
				} else {
					if vert[v].X == 0 && vert[v].Z == 0 {
						height := int((16.0/8.0)*float64(tl)) * 16
						vert[v].Y = int16(height + y*256)
					} else if vert[v].X != 0 && vert[v].Z == 0 {
						height := int((16.0/8.0)*float64(tr)) * 16
						vert[v].Y = int16(height + y*256)
					} else if vert[v].X == 0 && vert[v].Z != 0 {
						height := int((16.0/8.0)*float64(bl)) * 16
						vert[v].Y = int16(height + y*256)
					} else {
						height := int((16.0/8.0)*float64(br)) * 16
						vert[v].Y = int16(height + y*256)
					}
				}

				if vert[v].X == 0 {
					vert[v].X = int16(0 + x*256)
				} else {
					vert[v].X = int16(256 + x*256)
				}
				if vert[v].Z == 0 {
					vert[v].Z = int16(0 + z*256)
				} else {
					vert[v].Z = int16(256 + z*256)
				}

				vert[v].BlockLight, vert[v].SkyLight = calculateLight(
					bs,
					x, y, z,
					float64(vert[v].X)/256.0,
					float64(vert[v].Y)/256.0,
					float64(vert[v].Z)/256.0,
					1, true, l.ForceShade(),
				)

				if vert[v].TOffsetX == 0 {
					vert[v].TOffsetX = int16(ux1)
				} else {
					vert[v].TOffsetX = int16(ux2)
				}
				if vert[v].TOffsetY == 0 {
					vert[v].TOffsetY = int16(uy1)
				} else {
					vert[v].TOffsetY = int16(uy2)
				}
			}
			out = append(out, vert[:]...)
		}
	}
	return out
}

func (l *blockLiquid) averageLiquidLevel(bs *blocksSnapshot, x, y, z int) int {
	level := 0
	for xx := -1; xx <= 0; xx++ {
		for zz := -1; zz <= 0; zz++ {
			b := bs.block(x+xx, y+1, z+zz)
			if o, ok := b.(*blockLiquid); ok && l.Lava == o.Lava {
				return 8
			}
			b = bs.block(x+xx, y, z+zz)
			if o, ok := b.(*blockLiquid); ok && l.Lava == o.Lava {
				nl := 7 - (o.Level & 0x7)
				if nl > level {
					level = nl
				}
			}
		}
	}
	return level
}
