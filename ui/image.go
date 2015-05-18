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

package ui

import "github.com/thinkofdeath/steven/render"

// Image is a drawable that draws a texture.
type Image struct {
	baseElement
	texture        *render.TextureInfo
	x, y, w, h     float64
	tx, ty, tw, th float64
	r, g, b, a     int
}

// NewImage creates a new image drawable.
func NewImage(texture *render.TextureInfo, x, y, w, h, tx, ty, tw, th float64, r, g, b int) *Image {
	return &Image{
		texture: texture,
		r:       r, g: g, b: b, a: 255,
		x: x, y: y, w: w, h: h,
		tx: tx, ty: ty, tw: tw, th: th,
		baseElement: baseElement{
			visible: true,
			isNew:   true,
		},
	}
}

// Attach changes the location where this is attached to.
func (i *Image) Attach(vAttach, hAttach AttachPoint) *Image {
	i.vAttach, i.hAttach = vAttach, hAttach
	return i
}

func (i *Image) Texture() *render.TextureInfo { return i.texture }
func (i *Image) SetTexture(t *render.TextureInfo) {
	if i.texture != t {
		i.texture = t
		i.dirty = true
	}
}
func (i *Image) X() float64 { return i.x }
func (i *Image) SetX(x float64) {
	if i.x != x {
		i.x = x
		i.dirty = true
	}
}
func (i *Image) Y() float64 { return i.y }
func (i *Image) SetY(y float64) {
	if i.y != y {
		i.y = y
		i.dirty = true
	}
}
func (i *Image) Width() float64 { return i.w }
func (i *Image) SetWidth(w float64) {
	if i.w != w {
		i.w = w
		i.dirty = true
	}
}
func (i *Image) Height() float64 { return i.h }
func (i *Image) SetHeight(h float64) {
	if i.h != h {
		i.h = h
		i.dirty = true
	}
}
func (i *Image) TextureX() float64 { return i.tx }
func (i *Image) SetTextureX(x float64) {
	if i.tx != x {
		i.tx = x
		i.dirty = true
	}
}
func (i *Image) TextureY() float64 { return i.ty }
func (i *Image) SetTextureY(y float64) {
	if i.ty != y {
		i.ty = y
		i.dirty = true
	}
}
func (i *Image) TextureWidth() float64 { return i.tw }
func (i *Image) SetTextureWidth(w float64) {
	if i.tw != w {
		i.tw = w
		i.dirty = true
	}
}
func (i *Image) TextureHeight() float64 { return i.th }
func (i *Image) SetTextureHeight(h float64) {
	if i.th != h {
		i.th = h
		i.dirty = true
	}
}
func (i *Image) R() int { return i.r }
func (i *Image) SetR(r int) {
	if i.r != r {
		i.r = r
		i.dirty = true
	}
}
func (i *Image) G() int { return i.g }
func (i *Image) SetG(g int) {
	if i.g != g {
		i.g = g
		i.dirty = true
	}
}
func (i *Image) B() int { return i.b }
func (i *Image) SetB(b int) {
	if i.b != b {
		i.b = b
		i.dirty = true
	}
}
func (i *Image) A() int { return i.a }
func (i *Image) SetA(a int) {
	if i.a != a {
		i.a = a
		i.dirty = true
	}
}

// Draw draws this to the target region.
func (i *Image) Draw(r Region, delta float64) {
	if i.isNew || i.isDirty() || forceDirty {
		i.isNew = false
		e := render.NewUIElement(i.texture, r.X, r.Y, r.W, r.H, i.tx, i.ty, i.tw, i.th)
		e.R = byte(i.r)
		e.G = byte(i.g)
		e.B = byte(i.b)
		e.A = byte(i.a)
		i.data = e.Bytes()
	}
	render.UIAddBytes(i.data)
}

// Offset returns the offset of this drawable from the attachment
// point.
func (i *Image) Offset() (float64, float64) {
	return i.x, i.y
}

// Size returns the size of this drawable.
func (i *Image) Size() (float64, float64) {
	return i.w, i.h
}

// Remove removes the image element from the draw list.
func (i *Image) Remove() {
	Remove(i)
}
