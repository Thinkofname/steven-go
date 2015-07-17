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

import "github.com/go-gl/glfw/v3.1/glfw"

type focusable interface {
	Drawable
	setFocused(bool)
	handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
	handleChar(w *glfw.Window, char rune)
}

func HandleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyTab && action == glfw.Release {
		CycleFocus()
		return
	}
	if f := getFocused(); f != nil {
		f.handleKey(w, key, scancode, action, mods)
	}
}

func HandleChar(w *glfw.Window, char rune) {
	if f := getFocused(); f != nil {
		f.handleChar(w, char)
	}
}

var focused focusable

func getFocused() focusable {
	if focused != nil {
		for _, d := range drawables {
			if d.Drawable == focused {
				return d.Drawable.(focusable)
			}
		}
	}
	for _, d := range drawables {
		if f, ok := d.Drawable.(focusable); ok {
			focus(f)
			return focused
		}
	}
	return nil
}

func CycleFocus() {
	pos := 0
	l := len(drawables)
	for i, d := range drawables {
		if d.Drawable == focused {
			pos = (i + 1) % l
			break
		}
	}
	max := l
	for ; max >= 0; max-- {
		dr := drawables[pos]
		if f, ok := dr.Drawable.(focusable); ok {
			focus(f)
			break
		}
		pos = (pos + 1) % l
	}
}

func focus(f focusable) {
	if focused != nil {
		focused.setFocused(false)
	}
	focused = f
	focused.setFocused(true)
}
