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

// Button is a drawable that draws a button.
type Button struct {
	baseElement
	x, y, w, h float64
	disabled   bool
	currentTex render.TextureInfo
	hovered    bool
	hoverFuncs []func(over bool)
	clickFuncs []func()
}

// NewButton creates a new Text drawable.
func NewButton(x, y, w, h float64) *Button {
	return &Button{
		x: x, y: y, w: w, h: h,
		baseElement: baseElement{
			visible: true,
			isNew:   true,
		},
	}
}

// AddClick adds a function to be called when the button is clicked
func (b *Button) AddClick(f func()) {
	b.clickFuncs = append(b.clickFuncs, f)
}

// AddHover adds a function to be called when the button is hovered
func (b *Button) AddHover(f func(over bool)) {
	b.hoverFuncs = append(b.hoverFuncs, f)
}

// Attach changes the location where this is attached to.
func (b *Button) Attach(vAttach, hAttach AttachPoint) *Button {
	b.vAttach, b.hAttach = vAttach, hAttach
	return b
}
func (b *Button) X() float64 { return b.x }
func (b *Button) SetX(x float64) {
	if b.x != x {
		b.x = x
		b.dirty = true
	}
}
func (b *Button) Y() float64 { return b.y }
func (b *Button) SetY(y float64) {
	if b.y != y {
		b.y = y
		b.dirty = true
	}
}
func (b *Button) Width() float64 { return b.w }
func (b *Button) SetWidth(w float64) {
	if b.w != w {
		b.w = w
		b.dirty = true
	}
}
func (b *Button) Height() float64 { return b.h }
func (b *Button) SetHeight(h float64) {
	if b.h != h {
		b.h = h
		b.dirty = true
	}
}
func (b *Button) Disabled() bool { return b.disabled }
func (b *Button) SetDisabled(d bool) {
	if b.disabled != d {
		b.disabled = d
		b.dirty = true
	}
}

func (b *Button) Click(r Region, x, y float64) {
	for _, f := range b.clickFuncs {
		f()
	}
}
func (b *Button) Hover(r Region, x, y float64, over bool) {
	if b.hovered != over {
		b.dirty = true
	}
	b.hovered = over
	for _, f := range b.hoverFuncs {
		f(over)
	}
}

func (b *Button) newUIElement(tex render.TextureInfo, x, y, width, height float64, tx, ty, tw, th float64) *render.UIElement {
	u := render.NewUIElement(tex, x, y, width, height, tx, ty, tw, th)
	u.Layer = b.Layer()
	return u
}

// Draw draws this to the target region.
func (b *Button) Draw(r Region, delta float64) {
	if b.isNew || b.isDirty() || forceDirty {
		b.isNew = false
		if b.disabled {
			b.currentTex = render.RelativeTexture(render.GetTexture("gui/widgets"), 256, 256).
				Sub(0, 46, 200, 20)
		} else {
			off := 66
			if b.hovered {
				off += 20
			}
			b.currentTex = render.RelativeTexture(render.GetTexture("gui/widgets"), 256, 256).
				Sub(0, off, 200, 20)
		}
		b.data = b.data[:0]

		cw, ch := b.Size()
		sx, sy := r.W/cw, r.H/ch
		b.data = append(b.data, b.newUIElement(b.currentTex, r.X, r.Y, 4*sx, 4*sy, 0, 0, 2/200.0, 2/20.0).Bytes()...)
		b.data = append(b.data, b.newUIElement(b.currentTex, r.X+r.W-4*sx, r.Y, 4*sx, 4*sy, 198/200.0, 0, 2/200.0, 2/20.0).Bytes()...)
		b.data = append(b.data, b.newUIElement(b.currentTex, r.X, r.Y+r.H-6*sy, 4*sx, 6*sy, 0, 17/20.0, 2/200.0, 3/20.0).Bytes()...)
		b.data = append(b.data, b.newUIElement(b.currentTex, r.X+r.W-4*sx, r.Y+r.H-6*sy, 4*sx, 6*sy, 198/200.0, 17/20.0, 2/200.0, 3/20.0).Bytes()...)

		w := (r.W/sx)/2 - 4
		b.data = append(b.data, b.newUIElement(b.currentTex.Sub(2, 0, 196, 2), r.X+4*sx, r.Y, r.W-8*sx, 4*sy, 0, 0, w/196.0, 1.0).Bytes()...)
		b.data = append(b.data, b.newUIElement(b.currentTex.Sub(2, 17, 196, 3), r.X+4*sx, r.Y+r.H-6*sy, r.W-8*sx, 6*sy, 0, 0, w/196.0, 1.0).Bytes()...)

		h := (r.H/sy)/2 - 5
		b.data = append(b.data, b.newUIElement(b.currentTex.Sub(0, 2, 2, 15), r.X, r.Y+4*sy, 4*sx, r.H-10*sy, 0.0, 0.0, 1.0, h/16.0).Bytes()...)
		b.data = append(b.data, b.newUIElement(b.currentTex.Sub(198, 2, 2, 15), r.X+r.W-4*sx, r.Y+4*sy, 4*sx, r.H-10*sy, 0.0, 0.0, 1.0, h/16.0).Bytes()...)

		b.data = append(b.data, b.newUIElement(b.currentTex.Sub(2, 2, 196, 15), r.X+4*sx, r.Y+4*sy, r.W-8*sx, r.H-10*sy, 0.0, 0.0, w/196.0, h/16.0).Bytes()...)
	}
	render.UIAddBytes(b.data)
}

// Offset returns the offset of this drawable from the attachment
// point.
func (b *Button) Offset() (float64, float64) {
	return b.x, b.y
}

// Size returns the size of this drawable.
func (b *Button) Size() (float64, float64) {
	return b.w, b.h
}

// Remove removes the Button element from the draw list.
func (b *Button) Remove() {
	Remove(b)
}
