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
	"fmt"
	"image"
	icolor "image/color"
	"image/png"

	"github.com/thinkofdeath/steven/resource"
)

var (
	grassBiomeColors   *image.NRGBA
	foliageBiomeColors *image.NRGBA
)

func init() {
	grassBiomeColors = loadBiomeColors("grass")
	foliageBiomeColors = loadBiomeColors("foliage")
}

func loadBiomeColors(name string) *image.NRGBA {
	f, _ := resource.Open("minecraft", fmt.Sprintf("textures/colormap/%s.png", name))
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	i, ok := img.(*image.NRGBA)
	if !ok {
		i = convertImage(img)
	}
	return i
}

func convertImage(img image.Image) *image.NRGBA {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	pix := make([]byte, width*height*4)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			col := icolor.NRGBAModel.Convert(img.At(x, y)).(icolor.NRGBA)
			index := (y*width + x) * 4
			pix[index] = col.R
			pix[index+1] = col.G
			pix[index+2] = col.B
			pix[index+3] = col.A
		}
	}
	i := image.NewNRGBA(img.Bounds())
	i.Pix = pix
	return i
}
