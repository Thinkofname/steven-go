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

// Package ui provides methods to draw a user interface onto the
// the screen and manage resizing.
package ui

const (
	scaledWidth, scaledHeight = 854, 480
)

var (
	// DrawMode is the scaling mode used.
	DrawMode = Scaled
	// Scale controls the scaling manually when DrawModel is Unscaled
	Scale = 1.0

	drawables []drawRef
)

func ForceDraw() {
	lastWidth = -1
}

type drawRef struct {
	Drawable
	removeHook func(d Drawable)
}

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

	isDirty() bool
	flagDirty()
	clearDirty()
}

type Interactable interface {
	Click(r Region, x, y float64)
	Hover(r Region, x, y float64, over bool)
}

// AddDrawable adds the drawable to the draw list.
func AddDrawable(d Drawable) {
	d.flagDirty()
	drawables = append(drawables, drawRef{Drawable: d})
}

// AddDrawableHook adds the drawable to the draw list.
// The passed function will be called when the drawable
// is removed.
func AddDrawableHook(d Drawable, hook func(d Drawable)) {
	d.flagDirty()
	drawables = append(drawables, drawRef{Drawable: d, removeHook: hook})
}

var screen = Region{W: scaledWidth, H: scaledHeight}

var (
	lastWidth, lastHeight int
	forceDirty            bool
)

// Draw draws all drawables in the draw list to the screen.
func Draw(width, height int, delta float64) {
	sw := scaledWidth / float64(width)
	sh := scaledHeight / float64(height)
	if DrawMode == Unscaled {
		sw, sh = Scale, Scale
	}

	for _, d := range drawables {
		if !d.ShouldDraw() {
			continue
		}
		r := getDrawRegion(d, sw, sh)
		if r.intersects(screen) {
			d.Draw(r, delta)
		}
	}

	for _, d := range drawables {
		// Handle parents that aren't drawing too
		for r := d.Drawable; r != nil; r = r.AttachedTo() {
			r.clearDirty()
		}
	}
	forceDirty = false
	if lastWidth != width || lastHeight != height {
		forceDirty = true
		lastWidth, lastHeight = width, height
	}
}

func (r Region) intersects(o Region) bool {
	return !(r.X+r.W < o.X ||
		r.X > o.X+o.W ||
		r.Y+r.H < o.Y ||
		r.Y > o.Y+o.H)
}

// Hover calls Hover on all interactables at the passed location.
func Hover(x, y float64, width, height int) {
	sw := scaledWidth / float64(width)
	sh := scaledHeight / float64(height)
	if DrawMode == Unscaled {
		sw, sh = Scale, Scale
	}
	x = (x / float64(width)) * scaledWidth
	y = (y / float64(height)) * scaledHeight
	for i := range drawables {
		d := drawables[len(drawables)-1-i]
		inter, ok := d.Drawable.(Interactable)
		if !ok {
			continue
		}
		r := getDrawRegion(d, sw, sh)
		if x >= r.X && x <= r.X+r.W && y >= r.Y && y <= r.Y+r.H {
			inter.Hover(r, x, y, true)
		} else {
			inter.Hover(r, x, y, false)
		}
	}
}

// Click calls Click on all interactables at the passed location.
func Click(x, y float64, width, height int) {
	sw := scaledWidth / float64(width)
	sh := scaledHeight / float64(height)
	if DrawMode == Unscaled {
		sw, sh = Scale, Scale
	}
	x = (x / float64(width)) * scaledWidth
	y = (y / float64(height)) * scaledHeight
	for i := range drawables {
		d := drawables[len(drawables)-1-i]
		inter, ok := d.Drawable.(Interactable)
		if !ok {
			continue
		}
		r := getDrawRegion(d, sw, sh)
		if x >= r.X && x <= r.X+r.W && y >= r.Y && y <= r.Y+r.H {
			inter.Click(r, x, y)
			break
		}
	}
}

// Intersects returns whether the point x,y intersects with the drawable
func Intersects(d Drawable, x, y float64, width, height int) (float64, float64, bool) {
	sw := scaledWidth / float64(width)
	sh := scaledHeight / float64(height)
	if DrawMode == Unscaled {
		sw, sh = Scale, Scale
	}
	x = (x / float64(width)) * scaledWidth
	y = (y / float64(height)) * scaledHeight
	r := getDrawRegion(d, sw, sh)
	if x >= r.X && x <= r.X+r.W && y >= r.Y && y <= r.Y+r.H {
		w, h := d.Size()
		ox := ((x - r.X) / r.W) * w
		oy := ((y - r.Y) / r.H) * h
		return ox, oy, true
	}
	return 0, 0, false
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

// Remove removes the drawable from the screen.
func Remove(d Drawable) {
	for i, dd := range drawables {
		if dd.Drawable == d {
			if dd.removeHook != nil {
				dd.removeHook(d)
			}
			drawables = append(drawables[:i], drawables[i+1:]...)
			return
		}
	}
}

type baseElement struct {
	parent           Drawable
	visible          bool
	vAttach, hAttach AttachPoint

	dirty bool
	isNew bool
	data  []byte
}

// Attachment returns the sides where this element is attached too.
func (b *baseElement) Attachment() (vAttach, hAttach AttachPoint) {
	return b.vAttach, b.hAttach
}

// ShouldDraw returns whether this should be drawn at this time.
func (b *baseElement) ShouldDraw() bool {
	return b.visible
}

func (b *baseElement) SetDraw(shouldDraw bool) {
	if shouldDraw != b.visible {
		b.visible = shouldDraw
		b.dirty = true
	}
}

// AttachedTo returns the Drawable this is attached to or nil.
func (b *baseElement) AttachedTo() Drawable {
	return b.parent
}

func (b *baseElement) AttachTo(d Drawable) {
	if b.parent != d {
		b.parent = d
		b.dirty = true
	}
}

func (b *baseElement) isDirty() bool {
	return b.dirty || (b.parent != nil && b.parent.isDirty())
}

func (b *baseElement) flagDirty() {
	b.dirty = true
}

func (b *baseElement) clearDirty() {
	b.dirty = false
}
