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
	"time"

	"github.com/thinkofdeath/steven/render"
)

var fakeGenDistance = 7

func fakeGen() {
	render.Camera.X = 0.5
	render.Camera.Z = 0.5
	render.Camera.Y = 70
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	liquid := Blocks.Water
	top := Blocks.Grass
	mid := Blocks.Dirt
	bot := Blocks.Stone
	ra := r.Intn(3)
	if ra == 1 {
		liquid = Blocks.Lava
		render.LightLevel = 0.9
		render.ClearColour.R, render.ClearColour.G, render.ClearColour.B = 52/255.0, 8/255.0, 8/255.0
		top = Blocks.Netherrack
		mid = Blocks.Netherrack
		bot = Blocks.SoulSand
	} else if ra == 2 {
		liquid = Blocks.Air
		render.LightLevel = 0.8
		render.ClearColour.R, render.ClearColour.G, render.ClearColour.B = 23/255.0, 0, 23/255.0
		top = Blocks.EndStone
		mid = Blocks.EndStone
		bot = Blocks.EndStone
	}
	go func() {
		randGrid := make([]int, (fakeGenDistance*2+1)*(fakeGenDistance*2+1))
		for i := range randGrid {
			randGrid[i] = r.Intn(10) + 54
		}
		get := func(cx, cz int) int {
			if cx < -fakeGenDistance || cz < -fakeGenDistance || cx > fakeGenDistance || cz > fakeGenDistance {
				return 63
			}
			cx += fakeGenDistance
			cz += fakeGenDistance
			return randGrid[cx+cz*(fakeGenDistance*2+1)]
		}
		smooth := func(cx, cz, x, y, z int) int {
			tl := float64(get(cx, cz))
			tr := float64(get(cx+1, cz))
			bl := float64(get(cx, cz+1))
			br := float64(get(cx+1, cz+1))
			t := tl*((15-float64(x))/15.0) + tr*(float64(x)/15.0)
			b := bl*((15-float64(x))/15.0) + br*(float64(x)/15.0)
			return int(t*((15-float64(z))/15.0) + b*(float64(z)/15.0))
		}

		for cx := -fakeGenDistance; cx <= fakeGenDistance; cx++ {
			for cz := -fakeGenDistance; cz <= fakeGenDistance; cz++ {
				c := &chunk{
					chunkPosition: chunkPosition{
						X: cx, Z: cz,
					},
				}

				for i := 0; i < 4; i++ {
					cs := newChunkSection(c, i)
					c.Sections[i] = cs
					for y := 0; y < 16; y++ {
						for z := 0; z < 16; z++ {
							for x := 0; x < 16; x++ {
								height := smooth(cx, cz, x, y, z)
								ry := y + i<<4
								var block Block
								cs.setSkyLight(0, x, y, z)
								switch {
								case ry <= height-5:
									block = bot.Base
								case ry <= height-1:
									block = mid.Base
								case ry == height:
									block = top.Base
								default:
									level := 0xF
									if ry >= 60 {
										block = Blocks.Air.Base
									} else {
										block = liquid.Base
										if liquid == Blocks.Water {
											level = 13 - (60-ry)*2
										}
									}
									if level < 0 {
										level = 0
									}
									sky := (16*16*16*2 + 16*16*8) * 4
									sky += 16 * 16 * 8 * i
									cs.setSkyLight(byte(level), x, y, z)
								}
								cs.setBlock(block, x, y, z)
							}
						}
					}
				}
				c.calcHeightmap()
				syncChan <- c.postLoad
			}
		}
	}()
}
