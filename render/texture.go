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
	"strings"
	"sync"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/render/atlas"
	"github.com/thinkofdeath/steven/render/gl"
	"github.com/thinkofdeath/steven/resource"
)

var (
	textures         []*atlas.Type
	textureCount     int
	textureMap       = map[string]*textureInfo{}
	loadedTextures   []*loadedTexture
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
		ox:     float64(x) / float64(ti.rect.Width),
		oy:     float64(y) / float64(ti.rect.Height),
		w:      float64(w) / float64(ti.rect.Width),
		h:      float64(h) / float64(ti.rect.Height),
		parent: ti,
	}
}

type subTextureInfo struct {
	ox, oy, w, h float64
	parent       TextureInfo
}

func (s *subTextureInfo) Atlas() int { return s.parent.Atlas() }
func (s *subTextureInfo) Rect() atlas.Rect {
	rect := s.parent.Rect()
	rect.X += int(s.ox * float64(rect.Width))
	rect.Y += int(s.oy * float64(rect.Height))
	rect.Width = int(s.w * float64(rect.Width))
	rect.Height = int(s.h * float64(rect.Height))
	return rect
}

// Sub returns a subsection of this texture.
func (s *subTextureInfo) Sub(x, y, w, h int) TextureInfo {
	rect := s.Rect()
	return &subTextureInfo{
		ox:     float64(x) / float64(rect.Width),
		oy:     float64(y) / float64(rect.Height),
		w:      float64(w) / float64(rect.Width),
		h:      float64(h) / float64(rect.Height),
		parent: s,
	}
}

func RelativeTexture(ti TextureInfo, w, h int) TextureInfo {
	return &relativeTexture{
		w: 1, h: 1,
		fakeW: float64(w), fakeH: float64(h),
		parent: ti,
	}
}

type relativeTexture struct {
	ox, oy, w, h float64
	fakeW, fakeH float64
	parent       TextureInfo
}

func (r *relativeTexture) Atlas() int { return r.parent.Atlas() }
func (r *relativeTexture) Rect() atlas.Rect {
	rect := r.parent.Rect()
	rect.X += int(r.ox * float64(rect.Width))
	rect.Y += int(r.oy * float64(rect.Height))
	rect.Width = int(r.w * float64(rect.Width))
	rect.Height = int(r.h * float64(rect.Height))
	return rect
}

func (r *relativeTexture) Sub(x, y, w, h int) TextureInfo {
	return &relativeTexture{
		ox:     float64(x) / float64(r.fakeW),
		oy:     float64(y) / float64(r.fakeH),
		w:      float64(w) / float64(r.fakeW),
		h:      float64(h) / float64(r.fakeH),
		fakeW:  float64(w),
		fakeH:  float64(h),
		parent: r,
	}
}

// GetTexture returns the related TextureInfo for the requested texture.
// If the texture isn't found a placeholder is returned instead.
// The plugin prefix of 'minecraft:' is defualt
func GetTexture(name string) TextureInfo {
	textureLock.RLock()
	defer textureLock.RUnlock()
	t, ok := textureMap[name]
	if !ok {
		textureLock.RUnlock()
		ret := make(chan struct{}, 1)
		f := func() {
			textureLock.Lock()
			defer textureLock.Unlock()
			// Check to see if it was already loaded between
			// requesting it
			if _, ok = textureMap[name]; ok {
				ret <- struct{}{}
				return
			}
			ns := name
			plugin := "minecraft"
			if pos := strings.IndexRune(name, ':'); pos != -1 {
				plugin = name[:pos]
				ns = name[pos+1:]
			}
			r, err := resource.Open(plugin, "textures/"+ns+".png")
			if err == nil {
				defer r.Close()
				img, err := png.Decode(r)
				if err != nil {
					panic(fmt.Sprintf("(%s): %s", name, err))
				}
				s := &loadedTexture{
					Plugin: plugin,
					File:   "textures/" + ns + ".png",
					Image:  img,
				}
				loadedTextures = append(loadedTextures, s)
				loadTexFile(s)
			}
			t, ok = textureMap[name]
			if !ok {
				t = textureMap["missing_texture"]
				textureMap[name] = t
			}
			ret <- struct{}{}
		}
		if w := glfw.GetCurrentContext(); w == nil {
			syncChan <- f
		} else {
			f()
		}
		<-ret
		textureLock.RLock()
	}
	return t
}

type loadedTexture struct {
	Plugin string
	File   string
	Image  image.Image
}

// LoadTextures (re)loads all the block textures from the resource pack(s)
// TODO(Think) better error handling (if possible to recover?)
func LoadTextures() {
	textureLock.Lock()

	if texturesCreated {
		glTexture.Bind(gl.Texture2DArray)
		data := make([]byte, AtlasSize*AtlasSize*textureCount*4)
		glTexture.Image3D(0, AtlasSize, AtlasSize, textureCount, gl.RGBA, gl.UnsignedByte, data)
	} else {
		glTexture = gl.CreateTexture()
		glTexture.Bind(gl.Texture2DArray)
		textureDepth = len(textures)
		glTexture.Image3D(0, AtlasSize, AtlasSize, len(textures), gl.RGBA, gl.UnsignedByte, make([]byte, AtlasSize*AtlasSize*len(textures)*4))
		glTexture.Parameter(gl.TextureMagFilter, gl.Nearest)
		glTexture.Parameter(gl.TextureMinFilter, gl.Nearest)
		glTexture.Parameter(gl.TextureWrapS, gl.ClampToEdge)
		glTexture.Parameter(gl.TextureWrapT, gl.ClampToEdge)
		for i, tex := range textures {
			glTexture.SubImage3D(0, 0, 0, i, AtlasSize, AtlasSize, 1, gl.RGBA, gl.UnsignedByte, tex.Buffer)
			textures[i] = nil
		}
		texturesCreated = true
	}
	freeTextures = nil
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

	for _, s := range loadedTextures {
		func() {
			r, err := resource.Open(s.Plugin, s.File)
			if err == nil {
				defer r.Close()
				img, err := png.Decode(r)
				if err != nil {
					panic(fmt.Sprintf("(%s): %s", s.File, err))
				}
				s.Image = img
				loadTexFile(s)
			}
		}()
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
		info := getFreeTexture(skin.info.W, skin.info.H)
		uploadTexture(info.info, skin.data)
		i := skin.info.info
		skin.info = info
		i.rect = info.info.rect
		i.atlas = info.info.atlas
		skin.info.info = i
	}

	textureLock.Unlock()

	loadFontInfo()
	loadFontPage(0)

	for i := range isFontLoaded {
		if i == 0 || !isFontLoaded[i] {
			continue
		}
		isFontLoaded[i] = false
		loadFontPage(i)
	}
}

func loadTexFile(st *loadedTexture) {
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
	if st.Plugin != "minecraft" {
		name = st.Plugin + ":" + name
	}
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
