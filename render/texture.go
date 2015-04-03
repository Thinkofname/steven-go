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
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"sync"

	"github.com/thinkofdeath/steven/platform/gl"
	"github.com/thinkofdeath/steven/render/atlas"
	"github.com/thinkofdeath/steven/resource"
)

var (
	textures         []*atlas.Type
	textureViews     []*image.NRGBA
	textureMap       = map[string]TextureInfo{}
	textureLock      sync.RWMutex
	animatedTextures []*animatedTexture
)

const atlasSize = 1024

// TextureInfo returns information about a texture in an atlas
type TextureInfo struct {
	Atlas     int
	imageView *image.NRGBA
	*atlas.Rect
}

// GetTexture returns the related TextureInfo for the requested texture.
// If the texture isn't found a placeholder is returned instead.
func GetTexture(name string) *TextureInfo {
	textureLock.RLock()
	defer textureLock.RUnlock()
	t, ok := textureMap[name]
	if !ok {
		t = textureMap["missing_texture"]
	}
	return &t
}

// LoadTextures (re)loads all the block textures from the resource pack(s)
// TODO(Think) better error handling (if possible to recover?)
func LoadTextures() {
	textureLock.Lock()
	defer textureLock.Unlock()

	// Clear existing
	textures = nil
	textureMap = map[string]TextureInfo{}

	for _, file := range resource.Search("minecraft", "textures/blocks/", ".png") {
		r, err := resource.Open("minecraft", file)
		if err != nil {
			panic(err)
		}
		img, err := png.Decode(r)
		if err != nil {
			panic(err)
		}
		width, height := img.Bounds().Dx(), img.Bounds().Dy()
		var ani *animatedTexture
		if width != height {
			height = width
			old := img
			img := image.NewNRGBA(image.Rect(0, 0, width, width))
			draw.Draw(img, img.Bounds(), old, image.ZP, draw.Over)
			ani = loadAnimation(file, old.Bounds().Dy()/old.Bounds().Dx())
			ani.Image = old
			animatedTextures = append(animatedTextures, ani)
		}
		var pix []byte
		switch img := img.(type) {
		case *image.NRGBA:
			pix = img.Pix
		case *image.RGBA:
			pix = img.Pix
		default:
			panic(fmt.Sprintf("unsupported image type %T", img))
		}
		name := file[len("textures/blocks/") : len(file)-4]
		at, rect := addTexture(pix, width, height)
		info := TextureInfo{
			Rect:  rect,
			Atlas: at,
		}
		textureMap[name] = info
		if ani != nil {
			ani.Info = info
		}
	}

	at, rect := addTexture([]byte{
		0, 0, 0, 255,
		255, 0, 255, 255,
		255, 0, 255, 255,
		0, 0, 0, 255,
	}, 2, 2)
	textureMap["missing_texture"] = TextureInfo{
		Rect:  rect,
		Atlas: at,
	}
}

func addTexture(pix []byte, width, height int) (int, *atlas.Rect) {
	for i, a := range textures {
		rect, err := a.Add(pix, width, height)
		if err == nil {
			return i, rect
		}
	}
	a := atlas.New(atlasSize, atlasSize, 4)
	textures = append(textures, a)
	textureViews = append(textureViews, &image.NRGBA{
		Pix:    a.Buffer,
		Stride: 4 * atlasSize,
		Rect:   image.Rect(0, 0, atlasSize, atlasSize),
	})
	rect, err := a.Add(pix, width, height)
	if err != nil {
		panic(err)
	}
	return len(textures) - 1, rect
}

type animatedTexture struct {
	Info          TextureInfo
	Image         image.Image
	Interpolate   bool
	Frames        []textureFrame
	RemainingTime float64
	CurrentFrame  int
}

type textureFrame struct {
	Index int
	Time  int
}

func tickAnimatedTextures(delta float64) {
	delta /= 3 // default is 60 a second, minecraft is 20
	for _, ani := range animatedTextures {
		ani.RemainingTime -= delta
		if ani.RemainingTime < 0 {
			ani.CurrentFrame++
			ani.CurrentFrame %= len(ani.Frames)
			ani.RemainingTime += float64(ani.Frames[ani.CurrentFrame].Time)
			img := textureViews[ani.Info.Atlas]
			r := ani.Info.Rect
			y := r.Width * ani.Frames[ani.CurrentFrame].Index
			src := image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height)
			draw.Draw(img, src, ani.Image, image.Pt(0, y), draw.Src)
		}
	}
	for i, gt := range glTextures {
		gt.Bind(gl.Texture2D)
		gt.Image2D(0, atlasSize, atlasSize, gl.RGBA, gl.UnsignedByte, textures[i].Buffer)
	}
}

func loadAnimation(file string, max int) *animatedTexture {
	a := &animatedTexture{}
	defer func() {
		a.RemainingTime = float64(a.Frames[0].Time)
	}()

	type animation struct {
		FrameTime   int
		Interpolate bool
		Frames      []json.RawMessage
	}
	type base struct {
		Animation animation
	}

	meta, err := resource.Open("minecraft", file+".mcmeta")
	if err != nil {
		panic(err)
	}
	defer meta.Close()
	b := &base{}
	err = json.NewDecoder(meta).Decode(b)
	if err != nil {
		panic(err)
	}
	frameTime := b.Animation.FrameTime
	if frameTime == 0 {
		frameTime = 1
	}
	a.Interpolate = b.Animation.Interpolate

	if len(b.Animation.Frames) == 0 {
		a.Frames = make([]textureFrame, max)
		for i := range a.Frames {
			a.Frames[i] = textureFrame{
				Index: i,
				Time:  frameTime,
			}
		}
		return a
	}

	a.Frames = make([]textureFrame, len(b.Animation.Frames))
	for i := range a.Frames {
		a.Frames[i].Time = frameTime
		if b.Animation.Frames[i][0] == '{' {
			if err = json.Unmarshal(b.Animation.Frames[i], &a.Frames[i]); err != nil {
				panic(err)
			}
			a.Frames[i].Time *= frameTime
			continue
		}
		if err = json.Unmarshal(b.Animation.Frames[i], &a.Frames[i].Index); err != nil {
			panic(err)
		}
	}

	return a
}

type glTexture struct {
	Data          []byte
	Width, Height int
	Format        gl.TextureFormat
	Type          gl.Type
	Filter        gl.TextureValue
	MinFilter     gl.TextureValue
	Wrap          gl.TextureValue
}

func createTexture(t glTexture) gl.Texture {
	if t.Format == 0 {
		t.Format = gl.RGB
	}
	if t.Type == 0 {
		t.Type = gl.UnsignedByte
	}
	if t.Filter == 0 {
		t.Filter = gl.Nearest
	}
	if t.MinFilter == 0 {
		t.MinFilter = t.Filter
	}
	if t.Wrap == 0 {
		t.Wrap = gl.ClampToEdge
	}

	texture := gl.CreateTexture()
	texture.Bind(gl.Texture2D)
	texture.Image2D(0, t.Width, t.Height, t.Format, t.Type, t.Data)
	texture.Parameter(gl.TextureMagFilter, t.Filter)
	texture.Parameter(gl.TextureMinFilter, t.MinFilter)
	texture.Parameter(gl.TextureWrapS, t.Wrap)
	texture.Parameter(gl.TextureWrapT, t.Wrap)
	return texture
}
