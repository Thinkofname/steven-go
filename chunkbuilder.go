package main

import (
	"fmt"
	"sync"

	"github.com/thinkofdeath/steven/render/builder"
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
		b := builder.New(chunkVertexType...)

		for y := 0; y < 16; y++ {
			for x := 0; x < 16; x++ {
				for z := 0; z < 16; z++ {
					bl := bs.block(x, y, z)
					if bl.Is(BlockAir) {
						continue
					}

					if l, ok := bl.(*blockLiquid); ok {
						for _, v := range l.renderLiquid(bs, x, y, z) {
							// chunkVertexF(b, v)
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
							// chunkVertexF(b, v)
							buildVertex(b, v)
						}
						continue
					}
				}
			}
		}

		cs.Buffer.Upload(b.Data(), b.Count())
		complete <- buildPos{cs.chunk.X, cs.Y, cs.chunk.Z}
	}()
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

var (
	warnedBlockModels = map[Block]struct{}{}
	warnLock          sync.RWMutex
)

func warnMissingModel(b Block) {
	warnLock.RLock()
	if _, ok := warnedBlockModels[b]; ok {
		warnLock.RUnlock()
		return
	}
	warnLock.RUnlock()
	warnLock.Lock()
	// Check again in case of another worker warning between switching
	// locks
	if _, ok := warnedBlockModels[b]; ok {
		warnLock.Unlock()
		return
	}
	fmt.Printf("Missing block model for %s\n", b)
	warnedBlockModels[b] = struct{}{}
	warnLock.Unlock()
}
