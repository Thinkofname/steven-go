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
	fontPages               [0x100]*TextureInfo
	isFontLoaded            [0x100]bool
	fontCharacterInfo       [0x10000]fontInfo
	aPageWidth, aPageHeight float64
)

type fontInfo struct {
	Start, End int
}

// UIText is a collection of UI elements that make up a
// string of characters.
type UIText struct {
	elements []*UIElement
	Width    float64
}

// DrawUIText draws a UIText element to the screen with
// the passed text at the location. The text may be tinted too
func DrawUIText(str string, x, y float64, rr, gg, bb int) UIText {
	return DrawUITextScaled(str, x, y, 1.0, 1.0, rr, gg, bb)
}

// DrawUITextScaled draws a UIText element to the screen with
// the passed text at the location. The text may be tinted and/or
// scaled too
func DrawUITextScaled(str string, x, y float64, sx, sy float64, rr, gg, bb int) UIText {
	t := UIText{}
	offset := 0.0
	for _, r := range str {
		if r == ' ' {
			offset += 6
			continue
		}
		page := int(r >> 8)
		// Lazy loading to save memory
		if !isFontLoaded[page] {
			loadFontPage(page)
		}
		p := fontPages[page]
		// We don't have font pages for every character
		if p == nil {
			continue
		}
		c := int(r & 0xFF)
		var w float64
		var tx, ty, tw, th float64
		cx, cy := c&0xF, c>>4
		info := fontCharacterInfo[r]
		if page == 0 {
			sw, sh := aPageWidth/16, aPageHeight/16
			// The first page is 128x128 instead of 256x256
			tx = float64(cx*int(sw)+info.Start) / aPageWidth
			tw = float64(info.End-info.Start) / aPageWidth
			ty = float64(cy*int(sh)) / aPageHeight
			th = sh / aPageHeight
			w = (float64(info.End-info.Start) / sw) * 16
		} else {
			tx = float64(cx*16+info.Start) / 256.0
			tw = float64(info.End-info.Start) / 256.0
			ty = float64(cy*16) / 256.0
			th = 16.0 / 256.0
			w = float64(info.End - info.Start)
		}
		shadow := DrawUIElement(p, x+(offset+2)*sx, y+2*sy, w*sx, 16*sy, tx, ty, tw, th)
		// Tint the shadow to a darker shade of the original color
		shadow.R = byte(float64(rr) * 0.25)
		shadow.G = byte(float64(gg) * 0.25)
		shadow.B = byte(float64(bb) * 0.25)
		t.elements = append(t.elements, shadow)
		text := DrawUIElement(p, x+offset*sx, y, w*sx, 16*sy, tx, ty, tw, th)
		text.R = byte(rr)
		text.G = byte(gg)
		text.B = byte(bb)
		t.elements = append(t.elements, text)
		offset += w + 2
	}
	t.Width = (offset - 2) * sx
	return t
}

// Returns the size of the passed character in pixels.
func SizeOfCharacter(r rune) float64 {
	if r == ' ' {
		return 4
	}
	info := fontCharacterInfo[r]
	if r>>8 == 0 {
		sw := aPageWidth / 16
		return (float64(info.End-info.Start) / sw) * 16
	}
	return float64(info.End - info.Start)
}

// Returns the size of the passed string in pixels.
func SizeOfString(str string) float64 {
	size := 0.0
	for _, r := range str {
		size += SizeOfCharacter(r) + 2
	}
	return size - 2
}

// Shift moves all the elements belonging to this UIText by the
// passed amounts.
func (u UIText) Shift(x, y float64) {
	for _, e := range u.elements {
		e.Shift(x, y)
	}
}

// Alpha changes the alpha of all theelements belonging to this UIText
func (u UIText) Alpha(a float64) {
	for _, e := range u.elements {
		e.Alpha(a)
	}
}

func loadFontPage(page int) {
	textureLock.Lock()
	defer textureLock.Unlock()
	isFontLoaded[page] = true
	var p string
	if page == 0 {
		// The ascii font is the minecraft style one
		// which is the default page 0 instead of the
		// unicode one for the english locales.
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
	info := addTexture(pix, width, height)
	fontPages[page] = info
	if p == "ascii" {
		// The font map file included with minecraft has the
		// wide of the unicode page 0 instead of the ascii one
		// we need to work this out ourselves
		calculateFontSizes(img)
	}
}

// Scans through each character computing the sizes
func calculateFontSizes(img image.Image) {
	aPageWidth, aPageHeight = float64(img.Bounds().Dx()), float64(img.Bounds().Dy())
	sw := img.Bounds().Dx() / 16
	sh := img.Bounds().Dy() / 16
	for i := 0; i <= 255; i++ {
		cx := (i & 0xF) * sw
		cy := (i >> 4) * sh
		start := true
	xloop:
		for x := 0; x < sw; x++ {
			for y := 0; y < sh; y++ {
				_, _, _, a := img.At(cx+x, cy+y).RGBA()
				if start && a != 0 {
					fontCharacterInfo[i].Start = x
					start = false
					continue xloop
				} else if !start && a != 0 {
					continue xloop
				}
			}
			if !start {
				fontCharacterInfo[i].End = x
				break
			}
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
		// Top nibble - start position
		// Bottom nibble - end position
		fontCharacterInfo[i].Start = int(data[i] >> 4)
		fontCharacterInfo[i].End = int(data[i]&0xF) + 1
	}
}
