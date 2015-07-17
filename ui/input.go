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

// focusable is a drawable that can be focsued for keyboard input
type focusable interface {
	Drawable
	setFocused(bool)
	handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
	handleChar(w *glfw.Window, char rune)
}

// HandleKey passes the input to the focused drawable
func HandleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyTab && action == glfw.Release {
		CycleFocus()
		return
	}
	if f := getFocused(); f != nil {
		f.handleKey(w, key, scancode, action, mods)
	}
}

// HandleChar passes the input to the focused drawable
func HandleChar(w *glfw.Window, char rune) {
	if f := getFocused(); f != nil {
		f.handleChar(w, char)
	}
}

// Currently focused drawable
var focused focusable

// Returns the currently focused drawable or nil
// Tries to focus an drawable if one isn't currently focused
func getFocused() focusable {
	if focused != nil {
		// Ensure the drawable is focused
		for _, d := range drawables {
			if d.Drawable == focused {
				return d.Drawable.(focusable)
			}
		}
	}
	// Try to focus another drawable
	for _, d := range drawables {
		if f, ok := d.Drawable.(focusable); ok {
			focus(f)
			return focused
		}
	}
	// Clear the focus if one isn't found incase this
	// fell through from the first check
	focus(nil)
	return nil
}

// CycleFocus changes the focus to the next drawable
func CycleFocus() {
	pos := 0
	l := len(drawables)
	// Find our drawable
	for i, d := range drawables {
		if d.Drawable == focused {
			pos = (i + 1) % l
			break
		}
	}
	// Make sure we don't loop forever on scenes without focusable drawables
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

// Changes the focused drawable
func focus(f focusable) {
	if focused != nil {
		focused.setFocused(false)
	}
	focused = f
	if focused != nil {
		focused.setFocused(true)
	}
}
