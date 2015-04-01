package main

import (
	"fmt"

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
					if bl.Is(BlockAir) {
						continue
					}
					b := bO
					if bl.IsTranslucent() {
						b = bT
					}

					if l, ok := bl.(*blockLiquid); ok {
						for _, v := range l.renderLiquid(bs, x, y, z) {
							buildVertex(b, v)
						}
						continue
					}

					if bl.Model() == nil {
						continue
					}

					seed := (cs.chunk.X<<4 + x) ^ (cs.Y + y) ^ (cs.chunk.Z<<4 + z)
					if variant := bl.Model().variant(bl.ModelVariant(), seed); variant != nil {
						for _, v := range variant.render(x, y, z, bs) {
							buildVertex(b, v)
						}
						continue
					} else {
						fmt.Printf("Missing variant %s for %s\n", bl.ModelVariant(), bl)
					}
				}
			}
		}

		cullBits := buildCullBits(bs)

		cs.Buffer.Upload(bO.Data(), bO.Count(), cullBits)
		cs.Buffer.UploadTrans(bT.Data(), bT.Count())
		complete <- buildPos{cs.chunk.X, cs.Y, cs.chunk.Z}
	}()
}

func buildCullBits(bs *blocksSnapshot) uint64 {
	bits := uint64(0)
	set := func(from, to direction.Type) {
		bits |= 1 << (from*6 + to)
	}

	visited := map[position]struct{}{}
	for y := 0; y < 16; y++ {
		for z := 0; z < 16; z++ {
			for x := 0; x < 16; x++ {
				if _, ok := visited[position{x, y, z}]; ok {
					continue
				}
				touched := floodFill(bs, visited, x, y, z)

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
	if _, ok := visited[pos]; ok || x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return 0
	}
	visited[pos] = struct{}{}

	if bs.block(x, y, z).ShouldCullAgainst() {
		return 0
	}

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

	for _, d := range direction.Values {
		ox, oy, oz := d.Offset()
		touched |= floodFill(bs, visited, x+ox, y+oy, z+oz)
	}

	return touched
}

type position struct {
	X, Y, Z int
}

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
