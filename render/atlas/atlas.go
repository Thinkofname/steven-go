// Package atlas provides a basic texture atlas
package atlas

import (
	"errors"
	"math"
)

var (
	// ErrAtlasFull is returned when an image can't fit in an atlas
	ErrAtlasFull = errors.New("atlas full")
)

// Type is a texture atlas for storing quake
// pictures (from the bsp package). The buffer
// is public to allow easy uploading to the
// gpu.
type Type struct {
	width, height int
	Buffer        []byte
	pixelSize     int
	freeSpace     []*Rect
	padding       int
}

// Rect represents a location in a texture
// atlas.
type Rect struct {
	X, Y          int
	Width, Height int
}

// New creates an atlas of the specified size
// with zero padding around the textures.
// pixelSize controls the number of bytes per
// a pixel
func New(width, height, pixelSize int) *Type {
	return NewPadded(width, height, pixelSize, 0)
}

// NewPadded creates an atlas of the specified
// size. Textures are padded with the passed
// number of pixels around each size. This is
// useful for filtering textures without other
// textures bleeding through.
// pixelSize controls the number of bytes per
// a pixel
func NewPadded(width, height, pixelSize, padding int) *Type {
	a := &Type{
		width:     width,
		height:    height,
		padding:   padding,
		pixelSize: pixelSize,
		Buffer:    make([]byte, width*height*pixelSize),
	}
	a.freeSpace = append(a.freeSpace, &Rect{
		X:      0,
		Y:      0,
		Width:  width,
		Height: height,
	})
	return a
}

// Add adds the passed texture to the atlas and
// returns the location in the atlas. This method
// panics if the atlas has been baked.
func (a *Type) Add(image []byte, width, height int) (*Rect, error) {
	// Double the padding since its for both
	// sides
	w := width + (a.padding * 2)
	h := height + (a.padding * 2)

	var target *Rect
	targetIndex := 0
	priority := math.MaxInt32
	// Search through and find the best fit for this texture
	for i, free := range a.freeSpace {
		if free.Width >= w && free.Height >= h {
			currentPriority := (free.Width - w) * (free.Height - h)
			if target == nil || currentPriority < priority {
				target = free
				priority = currentPriority
				targetIndex = i
			}

			// Perfect match, we can break early
			if priority == 0 {
				break
			}
		}
	}

	if target == nil {
		return nil, ErrAtlasFull
	}

	// Copy the image into the atlas
	CopyImage(image, a.Buffer, target.X, target.Y, w, h, a.width, a.height, a.pixelSize, a.padding)

	tx := target.X + a.padding
	ty := target.Y + a.padding

	if w == target.Width {
		target.Y += h
		target.Height -= h
		if target.Height == 0 {
			// Remove empty sections
			a.freeSpace = append(a.freeSpace[:targetIndex], a.freeSpace[targetIndex+1:]...)
		}
	} else {
		if target.Height > h {
			// Split by height
			a.freeSpace = append(
				[]*Rect{&Rect{
					target.X, target.Y + h,
					w, target.Height - h,
				}},
				a.freeSpace...,
			)
		}
		target.X += w
		target.Width -= w
	}

	return &Rect{
		X:      tx,
		Y:      ty,
		Width:  width,
		Height: height,
	}, nil
}

// helper method that allows for out of bounds access
// to a picture. The coordinates to be changed to the
// nearest edge when out of bounds.
func safeGetPixel(x, y, w, h int) int {
	if x < 0 {
		x = 0
	}
	if x >= w {
		x = w - 1
	}
	if y < 0 {
		y = 0
	}
	if y >= h {
		y = h - 1
	}
	return y*w + x
}

// CopyImage copies the passed image data to the
// target buffer, accounting for padding.
func CopyImage(data, buffer []byte, targetX, targetY, w, h, width, height, pixelSize, padding int) {
	for y := 0; y < h; y++ {
		index := (targetY+y)*width + targetX
		for x := 0; x < w; x++ {
			px := x - padding
			py := y - padding
			pw := w - (padding << 1)
			ph := h - (padding << 1)
			bo := (index + x) * pixelSize
			do := safeGetPixel(px, py, pw, ph) * pixelSize
			copy(buffer[bo:bo+pixelSize], data[do:do+pixelSize])
		}
	}
}
