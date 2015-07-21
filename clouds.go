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
	"bytes"
	"image"
	"image/png"
	"io/ioutil"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/resource"
)

var (
	cloudOffset float64
	cloudIndex  int
	clouds      []*cloud
	cloudImage  *image.NRGBA
)

type cloud struct {
	*render.StaticModel

	x, y int
}

func tickClouds(delta float64) {
	if Client != nil && Client.WorldType != wtOverworld {
		for _, c := range clouds {
			c.Free()
		}
		clouds = nil
		return
	}
	if cloudImage == nil {
		f, err := resource.Open("minecraft", "textures/environment/clouds.png")
		if err != nil {
			cloudImage = image.NewNRGBA(image.Rect(0, 0, 256, 256))
		} else {
			defer f.Close()
			img, err := png.Decode(f)
			if err != nil {
				panic(err)
			}
			i, ok := img.(*image.NRGBA)
			if !ok {
				i = convertImage(img)
			}
			cloudImage = i
		}

		var buf bytes.Buffer
		png.Encode(&buf, cloudImage)
		ioutil.WriteFile("test.png", buf.Bytes(), 0777)
	}
	if clouds == nil {
		clouds = make([]*cloud, 25*25)
		for i := range clouds {
			tex := render.GetTexture("environment/clouds")
			data := appendBox(nil, -6, -2, -6, 12, 4, 12, [6]render.TextureInfo{tex, tex, tex, tex, tex, tex})
			clouds[i] = &cloud{StaticModel: render.NewStaticModel([][]*render.StaticVertex{data})}
			clouds[i].Colors[0] = [4]float32{1.0, 1.0, 1.0, 1.0}
		}
	}
	cloudIndex = 0

	cloudOffset += delta

	for x := -12; x <= 12; x++ {
		for y := -12; y <= 12; y++ {
			fx, fy := float64(x)/256.0, float64(y)/256.0
			fx += -math.Floor(Client.X/12.0) / 256.0
			fy += -math.Floor(Client.Z/12.0) / 256.0
			fy += cloudOffset / 500.0 / 256
			c := getCloud(
				math.Mod(1+math.Mod(fx, 1), 1),
				math.Mod(1+math.Mod(fy, 1), 1),
			)
			if c == nil {
				continue
			}
			c.Y = -128
			c.X = -float32(math.Floor((Client.X-float64(x*12))/12) * 12)
			c.Z = float32(math.Floor((Client.Z-float64(y*12))/12)*12) + float32(math.Mod(cloudOffset/500.0, 1)*12)
			c.Radius = 20
			c.Matrix[0] = mgl32.Translate3D(-c.X, c.Y, c.Z)
			c.SkyLight = 15
		}
	}

	for i := cloudIndex; i < len(clouds); i++ {
		c := clouds[i]
		c.X = 0
		c.Z = 0
		c.Y = 9999
		c.Radius = 0.01
	}
}

func getCloud(x, y float64) *cloud {
	px, py := int(256*x)%255, int(256*y)%255

	sx := cloudImage.Bounds().Dx() / 256
	sy := cloudImage.Bounds().Dy() / 256
	var ok bool
check:
	for xx := 0; xx < sx; xx++ {
		for yy := 0; yy < sy; yy++ {
			col := cloudImage.NRGBAAt(xx+px, yy+py)
			if col.A > 20 {
				ok = true
				break check
			}
		}
	}
	if !ok {
		return nil
	}
	c := clouds[cloudIndex]
	cloudIndex++
	if c.x != px || c.y != py {
		tex := render.RelativeTexture(render.GetTexture("environment/clouds"), 256, 256).
			Sub(px, py, 1, 1)
		for _, v := range c.Verts {
			v.Texture = tex
		}
		c.Refresh()
		c.x = px
		c.y = py
	}
	return c
}
