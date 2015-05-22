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

	"github.com/go-gl/mathgl/mgl32"
	"github.com/thinkofdeath/steven/native"
	"github.com/thinkofdeath/steven/render"
)

type Model struct {
	baseElement
	x, y  float64
	verts []*ModelVertex
	mat   mgl32.Mat4
}

type ModelVertex struct {
	X, Y, Z                    float32
	TX, TY, TW, TH             uint16
	TOffsetX, TOffsetY, TAtlas int16
	R, G, B, A                 byte
}

func NewModel(x, y float64, verts []*ModelVertex, mat mgl32.Mat4) *Model {
	return &Model{
		x: x, y: y,
		verts: verts,
		mat:   mat,
		baseElement: baseElement{
			visible: true,
			isNew:   true,
		},
	}
}

// Attach changes the location where this is attached to.
func (m *Model) Attach(vAttach, hAttach AttachPoint) *Model {
	m.vAttach, m.hAttach = vAttach, hAttach
	return m
}

func (m *Model) X() float64 { return m.x }
func (m *Model) SetX(x float64) {
	if m.x != x {
		m.x = x
		m.dirty = true
	}
}
func (m *Model) Y() float64 { return m.y }
func (m *Model) SetY(y float64) {
	if m.y != y {
		m.y = y
		m.dirty = true
	}
}

// Draw draws this to the target region.
func (m *Model) Draw(r Region, delta float64) {
	if m.isNew || m.isDirty() || forceDirty {
		m.isNew = false
		data := m.data[:0]

		for _, v := range m.verts {
			vec := mgl32.Vec3{v.X - 0.5, v.Y - 0.5, v.Z - 0.5}
			vec = m.mat.Mul4x1(vec.Vec4(1)).Vec3().
				Add(mgl32.Vec3{0.5, 0.5, 0.5})
			vX, vY, vZ := vec[0], 1.0-vec[1], vec[2]

			dx := r.X + r.W*float64(vX)
			dy := r.Y + r.H*float64(vY)
			dx /= scaledWidth
			dy /= scaledHeight
			data = appendShort(data, int16(math.Floor((dx*float64(lastWidth))+0.5)))
			data = appendShort(data, int16(math.Floor((dy*float64(lastHeight))+0.5)))
			data = appendShort(data, 256*int16(m.layer)+int16(256*vZ))
			data = appendUnsignedShort(data, v.TX)
			data = appendUnsignedShort(data, v.TY)
			data = appendUnsignedShort(data, v.TW)
			data = appendUnsignedShort(data, v.TH)
			data = appendShort(data, v.TOffsetX)
			data = appendShort(data, v.TOffsetY)
			data = appendShort(data, v.TAtlas)
			data = appendUnsignedByte(data, v.R)
			data = appendUnsignedByte(data, v.G)
			data = appendUnsignedByte(data, v.B)
			data = appendUnsignedByte(data, v.A)
		}
		m.data = data
	}
	render.UIAddBytes(m.data)
}

// Offset returns the offset of this drawable from the attachment
// point.
func (m *Model) Offset() (float64, float64) {
	return m.x, m.y
}

// Size returns the size of this drawable.
func (m *Model) Size() (float64, float64) {
	return 32, 32
}

// Remove removes the image element from the draw list.
func (m *Model) Remove() {
	Remove(m)
}

func appendUnsignedByte(data []byte, i byte) []byte {
	return append(data, i)
}

func appendByte(data []byte, i int8) []byte {
	return appendUnsignedByte(data, byte(i))
}

var scratch [8]byte

func appendUnsignedShort(data []byte, i uint16) []byte {
	d := scratch[:2]
	native.Order.PutUint16(d, i)
	return append(data, d...)
}

func appendShort(data []byte, i int16) []byte {
	return appendUnsignedShort(data, uint16(i))
}

func appendFloat(data []byte, f float32) []byte {
	d := scratch[:4]
	i := math.Float32bits(f)
	native.Order.PutUint32(d, i)
	return append(data, d...)
}
