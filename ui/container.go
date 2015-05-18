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
	baseElement
	x, y, h, w float64
	hovered    bool
	HoverFunc  func(over bool)
	ClickFunc  func()
}

func NewContainer(x, y, w, h float64) *Container {
	return &Container{
		x: x, y: y, w: w, h: h,
		baseElement: baseElement{
			visible: false,
			isNew:   false,
		},
	}
}

func (c *Container) X() float64 { return c.x }
func (c *Container) SetX(x float64) {
	if c.x != x {
		c.x = x
		c.dirty = true
	}
}
func (c *Container) Y() float64 { return c.y }
func (c *Container) SetY(y float64) {
	if c.y != y {
		c.y = y
		c.dirty = true
	}
}
func (c *Container) Width() float64 { return c.w }
func (c *Container) SetWidth(w float64) {
	if c.w != w {
		c.w = w
		c.dirty = true
	}
}
func (c *Container) Height() float64 { return c.h }
func (c *Container) SetHeight(h float64) {
	if c.h != h {
		c.h = h
		c.dirty = true
	}
}

// Attach changes the location where this is attached to.
func (c *Container) Attach(vAttach, hAttach AttachPoint) *Container {
	c.vAttach, c.hAttach = vAttach, hAttach
	return c
}

// Draw draws this to the target region.
func (c *Container) Draw(r Region, delta float64) {
}

// Offset returns the offset of this drawable from the attachment
// point.
func (c *Container) Offset() (float64, float64) {
	return c.x, c.y
}

// Size returns the size of this drawable.
func (c *Container) Size() (float64, float64) {
	return c.w, c.h
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
