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
	"github.com/thinkofdeath/steven/chat"
	"github.com/thinkofdeath/steven/render"
)

// Formatted is a drawable that draws a string.
type Formatted struct {
	Parent           Drawable
	X, Y             float64
	Width, Height    float64
	MaxWidth         float64
	Visible          bool
	ScaleX, ScaleY   float64
	vAttach, hAttach AttachPoint

	text []*Text
}

// NewFormatted creates a new Formatted drawable.
func NewFormatted(val chat.AnyComponent, x, y float64) *Formatted {
	f := &Formatted{
		X: x, Y: y,
		ScaleX: 1, ScaleY: 1,
		Visible:  true,
		MaxWidth: -1,
	}
	f.Update(val)
	return f
}

// NewFormattedWidth creates a new Formatted drawable with a max width.
func NewFormattedWidth(val chat.AnyComponent, x, y, width float64) *Formatted {
	f := &Formatted{
		X: x, Y: y,
		ScaleX: 1, ScaleY: 1,
		Visible:  true,
		MaxWidth: width,
	}
	f.Update(val)
	return f
}

// Attach changes the location where this is attached to.
func (f *Formatted) Attach(vAttach, hAttach AttachPoint) *Formatted {
	f.vAttach, f.hAttach = vAttach, hAttach
	return f
}

// Attachment returns the sides where this element is attached too.
func (f *Formatted) Attachment() (vAttach, hAttach AttachPoint) {
	return f.vAttach, f.hAttach
}

// ShouldDraw returns whether this should be drawn at this time.
func (f *Formatted) ShouldDraw() bool {
	return f.Visible
}

// Draw draws this to the target region.
func (f *Formatted) Draw(r Region, delta float64) {
	cw, ch := f.Size()
	sx, sy := r.W/cw, r.H/ch
	for _, t := range f.text {
		r := getDrawRegion(t, sx, sy)
		t.Draw(r, delta)
	}
}

// AttachedTo returns the Drawable this is attached to or nil.
func (f *Formatted) AttachedTo() Drawable {
	return f.Parent
}

// Offset returns the offset of this drawable from the attachment
// point.
func (f *Formatted) Offset() (float64, float64) {
	return f.X, f.Y
}

// Size returns the size of this drawable.
func (f *Formatted) Size() (float64, float64) {
	return (f.Width + 2) * f.ScaleX, f.Height * f.ScaleY
}

// Remove removes the Formatted element from the draw list.
func (f *Formatted) Remove() {
	Remove(f)
}

// Update updates the component drawn by this drawable.
func (f *Formatted) Update(val chat.AnyComponent) {
	f.text = f.text[:0]
	state := formatState{
		f: f,
	}
	state.build(val, func() chat.Color { return chat.White })
	f.Height = float64(state.lines+1) * 18
}

type formatState struct {
	f      *Formatted
	lines  int
	offset float64
}

func (f *formatState) build(c chat.AnyComponent, color getColorFunc) {
	switch c := c.Value.(type) {
	case *chat.TextComponent:
		gc := getColor(&c.Component, color)
		f.appendText(c.Text, gc)
		for _, e := range c.Extra {
			f.build(e, gc)
		}
	default:
		panic("unhandled component")
	}
}

func (f *formatState) appendText(text string, color getColorFunc) {
	width := 0.0
	last := 0
	for i, r := range text {
		s := render.SizeOfCharacter(r) + 2
		if (f.f.MaxWidth > 0 && f.offset+width+s > f.f.MaxWidth) || r == '\n' {
			rr, gg, bb := colorRGB(color())
			txt := NewText(text[last:i], f.offset, float64(f.lines*18+1), rr, gg, bb)
			txt.Parent = f.f
			last = i
			if r == '\n' {
				last++
			}
			f.f.text = append(f.f.text, txt)
			f.offset = 0
			f.lines++
			width = 0
		}
		width += s
	}
	if last != len(text) {
		r, g, b := colorRGB(color())
		txt := NewText(text[last:], f.offset, float64(f.lines*18+1), r, g, b)
		txt.Parent = f.f
		f.f.text = append(f.f.text, txt)
		f.offset += txt.Width + 2
	}
}

type getColorFunc func() chat.Color

func getColor(c *chat.Component, parent getColorFunc) getColorFunc {
	return func() chat.Color {
		if c.Color != "" {
			return c.Color
		}
		if parent != nil {
			return parent()
		}
		return chat.White
	}
}

func colorRGB(c chat.Color) (r, g, b int) {
	switch c {
	case chat.Black:
		return 0, 0, 0
	case chat.DarkBlue:
		return 0, 0, 170
	case chat.DarkGreen:
		return 0, 170, 0
	case chat.DarkAqua:
		return 0, 170, 170
	case chat.DarkRed:
		return 170, 0, 0
	case chat.DarkPurple:
		return 170, 0, 170
	case chat.Gold:
		return 255, 170, 0
	case chat.Gray:
		return 170, 170, 170
	case chat.DarkGray:
		return 85, 85, 85
	case chat.Blue:
		return 85, 85, 255
	case chat.Green:
		return 85, 255, 85
	case chat.Aqua:
		return 85, 255, 255
	case chat.Red:
		return 255, 85, 85
	case chat.LightPurple:
		return 255, 85, 255
	case chat.Yellow:
		return 255, 255, 85
	case chat.White:
		return 255, 255, 255
	}
	return 255, 255, 255
}
