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

import (
	"math"

	"github.com/thinkofdeath/steven/render"
)

// Text is a drawable that draws a string.
type Text struct {
	baseElement
	x, y           float64
	r, g, b, a     int
	value          string
	Width          float64
	scaleX, scaleY float64
	rotation       float64
}

// NewText creates a new Text drawable.
func NewText(val string, x, y float64, r, g, b int) *Text {
	return &Text{
		value: val,
		Width: render.SizeOfString(val),
		r:     r, g: g, b: b, a: 255,
		x: x, y: y,
		scaleX: 1, scaleY: 1,
		baseElement: baseElement{
			visible: true,
			isNew:   true,
		},
	}
}

// Attach changes the location where this is attached to.
func (t *Text) Attach(vAttach, hAttach AttachPoint) *Text {
	t.vAttach, t.hAttach = vAttach, hAttach
	return t
}

func (t *Text) Value() string { return t.value }
func (t *Text) X() float64    { return t.x }
func (t *Text) SetX(x float64) {
	if t.x != x {
		t.x = x
		t.dirty = true
	}
}
func (t *Text) Y() float64 { return t.y }
func (t *Text) SetY(y float64) {
	if t.y != y {
		t.y = y
		t.dirty = true
	}
}
func (t *Text) R() int { return t.r }
func (t *Text) SetR(r int) {
	if t.r != r {
		t.r = r
		t.dirty = true
	}
}
func (t *Text) G() int { return t.g }
func (t *Text) SetG(g int) {
	if t.g != g {
		t.g = g
		t.dirty = true
	}
}
func (t *Text) B() int { return t.b }
func (t *Text) SetB(b int) {
	if t.b != b {
		t.b = b
		t.dirty = true
	}
}
func (t *Text) A() int { return t.a }
func (t *Text) SetA(a int) {
	if a > 255 {
		a = 255
	}
	if a < 0 {
		a = 0
	}
	if t.a != a {
		t.a = a
		t.dirty = true
	}
}
func (t *Text) ScaleX() float64 { return t.scaleX }
func (t *Text) SetScaleX(s float64) {
	if t.scaleX != s {
		t.scaleX = s
		t.dirty = true
	}
}
func (t *Text) ScaleY() float64 { return t.scaleY }
func (t *Text) SetScaleY(s float64) {
	if t.scaleY != s {
		t.scaleY = s
		t.dirty = true
	}
}
func (t *Text) Rotation() float64 { return t.rotation }
func (t *Text) SetRotation(r float64) {
	if t.rotation != r {
		t.rotation = r
		t.dirty = true
	}
}

// Update updates the string drawn by this drawable.
func (t *Text) Update(val string) {
	if t.value == val {
		return
	}
	t.value = val
	t.Width = render.SizeOfString(val)
	t.dirty = true
}

// Draw draws this to the target region.
func (t *Text) Draw(r Region, delta float64) {
	if t.isNew || t.isDirty() || forceDirty {
		t.isNew = false
		cw, ch := t.Size()
		sx, sy := r.W/cw, r.H/ch
		var text render.UIText
		if t.rotation == 0 {
			text = render.NewUITextScaled(t.value, r.X, r.Y, sx*t.scaleX, sy*t.scaleY, t.r, t.g, t.b)
		} else {
			c := math.Cos(t.rotation)
			s := math.Sin(t.rotation)
			tmpx := r.W / 2
			tmpy := r.H / 2
			w := math.Abs(tmpx*c - tmpy*s)
			h := math.Abs(tmpy*c + tmpx*s)
			text = render.NewUITextRotated(t.value, r.X+w-(r.W/2), r.Y+h-(r.H/2), sx*t.scaleX, sy*t.scaleY, t.rotation, t.r, t.g, t.b)
		}
		text.Alpha(t.a)
		for _, txt := range text.Elements {
			txt.Layer = t.Layer()
		}
		t.data = text.Bytes()
	}
	render.UIAddBytes(t.data)
}

// Offset returns the offset of this drawable from the attachment
// point.
func (t *Text) Offset() (float64, float64) {
	return t.x, t.y
}

// Size returns the size of this drawable.
func (t *Text) Size() (float64, float64) {
	w, h := (t.Width + 2), 18.0
	return w * t.scaleX, h * t.scaleY
}

// Remove removes the text element from the draw list.
func (t *Text) Remove() {
	Remove(t)
}
