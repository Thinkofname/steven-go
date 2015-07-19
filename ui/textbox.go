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
	"strings"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/render"
)

// TextBox is a drawable that allows for input.
type TextBox struct {
	baseElement
	btn                 *Button
	text                *Text
	x, y, w, h          float64
	input               string
	focused, wasFocused bool
	password            bool
	cursorTick          float64
	SubmitFunc          func()
}

// NewTextBox creates a new TextBox drawable.
func NewTextBox(x, y, w, h float64) *TextBox {
	t := &TextBox{
		x: x, y: y, w: w, h: h,
		baseElement: baseElement{
			visible: true,
			isNew:   true,
		},
	}
	btn := NewButton(0, 0, w, h)
	btn.SetDisabled(true)
	btn.AttachTo(t)
	text := NewText("", 5, 0, 255, 255, 255).Attach(Middle, Left)
	text.AttachTo(t)
	t.btn, t.text = btn, text
	return t
}

// Attach changes the location where this is attached to.
func (t *TextBox) Attach(vAttach, hAttach AttachPoint) *TextBox {
	t.vAttach, t.hAttach = vAttach, hAttach
	return t
}

func (t *TextBox) Value() string { return t.input }
func (t *TextBox) X() float64    { return t.x }
func (t *TextBox) SetX(x float64) {
	if t.x != x {
		t.x = x
		t.dirty = true
	}
}
func (t *TextBox) Y() float64 { return t.y }
func (t *TextBox) SetY(y float64) {
	if t.y != y {
		t.y = y
		t.dirty = true
	}
}
func (t *TextBox) Width() float64 { return t.w }
func (t *TextBox) SetWidth(w float64) {
	if t.w != w {
		t.w = w
		t.dirty = true
	}
}
func (t *TextBox) Height() float64 { return t.h }
func (t *TextBox) SetHeight(h float64) {
	if t.h != h {
		t.h = h
		t.dirty = true
	}
}
func (t *TextBox) Password() bool { return t.password }
func (t *TextBox) SetPassword(p bool) {
	if t.password != p {
		t.password = p
		t.dirty = true
	}
}

// Update updates the string drawn by this drawable.
func (t *TextBox) Update(val string) {
	t.input = val
	t.text.Update(t.value())
}

func (t *TextBox) tick(delta float64) {
	if !t.focused {
		if t.wasFocused {
			t.wasFocused = t.focused
			t.text.Update(t.value())
		}
		return
	}
	t.wasFocused = true
	t.cursorTick += delta
	if int(t.cursorTick/30)%2 == 0 {
		t.text.Update(t.value() + "|")
	} else {
		t.text.Update(t.value())
	}
	// Lazy way of preventing rounding errors buiding up over time
	if t.cursorTick > 0xFFFFFF {
		t.cursorTick = 0
	}
}

// Draw draws this to the target region.
func (t *TextBox) Draw(r Region, delta float64) {
	t.tick(delta)
	if t.isNew || t.isDirty() || forceDirty {
		t.isNew = false
		cw, ch := t.Size()
		sx, sy := r.W/cw, r.H/ch
		t.data = t.data[:0]

		r := getDrawRegion(t.btn, sx, sy)
		t.SetLayer(t.layer)
		t.btn.dirty = true
		t.btn.Draw(r, delta)
		t.data = append(t.data, t.btn.data...)

		r = getDrawRegion(t.text, sx, sy)
		t.SetLayer(t.layer)
		t.text.dirty = true
		t.text.Draw(r, delta)
		t.data = append(t.data, t.text.data...)
	}
	render.UIAddBytes(t.data)
}

func (t *TextBox) setFocused(f bool) {
	t.focused = f
}

func (t *TextBox) Click(r Region, x, y float64) {
	focus(t)
}
func (t *TextBox) Hover(r Region, x, y float64, over bool) {}

func (t *TextBox) value() string {
	if t.password {
		return strings.Repeat("*", len(t.input))
	}
	return t.input
}

func (t *TextBox) handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyBackspace && action != glfw.Release {
		if len(t.input) > 0 {
			t.input = t.input[:len(t.input)-1]
			t.text.Update(t.value())
		}
	}
	if key == glfw.KeyEnter && action == glfw.Release {
		if t.SubmitFunc != nil {
			t.SubmitFunc()
			return
		}
		CycleFocus()
	}
}

func (t *TextBox) handleChar(w *glfw.Window, char rune) {
	t.input += string(char)
	t.text.Update(t.value())
}

func (t *TextBox) isDirty() bool {
	return t.baseElement.isDirty() || t.text.dirty || t.btn.dirty
}

func (t *TextBox) clearDirty() {
	t.dirty = false
	t.text.dirty = false
	t.btn.dirty = false
}

// Offset returns the offset of this drawable from the attachment
// point.
func (t *TextBox) Offset() (float64, float64) {
	return t.x, t.y
}

// Size returns the size of this drawable.
func (t *TextBox) Size() (float64, float64) {
	return t.w, t.h
}

// Remove removes the textbox element from the draw list.
func (t *TextBox) Remove() {
	Remove(t)
}
