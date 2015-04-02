package main

import (
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/render/builder"
	"github.com/thinkofdeath/steven/type/direction"
)

type chunkVertex struct {
	X, Y, Z              int16
	TX, TY, TW, TH       uint16
	TOffsetX, TOffsetY   int16
	R, G, B              byte
	BlockLight, SkyLight byte
}

type buildPos struct {
	X, Y, Z int
}

var _, chunkVertexType = builder.Struct(chunkVertex{})

func (cs *chunkSection) build(complete chan<- buildPos) {
	ox, oy, oz := (cs.chunk.X<<4)-2, (cs.Y<<4)-2, (cs.chunk.Z<<4)-2
	bs := getSnapshot(ox, oy, oz, 20, 20, 20)
	// Make relative
	bs.x = -2
	bs.y = -2
	bs.z = -2
	go func() {
		bO := builder.New(chunkVertexType...)
		bT := builder.New(chunkVertexType...)

		for y := 0; y < 16; y++ {
			for x := 0; x < 16; x++ {
				for z := 0; z < 16; z++ {
					bl := bs.block(x, y, z)
					// Air will never have a model
					if bl.Is(BlockAir) {
						continue
					}
					b := bO
					// Translucent models need special handling
					if bl.IsTranslucent() {
						b = bT
					}

					// Liquid can't be represented by the model system
					// due to the number of possible states they have
					if l, ok := bl.(*blockLiquid); ok {
						for _, v := range l.renderLiquid(bs, x, y, z) {
							buildVertex(b, v)
						}
						continue
					}

					// The seed is used to select a 'random' variant which is
					// constant for that position.
					seed := (cs.chunk.X<<4 + x) ^ (cs.Y<<4+y)*31 ^ (cs.chunk.Z<<4+z)*5

					if variant := bl.Models().selectModel(seed); variant != nil {
						for _, v := range variant.render(x, y, z, bs) {
							buildVertex(b, v)
						}
						continue
					}
				}
			}
		}

		// Update culling information
		cullBits := buildCullBits(bs)

		// Upload the buffers on the render goroutine
		render.Sync(func() {
			cs.Buffer.Upload(bO.Data(), bO.Count(), cullBits)
			cs.Buffer.UploadTrans(bT.Data(), bT.Count())
		})
		// Free up the builder
		complete <- buildPos{cs.chunk.X, cs.Y, cs.chunk.Z}
	}()
}

func buildCullBits(bs *blocksSnapshot) uint64 {
	bits := uint64(0)
	set := func(from, to direction.Type) {
		bits |= 1 << (from*6 + to)
	}

	visited := map[position]struct{}{}
	// This tries a flood fill on every block in the chunk
	// section with an optimization of not visiting a block
	// that was visited in a previous fill (as it would already
	// be accounted for).
	for y := 0; y < 16; y++ {
		for z := 0; z < 16; z++ {
			for x := 0; x < 16; x++ {
				if _, ok := visited[position{x, y, z}]; ok {
					continue
				}
				touched := floodFill(bs, visited, x, y, z)
				// Minor optimization for a common case
				if touched == 0 {
					continue
				}
				// Mark each face in the set as visible through
				// each other
				for _, d := range direction.Values {
					if touched&(1<<d) != 0 {
						for _, d2 := range direction.Values {
							if touched&(1<<d2) != 0 {
								set(d, d2)
							}
						}
					}
				}
			}
		}
	}

	return bits
}

func floodFill(bs *blocksSnapshot, visited map[position]struct{}, x, y, z int) uint8 {
	pos := position{x, y, z}
	// Make sure we aren't filling the same spot repeatedly or
	// going out of bounds.
	if _, ok := visited[pos]; ok || x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return 0
	}
	visited[pos] = struct{}{}

	// Can't fill into 'solid' spaces (ones that completely fill
	// the block)
	if bs.block(x, y, z).ShouldCullAgainst() {
		return 0
	}

	// bits are used to represent touched faces
	touched := uint8(0)
	if x == 0 {
		touched |= 1 << direction.West
	} else if x == 15 {
		touched |= 1 << direction.East
	}
	if y == 0 {
		touched |= 1 << direction.Down
	} else if y == 15 {
		touched |= 1 << direction.Up
	}
	if z == 0 {
		touched |= 1 << direction.North
	} else if z == 15 {
		touched |= 1 << direction.South
	}

	// Fill around us and add the touched faces to our
	// bits
	for _, d := range direction.Values {
		ox, oy, oz := d.Offset()
		touched |= floodFill(bs, visited, x+ox, y+oy, z+oz)
	}

	return touched
}

type position struct {
	X, Y, Z int
}

// builder.Struct works by reflection which is to slow for this
// as its called so often.
func buildVertex(b *builder.Buffer, v chunkVertex) {
	b.Short(v.X)
	b.Short(v.Y)
	b.Short(v.Z)
	b.UnsignedShort(v.TX)
	b.UnsignedShort(v.TY)
	b.UnsignedShort(v.TW)
	b.UnsignedShort(v.TH)
	b.Short(v.TOffsetX)
	b.Short(v.TOffsetY)
	b.UnsignedByte(v.R)
	b.UnsignedByte(v.G)
	b.UnsignedByte(v.B)
	b.UnsignedByte(v.BlockLight)
	b.UnsignedByte(v.SkyLight)
}
