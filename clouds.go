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
	"image/png"
	"math"

	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/resource"
)

var (
	cloudOffset float64
	cloudImage  *image.Gray

	cloudLastX, cloudLastZ int = math.MaxInt32, math.MaxInt32
)

func tickClouds(delta float64) {
	if Client != nil && Client.WorldType != wtOverworld {
		render.DrawClouds = false
		return
	}
	render.DrawClouds = true
	if cloudImage == nil {
		f, err := resource.Open("steven", "textures/environment/clouds.png")
		if err != nil {
			panic(err)
		} else {
			defer f.Close()
			img, err := png.Decode(f)
			if err != nil {
				panic(err)
			}
			cloudImage = img.(*image.Gray)
		}
	}
	cloudOffset += delta
	data := render.CloudData()
	cx, cz := int(render.Camera.X), int(render.Camera.Z)

	clx, clz := pmod(cx, 512), pmod(cz+int(cloudOffset/60.0), 512)
	if clx == cloudLastX && clz == cloudLastZ {
		return
	}
	cloudLastX, cloudLastZ = clx, clz

	for x := 0; x < 360; x++ {
		for z := 0; z < 360; z++ {
			col := cloudImage.Pix[cloudImage.PixOffset(
				(clx+x)%512,
				(clz+z)%512,
			)]
			h := chunkMap.HighestBlockAt(cx+(x-160), cz+(z-160))

			data[x+z*512] = 0
			if h < 127 {
				if int(col)+h > 250 {
					data[x+z*512] = 0xFF
				}
			}
		}
	}
}

func pmod(x, y int) int {
	return (y + (x % y)) % y
}
