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
	Parent  Drawable
	X, Y    float64
	R, G, B int
	value   string
	Width   float64
	Visible bool
	ref     DrawRef
}

// NewText creates a new Text drawable.
func NewText(val string, x, y float64, r, g, b int) *Text {
	return &Text{
		value: val,
		Width: render.SizeOfString(val),
		R:     r, G: g, B: b,
		X: x, Y: y,
		Visible: true,
	}
}
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
	render.DrawUITextScaled(t.value, r.X, r.Y, sx, sy, t.R, t.G, t.B)
}

// Offset returns the offset of this drawable from the attachment
// point. Includes the parent if it has one.
func (t *Text) Offset() (float64, float64) {
	if t.Parent != nil {
		ox, oy := t.Parent.Offset()
		return ox + t.X, oy + t.Y
	}
	return t.X, t.Y
}

// Size returns the size of this drawable.
func (t *Text) Size() (float64, float64) {
	return t.Width + 2, 18
}

func (t *Text) setRef(r DrawRef) {
	t.ref = r
}

// Remove removes the text element from the draw list.
func (t *Text) Remove() {
	t.ref.Remove()
}
