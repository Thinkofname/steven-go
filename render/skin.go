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

package render

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/thinkofdeath/steven/resource"
)

var (
	skins = map[string]*skin{}
)

type skin struct {
	info     *reusableTexture
	data     []byte
	refCount int
}

const skinCache = "skincache/"

var skinBuffer = make([]byte, 64*64*4)

func LoadSkinBuffer() {
	r, err := resource.Open("minecraft", "textures/entity/steve.png")
	if err != nil {
		skinBuffer = make([]byte, 64*64*4)
		return
	}
	defer r.Close()
	i, err := png.Decode(r)
	if err != nil {
		panic(err)
	}
	out := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	draw.Draw(out, out.Bounds(), i, image.ZP, draw.Over)
	skinBuffer = out.Pix
}

func RefSkin(hash string) {
	s := skins[hash]
	if s != nil {
		s.refCount++
		return
	}
	info := getFreeTexture(64, 64)
	uploadTexture(info.info, skinBuffer)
	s = &skin{
		info:     info,
		refCount: 1,
	}
	skins[hash] = s
	go obtainSkin(hash, s)
}

func obtainSkin(hash string, s *skin) {
	var r io.Reader
	var fromCache bool
	if f, err := os.Open(skinPath(hash)); err != nil {
		resp, err := http.Get(fmt.Sprintf("http://textures.minecraft.net/texture/%s", hash))
		if err != nil {
			fmt.Printf("Error downloading skin: %s\n", err)
			return
		}
		defer resp.Body.Close()
		r = resp.Body
	} else {
		defer f.Close()
		r = f
		fromCache = true
	}
	img, err := png.Decode(r)
	if err != nil {
		fmt.Printf("Error decoding skin: %s\n", err)
		return
	}
	if !fromCache {
		path := skinPath(hash)
		os.MkdirAll(filepath.Dir(path), 0777)
		f, err := os.Create(path)
		if err != nil {
			return
		}
		defer f.Close()
		err = png.Encode(f, img)
		if err != nil {
			panic(err)
		}
	}
	if img.Bounds().Dy() == 32 {
		ni := image.NewNRGBA(image.Rect(0, 0, 64, 64))
		draw.Draw(ni, img.Bounds(), img, image.ZP, draw.Over)
		draw.Draw(ni, image.Rect(16, 48, 32, 64), img, image.Pt(0, 16), draw.Over)
		draw.Draw(ni, image.Rect(32, 48, 48, 64), img, image.Pt(40, 16), draw.Over)
		img = ni
	}
	di := img.(draw.Image)
	for _, off := range [][4]int{
		// X, Y, W, H
		{0, 0, 32, 16},
		{16, 16, 24, 16},
		{0, 16, 16, 16},
		{16, 48, 16, 16},
		{32, 48, 16, 16},
		{40, 16, 16, 16},
	} {
		for x := off[0]; x < off[0]+off[2]; x++ {
			for y := off[1]; y < off[1]+off[3]; y++ {
				col := di.At(x, y)
				rgba := color.RGBAModel.Convert(col).(color.RGBA)
				di.Set(x, y, color.RGBA{R: rgba.R, B: rgba.B, G: rgba.G, A: 255})
			}
		}
	}
	pix := imgToBytes(img)
	Sync(func() {
		s.data = pix
		uploadTexture(s.info.info, pix)
	})
}

func skinPath(hash string) string {
	return filepath.Join(skinCache, hash[:2], hash+".png")
}

func Skin(hash string) TextureInfo {
	s := skins[hash]
	if s != nil {
		return s.info.info
	}
	return nil
}

func FreeSkin(hash string) {
	s := skins[hash]
	if s == nil {
		return
	}
	s.refCount--
	if s.refCount <= 0 {
		freeTexture(s.info)
		delete(skins, hash)
	}
}

// So I reuse the skins for icons
// hacky but it works

func FreeIcon(id string) {
	FreeSkin(id)
}

func Icon(id string) TextureInfo {
	return Skin(id)
}

func AddIcon(id string, pix image.Image) {
	s := skins[id]
	if s != nil {
		s.refCount++
		return
	}
	info := getFreeTexture(pix.Bounds().Dx(), pix.Bounds().Dy())
	data := imgToBytes(pix)
	uploadTexture(info.info, data)
	s = &skin{
		info:     info,
		refCount: 1,
		data:     data,
	}
	skins[id] = s
}
