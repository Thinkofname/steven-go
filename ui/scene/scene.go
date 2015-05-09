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

// Package scene provides methods to manage multiple ui scenes.
package scene

import "github.com/thinkofdeath/phteven/ui"

// Type stores a scene that can be removed and shown at any time.
type Type struct {
	visible bool

	drawables []ui.Drawable
}

// New creates a new scene.
func New(visible bool) *Type {
	return &Type{
		visible: visible,
	}
}

// Show shows all the drawables in the scene
func (t *Type) Show() {
	if t.visible {
		return
	}
	t.visible = true
	for _, d := range t.drawables {
		ui.AddDrawable(d)
	}
}

// Hide hides all the drawables in the scene
func (t *Type) Hide() {
	if !t.visible {
		return
	}
	t.visible = false
	for _, d := range t.drawables {
		ui.Remove(d)
	}
}

// AddDrawable adds the drawable to the draw list.
func (t *Type) AddDrawable(d ui.Drawable) {
	t.drawables = append(t.drawables, d)
	if t.visible {
		ui.AddDrawable(d)
	}
}

// IsVisible returns whether the scene is currently visible.
func (t *Type) IsVisible() bool {
	return t.visible
}
