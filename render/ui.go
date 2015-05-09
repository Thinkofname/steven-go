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

package render

import (
	"math"

	"github.com/thinkofdeath/phteven/native"
	"github.com/thinkofdeath/phteven/render/gl"
)

const (
	uiWidth, uiHeight = 800, 480
)

var (
	uiState = struct {
		program      gl.Program
		shader       *uiShader
		array        gl.VertexArray
		buffer       gl.Buffer
		count        int
		data         []byte
		prevSize     int
		elements     []UIElement
		elementCount int
	}{
		prevSize: -1,
	}
)

func initUI() {
	uiState.program = CreateProgram(vertexUI, fragmentUI)
	uiState.shader = &uiShader{}
	InitStruct(uiState.shader, uiState.program)

	uiState.array = gl.CreateVertexArray()
	uiState.array.Bind()
	uiState.buffer = gl.CreateBuffer()
	uiState.buffer.Bind(gl.ArrayBuffer)
	uiState.shader.Position.Enable()
	uiState.shader.TextureInfo.Enable()
	uiState.shader.TextureOffset.Enable()
	uiState.shader.Color.Enable()
	uiState.shader.Position.PointerInt(2, gl.Short, 22, 0)
	uiState.shader.TextureInfo.Pointer(4, gl.UnsignedShort, false, 22, 4)
	uiState.shader.TextureOffset.PointerInt(3, gl.Short, 22, 12)
	uiState.shader.Color.Pointer(4, gl.UnsignedByte, true, 22, 18)
}

func drawUI() {
	// Redraw everything
	uiState.count = 0
	uiState.data = uiState.data[:0]
	for i := 0; i < uiState.elementCount; i++ {
		uiState.elements[i].draw()
	}
	uiState.elementCount = 0

	// Prevent clipping with the world
	gl.Disable(gl.DepthTest)
	gl.Enable(gl.Blend)

	uiState.program.Use()
	uiState.shader.Texture.Int(0)
	if uiState.count > 0 {
		uiState.array.Bind()

		uiState.shader.ScreenSize.Float2(float32(lastWidth), float32(lastHeight))

		uiState.buffer.Bind(gl.ArrayBuffer)
		if len(uiState.data) > uiState.prevSize {
			uiState.prevSize = len(uiState.data)
			uiState.buffer.Data(uiState.data, gl.StreamDraw)
		} else {
			target := uiState.buffer.Map(gl.WriteOnly, len(uiState.data))
			copy(target, uiState.data)
			uiState.buffer.Unmap()
		}
		gl.DrawArrays(gl.Triangles, 0, uiState.count)
	}
	gl.Disable(gl.Blend)
	gl.Enable(gl.DepthTest)
}

// UIElement is a single element on the screen. It is a rectangle
// with a texture and a tint.
type UIElement struct {
	X, Y, W, H                 float64
	TX, TY, TW, TH             uint16
	TOffsetX, TOffsetY, TAtlas int16
	TSizeW, TSizeH             int16
	R, G, B, A                 byte
	Rotation                   float64
}

// DrawUIElement draws a single ui element onto the screen.
func DrawUIElement(tex *TextureInfo, x, y, width, height float64, tx, ty, tw, th float64) *UIElement {
	if len(uiState.elements) == uiState.elementCount {
		old := uiState.elements
		uiState.elements = make([]UIElement, (len(old)+1)<<1)
		copy(uiState.elements, old)
	}
	e := &uiState.elements[uiState.elementCount]
	// (Re)set the information for the element
	e.X = x / uiWidth
	e.Y = y / uiHeight
	e.W = width / uiWidth
	e.H = height / uiHeight
	e.TX = uint16(tex.X)
	e.TY = uint16(tex.Y)
	e.TW = uint16(tex.Width)
	e.TH = uint16(tex.Height)
	e.TAtlas = int16(tex.Atlas)
	e.TOffsetX = int16(tx * float64(tex.Width) * 16)
	e.TOffsetY = int16(ty * float64(tex.Height) * 16)
	e.TSizeW = int16(tw * float64(tex.Width) * 16)
	e.TSizeH = int16(th * float64(tex.Height) * 16)
	e.R = 255
	e.G = 255
	e.B = 255
	e.A = 255
	uiState.elementCount++
	return e
}

// Shift moves the element by the passed amounts.
func (u *UIElement) Shift(x, y float64) {
	u.X += x / uiWidth
	u.Y += y / uiHeight
}

// Alpha changes the alpha of this element
func (u *UIElement) Alpha(a float64) {
	if a > 1.0 {
		a = 1.0
	}
	if a < 0.0 {
		a = 0.0
	}
	u.A = byte(255.0 * a)
}

func (u *UIElement) draw() {
	u.appendVertex(u.X, u.Y, u.TOffsetX, u.TOffsetY)
	u.appendVertex(u.X+u.W, u.Y, u.TOffsetX+u.TSizeW, u.TOffsetY)
	u.appendVertex(u.X, u.Y+u.H, u.TOffsetX, u.TOffsetY+u.TSizeH)

	u.appendVertex(u.X+u.W, u.Y+u.H, u.TOffsetX+u.TSizeW, u.TOffsetY+u.TSizeH)
	u.appendVertex(u.X, u.Y+u.H, u.TOffsetX, u.TOffsetY+u.TSizeH)
	u.appendVertex(u.X+u.W, u.Y, u.TOffsetX+u.TSizeW, u.TOffsetY)
}

func (u *UIElement) appendVertex(x, y float64, tx, ty int16) {
	uiState.count++
	dx, dy := float64(x), float64(y)
	if u.Rotation != 0 {
		c := math.Cos(u.Rotation)
		s := math.Sin(u.Rotation)
		tmpx := dx - u.X - (u.W / 2)
		tmpy := dy - u.Y - (u.H / 2)
		dx = (u.W / 2) + (tmpx*c - tmpy*s) + u.X
		dy = (u.H / 2) + (tmpy*c + tmpx*s) + u.Y
	}
	uiState.data = appendShort(uiState.data, int16(math.Floor((dx*float64(lastWidth))+0.5)))
	uiState.data = appendShort(uiState.data, int16(math.Floor((dy*float64(lastHeight))+0.5)))
	uiState.data = appendUnsignedShort(uiState.data, u.TX)
	uiState.data = appendUnsignedShort(uiState.data, u.TY)
	uiState.data = appendUnsignedShort(uiState.data, u.TW)
	uiState.data = appendUnsignedShort(uiState.data, u.TH)
	uiState.data = appendShort(uiState.data, tx)
	uiState.data = appendShort(uiState.data, ty)
	uiState.data = appendShort(uiState.data, u.TAtlas)
	uiState.data = appendUnsignedByte(uiState.data, u.R)
	uiState.data = appendUnsignedByte(uiState.data, u.G)
	uiState.data = appendUnsignedByte(uiState.data, u.B)
	uiState.data = appendUnsignedByte(uiState.data, u.A)
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
