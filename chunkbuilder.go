package main

import (
	"fmt"
	"sync"

	"github.com/thinkofdeath/steven/render/builder"
)

type chunkVertex struct {
	X, Y, Z            int16
	TX, TY, TW, TH     uint16
	TOffsetX, TOffsetY int16
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

		get := func(x, y, z int) Block {
			return bs.block(ox+1+x, oy+1+y, oz+1+z)
		}

		for y := 0; y < 16; y++ {
			for x := 0; x < 16; x++ {
				for z := 0; z < 16; z++ {
					bl := get(x, y, z)
					if bl.Is(BlockAir) {
						continue
					}

					if model := findStateModel(bl.Plugin(), bl.ModelName()); model != nil {
						seed := (cs.chunk.X<<4 + x) ^ (cs.Y + y) ^ (cs.chunk.Z<<4 + z)
						if variant := model.variant(bl.ModelVariant(), seed); variant != nil {
							for _, v := range variant.render(x, y, z, get) {
								chunkVertexF(b, v)
							}
							continue
						}
					}
					warnMissingModel(bl)

					// TODO: Remove
					if model := findStateModel("steven", "missing_block"); model != nil {
						seed := (cs.chunk.X<<4 + x) ^ (cs.Y + y) ^ (cs.chunk.Z<<4 + z)
						if variant := model.variant("normal", seed); variant != nil {
							for _, v := range variant.render(x, y, z, get) {
								chunkVertexF(b, v)
							}
							continue
						}
					}
				}
			}
		}

		cs.Buffer.Upload(b.Data(), b.Count())
		complete <- buildPos{cs.chunk.X, cs.Y, cs.chunk.Z}
	}()
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
