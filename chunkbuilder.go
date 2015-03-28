package main

import (
	"fmt"
	"sync"

	"github.com/thinkofdeath/steven/render/builder"
)

type chunkVertex struct {
	X, Y, Z int16
	R, G, B byte
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

					if model := findStateModel(bl.ModelName()); model != nil {
						// model.render(x, y, z, bs)
						continue
					}
					warnMissingModel(bl)

					bb, gg, rr := byte(bl.Color()&0xFF), byte((bl.Color()>>8)&0xFF), byte((bl.Color()>>16)&0xFF)

					// Shitty test code

					if get(x, y+1, z).Is(BlockAir) {
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})

						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})
					}

					if get(x, y-1, z).Is(BlockAir) {
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})

						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})
					}

					if get(x-1, y, z).Is(BlockAir) {
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})

						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})
					}

					if get(x+1, y, z).Is(BlockAir) {
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})

						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})
					}

					if get(x, y, z-1).Is(BlockAir) {
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})

						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z) * 256),
							R: rr, G: gg, B: bb,
						})
					}

					if get(x, y, z+1).Is(BlockAir) {
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})

						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x + 1) * 256),
							Y: int16((y) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})
						chunkVertexF(b, chunkVertex{
							X: int16((x) * 256),
							Y: int16((y + 1) * 256),
							Z: int16((z + 1) * 256),
							R: rr, G: gg, B: bb,
						})
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
