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

	"github.com/thinkofdeath/steven/native"
	"github.com/thinkofdeath/steven/render/gl"
)

const (
	uiWidth, uiHeight = 854, 480
)

var (
	uiState = struct {
		program     gl.Program
		shader      *uiShader
		array       gl.VertexArray
		buffer      gl.Buffer
		indexBuffer gl.Buffer
		indexType   gl.Type
		maxIndex    int
		count       int
		data        []byte
		prevSize    int
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
	uiState.shader.Position.PointerInt(3, gl.Short, 24, 0)
	uiState.shader.TextureInfo.Pointer(4, gl.UnsignedShort, false, 24, 6)
	uiState.shader.TextureOffset.PointerInt(3, gl.Short, 24, 14)
	uiState.shader.Color.Pointer(4, gl.UnsignedByte, true, 24, 20)

	uiState.indexBuffer = gl.CreateBuffer()
	uiState.indexBuffer.Bind(gl.ElementArrayBuffer)
}

func drawUI() {
	// Prevent clipping with the world
	gl.Clear(gl.DepthBufferBit)
	gl.DepthFunc(gl.LessOrEqual)
	gl.Enable(gl.Blend)

	uiState.program.Use()
	uiState.shader.Texture.Int(0)
	if uiState.count > 0 {
		uiState.array.Bind()
		if uiState.maxIndex < uiState.count {
			var data []byte
			data, uiState.indexType = genElementBuffer(uiState.count)
			uiState.indexBuffer.Bind(gl.ElementArrayBuffer)
			uiState.indexBuffer.Data(data, gl.DynamicDraw)
			uiState.maxIndex = uiState.count
		}

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
		gl.DrawElements(gl.Triangles, uiState.count, uiState.indexType, 0)
	}
	gl.Disable(gl.Blend)
	gl.DepthFunc(gl.Less)
	uiState.count = 0
	uiState.data = uiState.data[:0]
}

// UIElement is a single element on the screen. It is a rectangle
// with a texture and a tint.
type UIElement struct {
	X, Y, W, H                 float64
	Layer                      int
	TX, TY, TW, TH             uint16
	TOffsetX, TOffsetY, TAtlas int16
	TSizeW, TSizeH             int16
	R, G, B, A                 byte
	Rotation                   float64
}

func UIAddBytes(data []byte) {
	uiState.data = append(uiState.data, data...)
	uiState.count += (len(data) / (24 * 4)) * 6
}

func NewUIElement(tex TextureInfo, x, y, width, height float64, tx, ty, tw, th float64) *UIElement {
	rect := tex.Rect()
	return &UIElement{
		X:        x / uiWidth,
		Y:        y / uiHeight,
		W:        width / uiWidth,
		H:        height / uiHeight,
		TX:       uint16(rect.X),
		TY:       uint16(rect.Y),
		TW:       uint16(rect.Width),
		TH:       uint16(rect.Height),
		TAtlas:   int16(tex.Atlas()),
		TOffsetX: int16(tx * float64(rect.Width) * 16),
		TOffsetY: int16(ty * float64(rect.Height) * 16),
		TSizeW:   int16(tw * float64(rect.Width) * 16),
		TSizeH:   int16(th * float64(rect.Height) * 16),
		R:        255,
		G:        255,
		B:        255,
		A:        255,
		Rotation: 0,
	}
}

func (u *UIElement) Bytes() []byte {
	data := make([]byte, 0, 24*4)
	data = u.appendVertex(data, u.X, u.Y, u.TOffsetX, u.TOffsetY)
	data = u.appendVertex(data, u.X+u.W, u.Y, u.TOffsetX+u.TSizeW, u.TOffsetY)
	data = u.appendVertex(data, u.X, u.Y+u.H, u.TOffsetX, u.TOffsetY+u.TSizeH)
	data = u.appendVertex(data, u.X+u.W, u.Y+u.H, u.TOffsetX+u.TSizeW, u.TOffsetY+u.TSizeH)
	return data
}

func (u *UIElement) appendVertex(data []byte, x, y float64, tx, ty int16) []byte {
	dx, dy := float64(x), float64(y)
	if u.Rotation != 0 {
		c := math.Cos(u.Rotation)
		s := math.Sin(u.Rotation)
		tmpx := dx - u.X - (u.W / 2)
		tmpy := dy - u.Y - (u.H / 2)
		dx = (u.W / 2) + (tmpx*c - tmpy*s) + u.X
		dy = (u.H / 2) + (tmpy*c + tmpx*s) + u.Y
	}
	data = appendShort(data, int16(math.Floor((dx*float64(lastWidth))+0.5)))
	data = appendShort(data, int16(math.Floor((dy*float64(lastHeight))+0.5)))
	data = appendShort(data, int16(u.Layer*256))
	data = appendUnsignedShort(data, u.TX)
	data = appendUnsignedShort(data, u.TY)
	data = appendUnsignedShort(data, u.TW)
	data = appendUnsignedShort(data, u.TH)
	data = appendShort(data, tx)
	data = appendShort(data, ty)
	data = appendShort(data, u.TAtlas)
	data = appendUnsignedByte(data, u.R)
	data = appendUnsignedByte(data, u.G)
	data = appendUnsignedByte(data, u.B)
	data = appendUnsignedByte(data, u.A)
	return data
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
