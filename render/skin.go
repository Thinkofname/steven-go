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
	"image/draw"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/thinkofdeath/steven/resource"
)

var (
	skins            = map[string]*skin{}
	freeSkinTextures []*TextureInfo
)

type skin struct {
	info     *TextureInfo
	refCount int
}

const skinCache = "skincache/"

var skinBuffer []byte

func init() {
	r, err := resource.Open("minecraft", "textures/entity/steve.png")
	if err != nil {
		panic(err)
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
	info := getSkinInfo()
	uploadTexture(info, skinBuffer)
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
	pix := imgToBytes(img)
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
	Sync(func() {
		uploadTexture(s.info, pix)
	})
}

func skinPath(hash string) string {
	return filepath.Join(skinCache, hash[:2], hash+".png")
}

func Skin(hash string) *TextureInfo {
	s := skins[hash]
	if s != nil {
		return s.info
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
		freeSkinTextures = append(freeSkinTextures, s.info)
		delete(skins, hash)
	}
}

func getSkinInfo() *TextureInfo {
	if len(freeSkinTextures) == 0 {
		return addTexture(skinBuffer, 64, 64)
	}
	var info *TextureInfo
	l := len(freeSkinTextures)
	info, freeSkinTextures = freeSkinTextures[l-1], freeSkinTextures[:l-1]
	return info
}
