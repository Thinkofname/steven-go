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

// Text is a drawable that draws a string.
type Text struct {
	Parent           Drawable
	X, Y             float64
	R, G, B          int
	A                float64
	value            string
	Width            float64
	Visible          bool
	ScaleX, ScaleY   float64
	Rotation         float64
	vAttach, hAttach AttachPoint
}

// NewText creates a new Text drawable.
func NewText(val string, x, y float64, r, g, b int) *Text {
	return &Text{
		value: val,
		Width: render.SizeOfString(val),
		R:     r, G: g, B: b,
		A: 1.0,
		X: x, Y: y,
		ScaleX: 1, ScaleY: 1,
		Visible: true,
	}
}

// Attach changes the location where this is attached to.
func (t *Text) Attach(vAttach, hAttach AttachPoint) *Text {
	t.vAttach, t.hAttach = vAttach, hAttach
	return t
}

// Attachment returns the sides where this element is attached too.
func (t *Text) Attachment() (vAttach, hAttach AttachPoint) {
	return t.vAttach, t.hAttach
}

// ShouldDraw returns whether this should be drawn at this time.
func (t *Text) ShouldDraw() bool {
	return t.Visible
}

// Update updates the string drawn by this drawable.
func (t *Text) Update(val string) {
	t.value = val
	t.Width = render.SizeOfString(val)
}

// Draw draws this to the target region.
func (t *Text) Draw(r Region, delta float64) {
	cw, ch := t.Size()
	sx, sy := r.W/cw, r.H/ch
	if t.Rotation == 0 {
		render.DrawUITextScaled(t.value, r.X, r.Y, sx*t.ScaleX, sy*t.ScaleY, t.R, t.G, t.B).Alpha(t.A)
	} else {
		render.DrawUITextRotated(t.value, r.X, r.Y, sx*t.ScaleX, sy*t.ScaleY, t.Rotation, t.R, t.G, t.B).Alpha(t.A)
	}
}

// AttachedTo returns the Drawable this is attached to or nil.
func (t *Text) AttachedTo() Drawable {
	return t.Parent
}

// Offset returns the offset of this drawable from the attachment
// point.
func (t *Text) Offset() (float64, float64) {
	return t.X, t.Y
}

// Size returns the size of this drawable.
func (t *Text) Size() (float64, float64) {
	return (t.Width + 2) * t.ScaleX, 18 * t.ScaleY
}

// Remove removes the text element from the draw list.
func (t *Text) Remove() {
	Remove(t)
}
