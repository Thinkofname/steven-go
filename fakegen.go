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
	"encoding/binary"
	"math/rand"
	"time"

	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/type/nibble"
)

var fakeGenDistance = 7

func fakeGen() {
	render.Camera.X = 0.5
	render.Camera.Z = 0.5
	render.Camera.Y = 70

	go func() {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
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
		liquid := uint16(9)
		if r.Float64() < 0.1 {
			liquid = 11
		}

		for cx := -fakeGenDistance; cx <= fakeGenDistance; cx++ {
			for cz := -fakeGenDistance; cz <= fakeGenDistance; cz++ {
				mask := uint16(0xF)
				data := make([]byte, 16*16*16*3*4+256)
				for i := 0; i < 4; i++ {
					for y := 0; y < 16; y++ {
						for z := 0; z < 16; z++ {
							for x := 0; x < 16; x++ {
								height := smooth(cx, cz, x, y, z)
								idx := 16 * 16 * 16 * 2 * i
								idx += (x | (z << 4) | (y << 8)) * 2
								ry := y + i<<4
								var val uint16
								switch {
								case ry <= height-5:
									val = (1 << 4)
								case ry <= height-1:
									val = (3 << 4)
								case ry == height:
									val = (2 << 4)
								default:
									level := 0xF
									if ry >= 60 {
										val = (0 << 4)
									} else {
										val = (liquid << 4)
										if liquid == 9 {
											level = 13 - (60-ry)*2
										}
									}
									if level < 0 {
										level = 0
									}
									sky := (16*16*16*2 + 16*16*8) * 4
									sky += 16 * 16 * 8 * i
									nibble.Array(data[sky:]).Set(x|(z<<4)|(y<<8), byte(level))
								}
								binary.LittleEndian.PutUint16(data[idx:], val)
							}
						}
					}
				}
				cx, cz := cx, cz
				syncChan <- func() { loadChunk(cx, cz, data, mask, true, true) }
			}
		}
	}()
}
