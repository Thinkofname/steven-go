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

// Container is a drawable that is used for aligning other drawables.
// Should never be drawn.
type Container struct {
	Parent           Drawable
	X, Y, W, H       float64
	vAttach, hAttach AttachPoint
	hovered          bool
	HoverFunc        func(over bool)
	ClickFunc        func()
}

// Attach changes the location where this is attached to.
func (c *Container) Attach(vAttach, hAttach AttachPoint) *Container {
	c.vAttach, c.hAttach = vAttach, hAttach
	return c
}

// Attachment returns the sides where this element is attached too.
func (c *Container) Attachment() (vAttach, hAttach AttachPoint) {
	return c.vAttach, c.hAttach
}

// ShouldDraw returns whether this should be drawn at this time.
func (c *Container) ShouldDraw() bool {
	return false
}

// Draw draws this to the target region.
func (c *Container) Draw(r Region, delta float64) {
}

// AttachedTo returns the Drawable this is attached to or nil.
func (c *Container) AttachedTo() Drawable {
	return c.Parent
}

// Offset returns the offset of this drawable from the attachment
// point.
func (c *Container) Offset() (float64, float64) {
	return c.X, c.Y
}

// Size returns the size of this drawable.
func (c *Container) Size() (float64, float64) {
	return c.W, c.H
}
func (c *Container) Click(r Region, x, y float64) {
	if c.ClickFunc != nil {
		c.ClickFunc()
	}
}
func (c *Container) Hover(r Region, x, y float64, over bool) {
	c.hovered = over
	if c.HoverFunc != nil {
		c.HoverFunc(over)
	}
}
