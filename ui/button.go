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

import "github.com/thinkofdeath/phteven/render"

// Button is a drawable that draws a button.
type Button struct {
	Parent           Drawable
	X, Y, W, H       float64
	Visible          bool
	vAttach, hAttach AttachPoint
	currentTex       *render.TextureInfo
	Disabled         bool
	hovered          bool
	HoverFunc        func(over bool)
	ClickFunc        func()
}

// NewButton creates a new Text drawable.
func NewButton(x, y, w, h float64) *Button {
	return &Button{
		X: x, Y: y, W: w, H: h,
		Visible: true,
	}
}

// Attach changes the location where this is attached to.
func (b *Button) Attach(vAttach, hAttach AttachPoint) *Button {
	b.vAttach, b.hAttach = vAttach, hAttach
	return b
}

// Attachment returns the sides where this element is attached too.
func (b *Button) Attachment() (vAttach, hAttach AttachPoint) {
	return b.vAttach, b.hAttach
}

// ShouldDraw returns whether this should be drawn at this time.
func (b *Button) ShouldDraw() bool {
	return b.Visible
}

func (b *Button) Click(r Region, x, y float64) {
	if b.ClickFunc != nil {
		b.ClickFunc()
	}
}
func (b *Button) Hover(r Region, x, y float64, over bool) {
	b.hovered = over
	if b.HoverFunc != nil {
		b.HoverFunc(over)
	}
}

// Draw draws this to the target region.
func (b *Button) Draw(r Region, delta float64) {
	if b.Disabled {
		b.currentTex = render.GetTexture("gui/widgets").Sub(0, 46, 200, 20)
	} else {
		off := 66
		if b.hovered {
			off += 20
		}
		b.currentTex = render.GetTexture("gui/widgets").Sub(0, off, 200, 20)
	}
	cw, ch := b.Size()
	sx, sy := r.W/cw, r.H/ch
	render.DrawUIElement(b.currentTex, r.X, r.Y, 4*sx, 4*sy, 0, 0, 2/200.0, 2/20.0)
	render.DrawUIElement(b.currentTex, r.X+r.W-4*sx, r.Y, 4*sx, 4*sy, 198/200.0, 0, 2/200.0, 2/20.0)
	render.DrawUIElement(b.currentTex, r.X, r.Y+r.H-6*sy, 4*sx, 6*sy, 0, 17/20.0, 2/200.0, 3/20.0)
	render.DrawUIElement(b.currentTex, r.X+r.W-4*sx, r.Y+r.H-6*sy, 4*sx, 6*sy, 198/200.0, 17/20.0, 2/200.0, 3/20.0)

	w := (r.W/sx)/2 - 4
	render.DrawUIElement(b.currentTex.Sub(2, 0, 196, 2), r.X+4*sx, r.Y, r.W-8*sx, 4*sy, 0, 0, w/196.0, 1.0)
	render.DrawUIElement(b.currentTex.Sub(2, 17, 196, 3), r.X+4*sx, r.Y+r.H-6*sy, r.W-8*sx, 6*sy, 0, 0, w/196.0, 1.0)

	h := (r.H/sy)/2 - 5
	render.DrawUIElement(b.currentTex.Sub(0, 2, 2, 15), r.X, r.Y+4*sy, 4*sx, r.H-10*sy, 0.0, 0.0, 1.0, h/16.0)
	render.DrawUIElement(b.currentTex.Sub(198, 2, 2, 15), r.X+r.W-4*sx, r.Y+4*sy, 4*sx, r.H-10*sy, 0.0, 0.0, 1.0, h/16.0)

	render.DrawUIElement(b.currentTex.Sub(2, 2, 196, 15), r.X+4*sx, r.Y+4*sy, r.W-8*sx, r.H-10*sy, 0.0, 0.0, w/196.0, h/16.0)
}

// AttachedTo returns the Drawable this is attached to or nil.
func (b *Button) AttachedTo() Drawable {
	return b.Parent
}

// Offset returns the offset of this drawable from the attachment
// point.
func (b *Button) Offset() (float64, float64) {
	return b.X, b.Y
}

// Size returns the size of this drawable.
func (i *Button) Size() (float64, float64) {
	return i.W, i.H
}

// Remove removes the Button element from the draw list.
func (b *Button) Remove() {
	Remove(b)
}
