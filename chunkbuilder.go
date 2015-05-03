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
	"math/rand"
	"sync"

	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/render/builder"
	"github.com/thinkofdeath/steven/type/bit"
	"github.com/thinkofdeath/steven/type/direction"
)

type chunkVertex struct {
	X, Y, Z                    int16
	TX, TY, TW, TH             uint16
	TOffsetX, TOffsetY, TAtlas int16
	R, G, B                    byte
	BlockLight, SkyLight       uint16
}

type buildPos struct {
	X, Y, Z int
}

var (
	_, chunkVertexType = builder.Struct(chunkVertex{})
	builderPool        = sync.Pool{
		New: func() interface{} {
			return builder.New(chunkVertexType...)
		},
	}
)

func (cs *chunkSection) build(complete chan<- buildPos) {
	ox, oy, oz := (cs.chunk.X<<4)-2, (cs.Y<<4)-2, (cs.chunk.Z<<4)-2
	bs := getPooledSnapshot(ox, oy, oz)
	// Make relative
	bs.x = -2
	bs.y = -2
	bs.z = -2
	go func() {
		bO := builderPool.Get().(*builder.Buffer)
		bT := builderPool.Get().(*builder.Buffer)

		r := rand.New(rand.NewSource(int64(cs.chunk.X) | (int64(cs.chunk.Z) << 32)))
		var tInfo []render.ObjectInfo

		for y := 0; y < 16; y++ {
			for x := 0; x < 16; x++ {
				for z := 0; z < 16; z++ {

					bl := bs.block(x, y, z)
					if !bl.Renderable() {
						// Use one step of the rng so that
						// if a block is placed in an empty
						// location is variant doesn't change
						r.Int()
						continue
					}
					b := bO
					// Translucent models need special handling
					if bl.IsTranslucent() {
						b = bT
					}
					offset := len(b.Data())

					// Liquids can't be represented by the model system
					// due to the number of possible states they have
					if l, ok := bl.(*blockLiquid); ok {
						l.renderLiquid(bs, x, y, z, b)
						r.Int() // See the comment above for air
						count := len(b.Data()) - offset
						if bl.IsTranslucent() && count > 0 {
							tInfo = append(tInfo, render.ObjectInfo{
								X:      (cs.chunk.X << 4) + x,
								Y:      (cs.Y << 4) + y,
								Z:      (cs.chunk.Z << 4) + z,
								Offset: offset,
								Count:  count,
							})
						}
						continue
					}

					// The index is used to select a 'random' variant which is
					// constant for that position.
					index := r.Intn(len(bl.Models()))

					if variant := bl.Models().selectModel(index); variant != nil {
						variant.Render(x, y, z, bs, b)
						count := len(b.Data()) - offset
						if bl.IsTranslucent() && count > 0 {
							tInfo = append(tInfo, render.ObjectInfo{
								X:      (cs.chunk.X << 4) + x,
								Y:      (cs.Y << 4) + y,
								Z:      (cs.chunk.Z << 4) + z,
								Offset: offset,
								Count:  count,
							})
						}
					}
				}
			}
		}

		// Update culling information
		cullBits := buildCullBits(bs)
		snapshotPool.Put(bs)

		// Upload the buffers on the render goroutine
		render.Sync(func() {
			cs.Buffer.Upload(bO.Data(), bO.Count(), cullBits)
			cs.Buffer.UploadTrans(tInfo, bT.Data(), bT.Count())
			bO.Reset()
			bT.Reset()
			builderPool.Put(bO)
			builderPool.Put(bT)
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

	visited := bit.NewSet(16 * 16 * 16)
	// This tries a flood fill on every block in the chunk
	// section with an optimization of not visiting a block
	// that was visited in a previous fill (as it would already
	// be accounted for).
	for y := 0; y < 16; y++ {
		for z := 0; z < 16; z++ {
			for x := 0; x < 16; x++ {
				if visited.Get(x | (z << 4) | (y << 8)) {
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

func floodFill(bs *blocksSnapshot, visited bit.Set, x, y, z int) uint8 {
	i := x | (z << 4) | (y << 8)
	// Make sure we aren't filling the same spot repeatedly or
	// going out of bounds.
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 || visited.Get(i) {
		return 0
	}
	visited.Set(i, true)

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
	b.Short(v.TAtlas)
	b.UnsignedByte(v.R)
	b.UnsignedByte(v.G)
	b.UnsignedByte(v.B)
	b.UnsignedByte(255)
	b.UnsignedShort(v.BlockLight)
	b.UnsignedShort(v.SkyLight)
}
