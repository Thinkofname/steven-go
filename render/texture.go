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
	"sort"
	"strings"
	"sync"

	"github.com/thinkofdeath/steven/render/atlas"
	"github.com/thinkofdeath/steven/render/gl"
	"github.com/thinkofdeath/steven/resource"
)

var (
	textures         []*atlas.Type
	textureCount     int
	textureMap       = map[string]*textureInfo{}
	textureLock      sync.RWMutex
	animatedTextures []*animatedTexture
)

const (
	AtlasSize    = 1024
	atlasSizeStr = "1024"
)

// TextureInfo returns information about a texture in an atlas
type TextureInfo interface {
	Atlas() int
	Rect() atlas.Rect
	Sub(x, y, w, h int) TextureInfo
}

type textureInfo struct {
	atlas int
	rect  atlas.Rect
}

func (ti *textureInfo) Atlas() int       { return ti.atlas }
func (ti *textureInfo) Rect() atlas.Rect { return ti.rect }

// Sub returns a subsection of this texture.
func (ti *textureInfo) Sub(x, y, w, h int) TextureInfo {
	return &subTextureInfo{
		ox: x, oy: y, w: w, h: h,
		parent: ti,
	}
}

type subTextureInfo struct {
	ox, oy, w, h int
	parent       TextureInfo
}

func (s *subTextureInfo) Atlas() int { return s.parent.Atlas() }
func (s *subTextureInfo) Rect() atlas.Rect {
	rect := s.parent.Rect()
	rect.X += s.ox
	rect.Y += s.oy
	rect.Width, rect.Height = s.w, s.h
	return rect
}

// Sub returns a subsection of this texture.
func (s *subTextureInfo) Sub(x, y, w, h int) TextureInfo {
	return &subTextureInfo{
		ox: x, oy: y, w: w, h: h,
		parent: s,
	}
}

// GetTexture returns the related TextureInfo for the requested texture.
// If the texture isn't found a placeholder is returned instead.
func GetTexture(name string) TextureInfo {
	textureLock.RLock()
	defer textureLock.RUnlock()
	t, ok := textureMap[name]
	if !ok {
		t = textureMap["missing_texture"]
	}
	return t
}

type sortableTexture struct {
	Area  int
	File  string
	Image image.Image
}

type sortableTextures []sortableTexture

func (s sortableTextures) Swap(a, b int) {
	s[a], s[b] = s[b], s[a]
}

func (s sortableTextures) Len() int {
	return len(s)
}

func (s sortableTextures) Less(a, b int) bool {
	return s[a].Area > s[b].Area
}

// LoadTextures (re)loads all the block textures from the resource pack(s)
// TODO(Think) better error handling (if possible to recover?)
func LoadTextures() {
	textureLock.Lock()

	if texturesCreated {
		glTexture.Bind(gl.Texture2DArray)
		data := make([]byte, AtlasSize*AtlasSize*textureCount*4)
		glTexture.Image3D(0, AtlasSize, AtlasSize, textureCount, gl.RGBA, gl.UnsignedByte, data)
	}
	freeSkinTextures = nil
	for i := range isFontLoaded {
		isFontLoaded[i] = false
		fontPages[i] = nil
	}
	animatedTextures = nil
	textures = nil
	pix := []byte{
		0, 0, 0, 255,
		255, 0, 255, 255,
		255, 0, 255, 255,
		0, 0, 0, 255,
	}
	info := addTexture(pix, 2, 2)
	if t, ok := textureMap["missing_texture"]; ok {
		t.atlas = info.atlas
		t.rect = info.rect
	} else {
		textureMap["missing_texture"] = info
	}

	for _, t := range textureMap {
		t.atlas = info.atlas
		t.rect = info.rect
	}

	names := resource.Search("minecraft", "textures/", ".png")
	tList := make(sortableTextures, 0, len(names))
	for _, file := range names {
		if strings.HasPrefix(file, "textures/font") {
			continue
		}
		r, err := resource.Open("minecraft", file)
		if err != nil {
			panic(err)
		}
		defer r.Close()
		img, err := png.Decode(r)
		if err != nil {
			panic(err)
		}
		width, height := img.Bounds().Dx(), img.Bounds().Dy()
		tList = append(tList, sortableTexture{
			Area:  width * height,
			File:  file,
			Image: img,
		})
	}
	sort.Sort(tList)
	for _, st := range tList {
		loadTexFile(st)
	}

	pix = []byte{255, 255, 255, 255}
	info = addTexture(pix, 1, 1)
	if t, ok := textureMap["solid"]; ok {
		t.atlas = info.atlas
		t.rect = info.rect
	} else {
		textureMap["solid"] = info
	}

	for _, skin := range skins {
		info := getSkinInfo()
		uploadTexture(info, skin.data)
		skin.info.rect = info.rect
		skin.info.atlas = info.atlas
	}

	textureLock.Unlock()

	loadFontInfo()
	loadFontPage(0)
}

func loadTexFile(st sortableTexture) {
	file := st.File
	ii := st.Image
	img := ii.(draw.Image)
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	var ani *animatedTexture
	if (strings.HasPrefix(file, "textures/blocks") || strings.HasPrefix(file, "textures/items")) &&
		width != height {
		height = width
		old := img
		img = image.NewNRGBA(image.Rect(0, 0, width, width))
		draw.Draw(img, img.Bounds(), old, image.ZP, draw.Over)
		ani = loadAnimation(file, old.Bounds().Dy()/old.Bounds().Dx())
		if ani != nil {
			ani.Image = old
			ani.Buffer = imgToBytes(old)
			animatedTextures = append(animatedTextures, ani)
		} else {
			img = old
			width, height = img.Bounds().Dx(), img.Bounds().Dy()
		}
	}
	pix := imgToBytes(img)
	name := file[len("textures/") : len(file)-4]
	info := addTexture(pix, width, height)
	if t, ok := textureMap[name]; ok {
		t.atlas = info.atlas
		t.rect = info.rect
	} else {
		textureMap[name] = info
	}
	if ani != nil {
		ani.Info = info
	}
}

func imgToBytes(img image.Image) []byte {
	switch img := img.(type) {
	case *image.NRGBA:
		return img.Pix
	case *image.RGBA:
		return img.Pix
	default:
		temp := image.NewNRGBA(img.Bounds())
		draw.Draw(temp, img.Bounds(), img, image.ZP, draw.Over)
		return temp.Pix
	}
}

func addTexture(pix []byte, width, height int) *textureInfo {
	for i, a := range textures {
		if a == nil {
			continue
		}
		rect, err := a.Add(pix, width, height)
		if err == nil {
			info := &textureInfo{atlas: i, rect: *rect}
			if texturesCreated {
				uploadTexture(info, pix)
			}
			return info
		}
	}
	var a *atlas.Type
	if texturesCreated {
		a = atlas.NewLight(AtlasSize, AtlasSize, 0)
	} else {
		a = atlas.New(AtlasSize, AtlasSize, 4)
	}
	textures = append(textures, a)
	reupload := false
	if len(textures) > textureCount {
		textureCount = len(textures)
		reupload = true
	}
	rect, err := a.Add(pix, width, height)
	if err != nil {
		panic(fmt.Sprintf("Failed to place %d,%d: %s", width, height, err))
	}

	info := &textureInfo{atlas: len(textures) - 1, rect: *rect}
	if texturesCreated {
		if reupload {
			glTexture.Bind(gl.Texture2DArray)
			data := make([]byte, AtlasSize*AtlasSize*textureCount*4)
			glTexture.Get(0, gl.RGBA, gl.UnsignedByte, data)
			glTexture.Image3D(0, AtlasSize, AtlasSize, textureCount, gl.RGBA, gl.UnsignedByte, data)
			textureDepth = textureCount
		}
		uploadTexture(info, pix)
	}

	return info
}

func uploadTexture(info *textureInfo, data []byte) {
	glTexture.Bind(gl.Texture2DArray)
	r := info.rect
	glTexture.SubImage3D(0, r.X, r.Y, info.atlas, r.Width, r.Height, 1, gl.RGBA, gl.UnsignedByte, data)
}

type animatedTexture struct {
	Info          *textureInfo
	Image         image.Image
	Buffer        []byte
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
			r := ani.Info.rect
			glTexture.Bind(gl.Texture2DArray)
			offset := r.Width * r.Width * ani.Frames[ani.CurrentFrame].Index * 4
			offset2 := offset + r.Height*r.Width*4
			glTexture.SubImage3D(0, r.X, r.Y, ani.Info.atlas, r.Width, r.Height, 1, gl.RGBA, gl.UnsignedByte, ani.Buffer[offset:offset2])
		}
	}
}

func loadAnimation(file string, max int) *animatedTexture {
	a := &animatedTexture{}
	defer func() {
		if a != nil {
			a.RemainingTime = float64(a.Frames[0].Time)
		}
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
		fmt.Printf("%s: %s\n", file+".mcmeta", err)
		a = nil
		return nil
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
