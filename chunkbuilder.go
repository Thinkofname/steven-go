package main

import (
	"math"

	"github.com/thinkofdeath/steven/render/builder"
)

type chunkVertex struct {
	X, Y, Z int16
}

type buildPos struct {
	X, Y, Z int
}

var chunkVertexF, chunkVertexType = builder.Struct(chunkVertex{})

func (cs *chunkSection) build(complete chan<- buildPos) {
	ox, oy, oz := (cs.chunk.X<<4)-1, (cs.Y<<4)-1, (cs.chunk.Z<<4)-1
	bs := getSnapshot(ox, oy, oz, 18, 18, 18)
	go func() {
		b := builder.New(chunkVertexType...)

		block := func(x, y, z int) uint16 {
			return bs.block(ox+1+x, oy+1+y, oz+1+z)
		}

		for y := 0; y < 16; y++ {
			for x := 0; x < 16; x++ {
				for z := 0; z < 16; z++ {
					bl := block(x, y, z)
					if bl < 16 {
						continue
					}

					// Shitty test code

					if block(x, y+1, z) < 16 {
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z + 1) * 256),
						})

						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z + 1) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z + 1) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
						})
					}

					if block(x, y-1, z) < 16 {
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y) * 256),
							Z: int16((z) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z) * 256),
						})

						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
						})
					}

					if block(x-1, y, z) < 16 {
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y) * 256),
							Z: int16((z) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
						})

						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z + 1) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
						})
					}

					if block(x+1, y, z) < 16 {
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
						})

						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z + 1) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
						})
					}

					if block(x, y, z-1) < 16 {
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y) * 256),
							Z: int16((z) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
						})

						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z) * 256),
						})
					}

					if block(x, y, z+1) < 16 {
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z + 1) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
						})

						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z + 1) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z + 1) * 256),
						})
					}
				}
			}
		}

		cs.Buffer.Upload(b.Data(), b.Count())
		complete <- buildPos{cs.chunk.X, cs.Y, cs.chunk.Z}
	}()
}

type blocksSnapshot struct {
	Blocks []uint16

	x, y, z int
	w, h, d int
}

func getSnapshot(x, y, z, w, h, d int) *blocksSnapshot {
	bs := &blocksSnapshot{
		Blocks: make([]uint16, w*h*d),
		x:      x,
		y:      y,
		z:      z,
		w:      w,
		h:      h,
		d:      d,
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
							bs.setBlock(xx+(cx<<4), yy+(cy<<4), zz+(cz<<4), bl)
						}
					}
				}
			}
		}
	}

	return bs
}

func (bs *blocksSnapshot) block(x, y, z int) uint16 {
	return bs.Blocks[bs.index(x, y, z)]
}

func (bs *blocksSnapshot) setBlock(x, y, z int, b uint16) {
	bs.Blocks[bs.index(x, y, z)] = b
}

func (bs *blocksSnapshot) index(x, y, z int) int {
	x -= bs.x
	y -= bs.y
	z -= bs.z
	return x + z*bs.w + y*bs.w*bs.d
}
