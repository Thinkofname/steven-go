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
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/thinkofdeath/steven/render"
)

var sunModel *render.Model

func tickClouds(delta float64) {
	if Client != nil && Client.WorldType != wtOverworld {
		render.DrawClouds = false
		if sunModel != nil {
			sunModel.Free()
			sunModel = nil
		}
		return
	}
	render.DrawClouds = true
	if Client == nil && Client.entity != nil {
		return
	}
	if sunModel == nil {
		genSunModel()
	}
	x, y, z := Client.entity.Position()

	time := Client.WorldTime / 12000
	ox := math.Cos(time*math.Pi) * 300
	oy := math.Sin(time*math.Pi) * 300

	sunModel.Matrix[0] = mgl32.Translate3D(
		float32(x+ox),
		-float32(y+oy),
		float32(z),
	).Mul4(mgl32.Rotate3DZ(-float32(time * math.Pi)).Mat4())
}

func genSunModel() {
	const size = 50
	tex := render.GetTexture("environment/sun")
	sunModel = render.NewModelCollection([][]*render.ModelVertex{
		{
			{X: 0, Y: -size, Z: -size, TextureX: 0, TextureY: 1, Texture: tex, R: 255, G: 255, B: 255, A: 255},
			{X: 0, Y: size, Z: -size, TextureX: 0, TextureY: 0, Texture: tex, R: 255, G: 255, B: 255, A: 255},
			{X: 0, Y: -size, Z: size, TextureX: 1, TextureY: 1, Texture: tex, R: 255, G: 255, B: 255, A: 255},
			{X: 0, Y: size, Z: size, TextureX: 1, TextureY: 0, Texture: tex, R: 255, G: 255, B: 255, A: 255},
		},
	}, render.SunModels)
}
