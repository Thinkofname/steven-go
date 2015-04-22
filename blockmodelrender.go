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
	"image"
	"math"

	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/render/builder"
	"github.com/thinkofdeath/steven/type/direction"
)

type processedModel struct {
	faces            []processedFace
	ambientOcclusion bool
}

type processedFace struct {
	cullFace  direction.Type
	facing    direction.Type
	vertices  []chunkVertex
	shade     bool
	tintIndex int
}

var faceRotation = []direction.Type{
	direction.North,
	direction.East,
	direction.South,
	direction.West,
}

var faceRotationX = []direction.Type{
	direction.North,
	direction.Down,
	direction.South,
	direction.Up,
}

func rotateDirection(val direction.Type, offset int, rots []direction.Type, invalid ...direction.Type) direction.Type {
	for _, d := range invalid {
		if d == val {
			return val
		}
	}
	var pos int
	for di, d := range rots {
		if d == val {
			pos = di
			break
		}
	}
	return rots[(pos+offset)%len(rots)]
}

func precomputeModel(bm *blockModel) *processedModel {
	p := &processedModel{}
	p.ambientOcclusion = bm.ambientOcclusion
	for ei := range bm.elements {
		// Render the last element first so that
		// grass's overlay works correctly.
		el := bm.elements[len(bm.elements)-1-ei]
		for i, face := range el.faces {
			faceID := direction.Type(i)
			if face == nil {
				continue
			}
			pFace := processedFace{}
			cullFace := face.cullFace
			if bm.x > 0 {
				o := int(bm.x) / 90
				cullFace = rotateDirection(cullFace, o, faceRotationX, direction.East, direction.West, direction.Invalid)
				faceID = rotateDirection(faceID, o, faceRotationX, direction.East, direction.West, direction.Invalid)
			}
			if bm.y > 0 {
				o := int(bm.y) / 90
				cullFace = rotateDirection(cullFace, o, faceRotation, direction.Up, direction.Down, direction.Invalid)
				faceID = rotateDirection(faceID, o, faceRotation, direction.Up, direction.Down, direction.Invalid)
			}
			pFace.cullFace = cullFace
			pFace.facing = direction.Type(faceID)
			pFace.tintIndex = face.tintIndex
			pFace.shade = el.shade

			vert := faceVertices[i]
			tex := bm.lookupTexture(face.texture)

			ux1 := int16(face.uv[0] * float64(tex.Width))
			ux2 := int16(face.uv[2] * float64(tex.Width))
			uy1 := int16(face.uv[1] * float64(tex.Height))
			uy2 := int16(face.uv[3] * float64(tex.Height))

			tw, th := int16(tex.Width), int16(tex.Height)
			if face.rotation > 0 {
				x := ux1
				y := uy1
				w := ux2 - ux1
				h := uy2 - uy1
				switch face.rotation {
				case 90:
					uy2 = x + w
					ux1 = tw*16 - (y + h)
					ux2 = tw*16 - y
					uy1 = x
				case 180:
					uy1 = th*16 - (y + h)
					uy2 = th*16 - y
					ux1 = x + w
					ux2 = x
				case 270:
					uy2 = x
					uy1 = x + w
					ux2 = y + h
					ux1 = y
				}
			}

			var minX, minY, minZ int16 = math.MaxInt16, math.MaxInt16, math.MaxInt16
			var maxX, maxY, maxZ int16 = math.MinInt16, math.MinInt16, math.MinInt16

			for v := range vert {
				vert[v].TX = uint16(tex.X)
				vert[v].TY = uint16(tex.Y + tex.Atlas*render.AtlasSize)
				vert[v].TW = uint16(tex.Width)
				vert[v].TH = uint16(tex.Height)

				if vert[v].X == 0 {
					vert[v].X = int16(el.from[0] * 16)
				} else {
					vert[v].X = int16(el.to[0] * 16)
				}
				if vert[v].Y == 0 {
					vert[v].Y = int16(el.from[1] * 16)
				} else {
					vert[v].Y = int16(el.to[1] * 16)
				}
				if vert[v].Z == 0 {
					vert[v].Z = int16(el.from[2] * 16)
				} else {
					vert[v].Z = int16(el.to[2] * 16)
				}

				if el.rotation != nil {
					r := el.rotation
					switch r.axis {
					case "y":
						rotY := -r.angle * (math.Pi / 180)
						c := math.Cos(rotY)
						s := math.Sin(rotY)
						x := float64(vert[v].X) - r.origin[0]*16
						z := float64(vert[v].Z) - r.origin[2]*16
						vert[v].X = int16(r.origin[0]*16 + (x*c - z*s))
						vert[v].Z = int16(r.origin[2]*16 + (z*c + x*s))
					case "x":
						rotX := -r.angle * (math.Pi / 180)
						c := math.Cos(-rotX)
						s := math.Sin(-rotX)
						z := float64(vert[v].Z) - r.origin[2]*16
						y := float64(vert[v].Y) - r.origin[1]*16
						vert[v].Z = int16(r.origin[2]*16 + (z*c - y*s))
						vert[v].Y = int16(r.origin[1]*16 + (y*c + z*s))
					case "z":
						rotZ := -r.angle * (math.Pi / 180)
						c := math.Cos(-rotZ)
						s := math.Sin(-rotZ)
						x := float64(vert[v].X) - r.origin[0]*16
						y := float64(vert[v].Y) - r.origin[1]*16
						vert[v].X = int16(r.origin[0]*16 + (x*c - y*s))
						vert[v].Y = int16(r.origin[1]*16 + (y*c + x*s))
					}
				}

				if bm.x > 0 {
					rotX := bm.x * (math.Pi / 180)
					c := int16(math.Cos(rotX))
					s := int16(math.Sin(rotX))
					z := vert[v].Z - 8*16
					y := vert[v].Y - 8*16
					vert[v].Z = 8*16 + int16(z*c-y*s)
					vert[v].Y = 8*16 + int16(y*c+z*s)
				}

				if bm.y > 0 {
					rotY := bm.y * (math.Pi / 180)
					c := int16(math.Cos(rotY))
					s := int16(math.Sin(rotY))
					x := vert[v].X - 8*16
					z := vert[v].Z - 8*16
					vert[v].X = 8*16 + int16(x*c-z*s)
					vert[v].Z = 8*16 + int16(z*c+x*s)
				}

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

				if face.rotation > 0 {
					rotY := -float64(face.rotation) * (math.Pi / 180)
					c := int16(math.Cos(rotY))
					s := int16(math.Sin(rotY))
					x := vert[v].TOffsetX - 8*tw
					y := vert[v].TOffsetY - 8*th
					vert[v].TOffsetX = 8*tw + int16(x*c-y*s)
					vert[v].TOffsetY = 8*th + int16(y*c+x*s)
				}

				if bm.uvLock && bm.y > 0 &&
					(pFace.facing == direction.Up || pFace.facing == direction.Down) {
					rotY := float64(-bm.y) * (math.Pi / 180)
					c := int16(math.Cos(rotY))
					s := int16(math.Sin(rotY))
					x := vert[v].TOffsetX - 8*16
					y := vert[v].TOffsetY - 8*16
					vert[v].TOffsetX = 8*16 + int16(x*c+y*s)
					vert[v].TOffsetY = 8*16 + int16(y*c-x*s)
				}

				if bm.uvLock && bm.x > 0 &&
					(pFace.facing != direction.Up && pFace.facing != direction.Down) {
					rotY := float64(bm.x) * (math.Pi / 180)
					c := int16(math.Cos(rotY))
					s := int16(math.Sin(rotY))
					x := vert[v].TOffsetX - 8*16
					y := vert[v].TOffsetY - 8*16
					vert[v].TOffsetX = 8*16 + int16(x*c+y*s)
					vert[v].TOffsetY = 8*16 + int16(y*c-x*s)
				}

				if el.rotation != nil && el.rotation.rescale {
					if vert[v].X < minX {
						minX = vert[v].X
					} else if vert[v].X > maxX {
						maxX = vert[v].X
					}
					if vert[v].Y < minY {
						minY = vert[v].Y
					} else if vert[v].Y > maxY {
						maxY = vert[v].Y
					}
					if vert[v].Z < minZ {
						minZ = vert[v].Z
					} else if vert[v].Z > maxZ {
						maxZ = vert[v].Z
					}
				}
			}

			if el.rotation != nil && el.rotation.rescale {
				diffX := float64(maxX - minX)
				diffY := float64(maxY - minY)
				diffZ := float64(maxZ - minZ)
				for v := range vert {
					vert[v].X = int16((float64(vert[v].X-minX) / diffX) * 256)
					vert[v].Y = int16((float64(vert[v].Y-minY) / diffY) * 256)
					vert[v].Z = int16((float64(vert[v].Z-minZ) / diffZ) * 256)
				}
			}

			pFace.vertices = vert[:]

			p.faces = append(p.faces, pFace)
		}
	}
	return p
}

func (p processedModel) Render(x, y, z int, bs *blocksSnapshot, buf *builder.Buffer) {
	this := bs.block(x, y, z)
	for _, f := range p.faces {
		if f.cullFace != direction.Invalid {
			ox, oy, oz := f.cullFace.Offset()
			if b := bs.block(x+ox, y+oy, z+oz); b.ShouldCullAgainst() || b == this {
				continue
			}
		}

		var cr, cg, cb byte
		cr = 255
		cg = 255
		cb = 255
		if this.TintImage() != nil {
			switch f.tintIndex {
			case 0:
				cr, cg, cb = calculateBiome(bs, x, z, this.TintImage())
			}
		}
		if f.facing == direction.West || f.facing == direction.East {
			cr = byte(float64(cr) * 0.8)
			cg = byte(float64(cg) * 0.8)
			cb = byte(float64(cb) * 0.8)
		}

		for _, vert := range f.vertices {
			vert.R = cr
			vert.G = cg
			vert.B = cb

			vert.X += int16(x * 256)
			vert.Y += int16(y * 256)
			vert.Z += int16(z * 256)

			vert.BlockLight, vert.SkyLight = calculateLight(
				bs,
				x, y, z,
				float64(vert.X)/256.0,
				float64(vert.Y)/256.0,
				float64(vert.Z)/256.0,
				int(f.facing), p.ambientOcclusion, this.ForceShade(),
			)
			buildVertex(buf, vert)
		}
	}
}

// Takes an average of the biome colors of the surrounding area
func calculateBiome(bs *blocksSnapshot, x, z int, img *image.NRGBA) (byte, byte, byte) {
	count := 0
	var r, g, b int
	for xx := -2; xx <= 2; xx++ {
		for zz := -2; zz <= 2; zz++ {
			biome := bs.biome(x+xx, z+zz)
			ix := biome.ColorIndex & 0xFF
			iy := biome.ColorIndex >> 8
			col := img.NRGBAAt(ix, iy)
			r += int(col.R)
			g += int(col.G)
			b += int(col.B)
			count++
		}
	}
	return byte(r / count), byte(g / count), byte(b / count)
}

func calculateLight(bs *blocksSnapshot, origX, origY, origZ int,
	x, y, z float64, face int, smooth, force bool) (byte, byte) {
	blockLight := bs.blockLight(origX, origY, origZ)
	skyLight := bs.skyLight(origX, origY, origZ)
	if !smooth {
		return blockLight, skyLight
	}
	count := 1

	// TODO(Think) Document/cleanup this
	// it was taken from and older renderer of mine
	// (thinkmap).

	var pox, poy, poz, nox, noy, noz int

	switch face {
	case 0: // Up
		poz, pox = 0, 0
		noz, nox = -1, -1
		poy = 1
		noy = 0
	case 1: // Down
		poz, pox = 0, 0
		noz, nox = -1, -1
		poy = -1
		noy = -2
	case 2: // North
		poy, pox = 0, 0
		noy, nox = -1, -1
		poz = -1
		noz = -2
	case 3: // South
		poy, pox = 0, 0
		noy, nox = -1, -1
		poz = 1
		noz = 0
	case 4: // West
		poz, poy = 0, 0
		noz, noy = -1, -1
		pox = -1
		nox = -2
	case 5: // East
		poz, poy = 0, 0
		noz, noy = -1, -1
		pox = 1
		nox = 0
	}
	for ox := nox; ox <= pox; ox++ {
		for oy := noy; oy <= poy; oy++ {
			for oz := noz; oz <= poz; oz++ {
				bx := round(x + float64(ox))
				by := round(y + float64(oy))
				bz := round(z + float64(oz))
				count++
				blockLight += bs.blockLight(bx, by, bz)
				if !force {
					skyLight += bs.skyLight(bx, by, bz)
				} else if bl := bs.block(bx, by, bz); bl.Is(Blocks.Air) {
					skyLight += 15
				}
			}
		}

	}

	return blockLight / byte(count), skyLight / byte(count)
}

func round(f float64) int {
	return int(f + math.Copysign(0.5, f))
}

// Precomputed face vertices
var faceVertices = [6][6]chunkVertex{
	{ // Up
		{X: 0, Y: 1, Z: 0, TOffsetX: 0, TOffsetY: 0},
		{X: 1, Y: 1, Z: 0, TOffsetX: 1, TOffsetY: 0},
		{X: 0, Y: 1, Z: 1, TOffsetX: 0, TOffsetY: 1},

		{X: 1, Y: 1, Z: 1, TOffsetX: 1, TOffsetY: 1},
		{X: 0, Y: 1, Z: 1, TOffsetX: 0, TOffsetY: 1},
		{X: 1, Y: 1, Z: 0, TOffsetX: 1, TOffsetY: 0},
	},
	{ // Down
		{X: 0, Y: 0, Z: 0, TOffsetX: 0, TOffsetY: 1},
		{X: 0, Y: 0, Z: 1, TOffsetX: 0, TOffsetY: 0},
		{X: 1, Y: 0, Z: 0, TOffsetX: 1, TOffsetY: 1},

		{X: 1, Y: 0, Z: 1, TOffsetX: 1, TOffsetY: 0},
		{X: 1, Y: 0, Z: 0, TOffsetX: 1, TOffsetY: 1},
		{X: 0, Y: 0, Z: 1, TOffsetX: 0, TOffsetY: 0},
	},
	{ // North
		{X: 0, Y: 0, Z: 0, TOffsetX: 1, TOffsetY: 1},
		{X: 1, Y: 0, Z: 0, TOffsetX: 0, TOffsetY: 1},
		{X: 0, Y: 1, Z: 0, TOffsetX: 1, TOffsetY: 0},

		{X: 1, Y: 1, Z: 0, TOffsetX: 0, TOffsetY: 0},
		{X: 0, Y: 1, Z: 0, TOffsetX: 1, TOffsetY: 0},
		{X: 1, Y: 0, Z: 0, TOffsetX: 0, TOffsetY: 1},
	},
	{ // South
		{X: 0, Y: 0, Z: 1, TOffsetX: 0, TOffsetY: 1},
		{X: 0, Y: 1, Z: 1, TOffsetX: 0, TOffsetY: 0},
		{X: 1, Y: 0, Z: 1, TOffsetX: 1, TOffsetY: 1},

		{X: 1, Y: 1, Z: 1, TOffsetX: 1, TOffsetY: 0},
		{X: 1, Y: 0, Z: 1, TOffsetX: 1, TOffsetY: 1},
		{X: 0, Y: 1, Z: 1, TOffsetX: 0, TOffsetY: 0},
	},
	{ // West
		{X: 0, Y: 0, Z: 0, TOffsetX: 0, TOffsetY: 1},
		{X: 0, Y: 1, Z: 0, TOffsetX: 0, TOffsetY: 0},
		{X: 0, Y: 0, Z: 1, TOffsetX: 1, TOffsetY: 1},

		{X: 0, Y: 1, Z: 1, TOffsetX: 1, TOffsetY: 0},
		{X: 0, Y: 0, Z: 1, TOffsetX: 1, TOffsetY: 1},
		{X: 0, Y: 1, Z: 0, TOffsetX: 0, TOffsetY: 0},
	},
	{ // East
		{X: 1, Y: 0, Z: 0, TOffsetX: 1, TOffsetY: 1},
		{X: 1, Y: 0, Z: 1, TOffsetX: 0, TOffsetY: 1},
		{X: 1, Y: 1, Z: 0, TOffsetX: 1, TOffsetY: 0},

		{X: 1, Y: 1, Z: 1, TOffsetX: 0, TOffsetY: 0},
		{X: 1, Y: 1, Z: 0, TOffsetX: 1, TOffsetY: 0},
		{X: 1, Y: 0, Z: 1, TOffsetX: 0, TOffsetY: 1},
	},
}
