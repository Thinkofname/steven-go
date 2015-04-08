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
	"image/png"
	"io"

	"github.com/thinkofdeath/steven/resource"
)

var (
	fontPages         [0x100]*TextureInfo
	isFontLoaded      [0x100]bool
	fontCharacterInfo [0x10000]fontInfo
)

type fontInfo struct {
	Start, End int
}

type UIText struct {
	elements []*UIElement
	Width    int
}

func AddUIText(str string, x, y float64, rr, gg, bb int) *UIText {
	t := &UIText{}
	offset := 0.0
	for _, r := range str {
		if r == ' ' {
			offset += 6
			continue
		}
		page := int(r >> 8)
		if !isFontLoaded[page] {
			loadFontPage(page)
		}
		p := fontPages[page]
		if p == nil {
			continue
		}
		c := int(r & 0xFF)
		var w float64
		var tx, ty, tw, th int
		cx, cy := c&0xF, c>>4
		info := fontCharacterInfo[r]
		if page == 0 {
			tx = cx*8 + info.Start
			tw = info.End - info.Start
			ty = cy * 8
			th = 8
			w = float64(tw * 2)
		} else {
			tx = cx*16 + info.Start
			tw = info.End - info.Start
			ty = cy * 16
			th = 16
			w = float64(tw)
		}
		shadow := AddUIElement(p, x+offset+2, y+2, w, 16, tx, ty, tw, th)
		shadow.R = byte(float64(rr) * 0.2)
		shadow.G = byte(float64(gg) * 0.2)
		shadow.B = byte(float64(bb) * 0.2)
		t.elements = append(t.elements, shadow)
		text := AddUIElement(p, x+offset, y, w, 16, tx, ty, tw, th)
		text.R = byte(rr)
		text.G = byte(gg)
		text.B = byte(bb)
		t.elements = append(t.elements, text)
		offset += w + 2
	}
	t.Width = int(offset - 2)
	return t
}

func (u *UIText) Free() {
	for _, e := range u.elements {
		e.Free()
	}
}

func loadFontPage(page int) {
	textureLock.Lock()
	defer textureLock.Unlock()
	isFontLoaded[page] = true
	var p string
	if page == 0 {
		p = "ascii"
	} else {
		p = fmt.Sprintf("unicode_page_%02x", page)
	}
	r, err := resource.Open("minecraft", "textures/font/"+p+".png")
	if err != nil {
		return
	}
	defer r.Close()
	img, err := png.Decode(r)
	if err != nil {
		panic(err)
	}
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	pix := imgToBytes(img)
	at, rect := addTexture(pix, width, height)
	info := &TextureInfo{
		Rect:  rect,
		Atlas: at,
	}
	fontPages[page] = info
	if p == "ascii" {
		calculateFontSizes(img)
	}
}

func calculateFontSizes(img image.Image) {
	for i := 0; i <= 255; i++ {
		cx := (i & 0xF) * 8
		cy := (i >> 4) * 8
		start := true
	xloop:
		for x := 0; x < 8; x++ {
			for y := 0; y < 8; y++ {
				_, _, _, a := img.At(cx+x, cy+y).RGBA()
				if start && a != 0 {
					fontCharacterInfo[i].Start = x
					start = false
					continue xloop
				} else if !start && a != 0 {
					continue xloop
				}
			}
			fontCharacterInfo[i].End = x
			break
		}
	}
}

func loadFontInfo() {
	r, err := resource.Open("minecraft", "font/glyph_sizes.bin")
	if err != nil {
		panic(err)
	}
	var data [0x10000]byte
	_, err = io.ReadFull(r, data[:])
	if err != nil {
		panic(err)
	}
	for i := range fontCharacterInfo {
		fontCharacterInfo[i].Start = int(data[i] >> 4)
		fontCharacterInfo[i].End = int(data[i]&0xF) + 1
	}
}
