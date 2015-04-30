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

// Image is a drawable that draws a texture.
type Image struct {
	Parent           Drawable
	Texture          *render.TextureInfo
	X, Y, W, H       float64
	TX, TY, TW, TH   float64
	R, G, B, A       int
	Visible          bool
	vAttach, hAttach AttachPoint
}

// NewImage creates a new image drawable.
func NewImage(texture *render.TextureInfo, x, y, w, h, tx, ty, tw, th float64, r, g, b int) *Image {
	return &Image{
		Texture: texture,
		R:       r, G: g, B: b, A: 255,
		X: x, Y: y, W: w, H: h,
		TX: tx, TY: ty, TW: tw, TH: th,
		Visible: true,
	}
}

// Attach changes the location where this is attached to.
func (i *Image) Attach(vAttach, hAttach AttachPoint) *Image {
	i.vAttach, i.hAttach = vAttach, hAttach
	return i
}

// Attachment returns the sides where this element is attached too.
func (i *Image) Attachment() (vAttach, hAttach AttachPoint) {
	return i.vAttach, i.hAttach
}

// ShouldDraw returns whether this should be drawn at this time.
func (i *Image) ShouldDraw() bool {
	return i.Visible
}

// Draw draws this to the target region.
func (i *Image) Draw(r Region, delta float64) {
	e := render.DrawUIElement(i.Texture, r.X, r.Y, r.W, r.H, i.TX, i.TY, i.TW, i.TH)
	e.R = byte(i.R)
	e.G = byte(i.G)
	e.B = byte(i.B)
	e.A = byte(i.A)
}

// AttachedTo returns the Drawable this is attached to or nil.
func (i *Image) AttachedTo() Drawable {
	return i.Parent
}

// Offset returns the offset of this drawable from the attachment
// point.
func (i *Image) Offset() (float64, float64) {
	return i.X, i.Y
}

// Size returns the size of this drawable.
func (i *Image) Size() (float64, float64) {
	return i.W, i.H
}

// Remove removes the image element from the draw list.
func (i *Image) Remove() {
	Remove(i)
}
