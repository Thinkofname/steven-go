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
	"reflect"

	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/render/builder"
	"github.com/thinkofdeath/steven/type/direction"
)

type blockLiquid struct {
	baseBlock
	Lava  bool
	Level int `state:"level,0-15"`
}

func (l *blockLiquid) load(tag reflect.StructTag) {
	getBool := wrapTagBool(tag)
	l.Lava = getBool("lava", false)
	l.cullAgainst = false
	l.collidable = false
	if !l.Lava {
		l.translucent = true
	}
}

func (l *blockLiquid) LightReduction() int {
	if l.Lava {
		return 0
	}
	return 1
}

func (l *blockLiquid) LightEmitted() int {
	if l.Lava {
		return 15
	}
	return 0
}

func (l *blockLiquid) toData() int {
	return l.Level
}

func (l *blockLiquid) renderLiquid(bs *blocksSnapshot, x, y, z int, buf *builder.Buffer, indices *int) {
	var tex render.TextureInfo
	var b1, b2 *BlockSet
	if l.Lava {
		b1 = Blocks.Lava
		b2 = Blocks.FlowingLava
		tex = render.GetTexture("blocks/lava_still")
	} else {
		b1 = Blocks.Water
		b2 = Blocks.FlowingWater
		tex = render.GetTexture("blocks/water_still")
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

	for f, d := range direction.Values {
		ox, oy, oz := d.Offset()
		special := d == direction.Up && (tl < 8 || tr < 8 || bl < 8 || br < 8)
		if b := bs.block(x+ox, y+oy, z+oz); special || (!b.Is(b1) && !b.Is(b2) && !b.ShouldCullAgainst()) {
			vert := faceVertices[f]

			var cr, cg, cb byte
			cr = 255
			cg = 255
			cb = 255

			*indices += len(vert.indices)

			// TODO: Needs fixing (maybe?)
			rect := tex.Rect()
			ux1 := int16(0)
			ux2 := int16(16 * rect.Width)
			uy1 := int16(0)
			uy2 := int16(16 * rect.Height)
			for _, vert := range vert.verts {
				vert.TX = uint16(rect.X)
				vert.TY = uint16(rect.Y)
				vert.TW = uint16(rect.Width)
				vert.TH = uint16(rect.Height)
				vert.TAtlas = int16(tex.Atlas())
				vert.R = cr
				vert.G = cg
				vert.B = cb

				if vert.Y == 0 {
					vert.Y = int16(0 + y*256)
				} else {
					if vert.X == 0 && vert.Z == 0 {
						height := int((16.0/8.0)*float64(tl)) * 16
						vert.Y = int16(height + y*256)
					} else if vert.X != 0 && vert.Z == 0 {
						height := int((16.0/8.0)*float64(tr)) * 16
						vert.Y = int16(height + y*256)
					} else if vert.X == 0 && vert.Z != 0 {
						height := int((16.0/8.0)*float64(bl)) * 16
						vert.Y = int16(height + y*256)
					} else {
						height := int((16.0/8.0)*float64(br)) * 16
						vert.Y = int16(height + y*256)
					}
				}

				if vert.X == 0 {
					vert.X = int16(0 + x*256)
				} else {
					vert.X = int16(256 + x*256)
				}
				if vert.Z == 0 {
					vert.Z = int16(0 + z*256)
				} else {
					vert.Z = int16(256 + z*256)
				}

				vert.BlockLight, vert.SkyLight = calculateLight(
					bs,
					x, y, z,
					float64(vert.X)/256.0,
					float64(vert.Y)/256.0,
					float64(vert.Z)/256.0,
					1, !l.Lava, l.ForceShade(),
				)

				if vert.TOffsetX == 0 {
					vert.TOffsetX = int16(ux1)
				} else {
					vert.TOffsetX = int16(ux2)
				}
				if vert.TOffsetY == 0 {
					vert.TOffsetY = int16(uy1)
				} else {
					vert.TOffsetY = int16(uy2)
				}
				buildVertex(buf, vert)
			}
		}
	}
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
