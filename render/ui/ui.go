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
	drawables []Drawable
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
	AttachedTo() Drawable
	Attachment() (vAttach, hAttach AttachPoint)
}

// AddDrawable adds the drawable to the draw list and attaches it
// to the specified parts of the screen.
func AddDrawable(d Drawable) {
	drawables = append(drawables, d)
}

var screen = Region{W: scaledWidth, H: scaledHeight}

// Draw draws all drawables in the draw list to the screen.
func Draw(width, height int, delta float64) {
	sw := scaledWidth / float64(width)
	sh := scaledHeight / float64(height)
	if drawMode == mUnscaled {
		sw, sh = 1.0, 1.0
	}
	for _, d := range drawables {
		if !d.ShouldDraw() {
			continue
		}
		r := getDrawRegion(d, sw, sh)
		d.Draw(r, delta)
	}
}

func getDrawRegion(d Drawable, sw, sh float64) Region {
	parent := d.AttachedTo()
	var superR Region
	if parent != nil {
		superR = getDrawRegion(parent, sw, sh)
	} else {
		superR = screen
	}
	r := Region{}
	w, h := d.Size()
	ox, oy := d.Offset()
	r.W = w * sw
	r.H = h * sh
	vAttach, hAttach := d.Attachment()
	switch hAttach {
	case Left:
		r.X = ox * sw
	case Middle:
		r.X = (superR.W / 2) - (r.W / 2) + ox*sw
	case Right:
		r.X = superR.W - ox*sw - r.W
	}
	switch vAttach {
	case Top:
		r.Y = oy * sh
	case Middle:
		r.Y = (superR.H / 2) - (r.H / 2) + oy*sh
	case Right:
		r.Y = superR.H - oy*sh - r.H
	}
	r.X += superR.X
	r.Y += superR.Y
	return r
}

func Remove(d Drawable) {
	for i, dd := range drawables {
		if dd == d {
			drawables = append(drawables[:i], drawables[i+1:]...)
			return
		}
	}
}
