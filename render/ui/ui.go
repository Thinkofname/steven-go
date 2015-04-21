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

const (
	scaledWidth, scaledHeight = 800, 480
)

var (
	drawMode  = mScaled
	drawables []attachedDrawable
)

// Region is an area for a Drawable to draw to
type Region struct {
	X, Y, W, H float64
}

// Drawable is a scalable element that can be drawn to an
// area.
type Drawable interface {
	Draw(r Region, delta float64)
	Size() (float64, float64)
	// Offset is the offset from the attachment point on
	// each axis
	Offset() (float64, float64)
	ShouldDraw() bool
}

type refStorable interface {
	setRef(r DrawRef)
}

type attachedDrawable struct {
	d                Drawable
	vAttach, hAttach AttachPoint
}

// AddDrawable adds the drawable to the draw list and attaches it
// to the specified parts of the screen. This returns a reference
// that may be used to remove this from the draw list.
func AddDrawable(d Drawable, vAttach, hAttach AttachPoint) DrawRef {
	drawables = append(drawables, attachedDrawable{
		d:       d,
		vAttach: vAttach,
		hAttach: hAttach,
	})
	r := DrawRef{d: d}
	if ra, ok := d.(refStorable); ok {
		ra.setRef(r)
	}
	return r
}

// Draw draws all drawables in the draw list to the screen.
func Draw(width, height int, delta float64) {
	sw := scaledWidth / float64(width)
	sh := scaledHeight / float64(height)
	if drawMode == mUnscaled {
		sw, sh = 1.0, 1.0
	}
	for _, d := range drawables {
		if !d.d.ShouldDraw() {
			continue
		}
		r := Region{}
		w, h := d.d.Size()
		ox, oy := d.d.Offset()
		r.W = w * sw
		r.H = h * sh
		switch d.hAttach {
		case Left:
			r.X = ox * sw
		case Middle:
			r.X = (scaledWidth / 2) - (r.W / 2) + ox*sw
		case Right:
			r.X = scaledWidth - ox*sw - r.W
		}
		switch d.vAttach {
		case Top:
			r.Y = oy * sh
		case Middle:
			r.Y = (scaledHeight / 2) - (r.H / 2) + oy*sh
		case Right:
			r.Y = scaledHeight - oy*sh - r.H
		}
		d.d.Draw(r, delta)
	}
}

// DrawRef is a reference to a Drawable that can be used to
// remove it from the draw list.
type DrawRef struct {
	d Drawable
}

// Remove removes the referenced drawable from the draw list.
func (d DrawRef) Remove() {
	for i, dd := range drawables {
		if dd.d == d.d {
			drawables = append(drawables[:i], drawables[i+1:]...)
			return
		}
	}
}
