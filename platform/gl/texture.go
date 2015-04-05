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

package gl

import (
	"github.com/thinkofdeath/gl/v3.2-core/gl"
)

const (
	Texture2D TextureTarget = gl.TEXTURE_2D

	RGB   TextureFormat = gl.RGB
	RGBA  TextureFormat = gl.RGBA
	RGBA8 TextureFormat = gl.RGBA8

	TextureMinFilter TextureParameter = gl.TEXTURE_MIN_FILTER
	TextureMagFilter TextureParameter = gl.TEXTURE_MAG_FILTER
	TextureWrapS     TextureParameter = gl.TEXTURE_WRAP_S
	TextureWrapT     TextureParameter = gl.TEXTURE_WRAP_T

	Nearest              TextureValue = gl.NEAREST
	Linear               TextureValue = gl.LINEAR
	LinearMipmapLinear   TextureValue = gl.LINEAR_MIPMAP_LINEAR
	LinearMipmapNearest  TextureValue = gl.LINEAR_MIPMAP_NEAREST
	NearestMipmapNearest TextureValue = gl.NEAREST_MIPMAP_NEAREST
	NearestMipmapLinear  TextureValue = gl.NEAREST_MIPMAP_LINEAR
	ClampToEdge          TextureValue = gl.CLAMP_TO_EDGE
)

// State tracking
var (
	currentTexture       Texture
	currentTextureTarget TextureTarget
)

type TextureTarget uint32
type TextureFormat uint32
type TextureParameter uint32
type TextureValue int32

type Texture struct {
	internal uint32
}

func CreateTexture() Texture {
	var texture Texture
	gl.GenTextures(1, &texture.internal)
	return texture
}

func (t Texture) Bind(target TextureTarget) {
	if currentTexture == t && currentTextureTarget == target {
		return
	}
	gl.BindTexture(uint32(target), t.internal)
	currentTexture = t
	currentTextureTarget = target
}

func (t Texture) Image2D(level int, width, height int, format TextureFormat, ty Type, pix []byte) {
	if t != currentTexture {
		panic("texture not bound")
	}
	gl.TexImage2D(
		uint32(currentTextureTarget),
		int32(level),
		int32(format),
		int32(width),
		int32(height),
		0,
		uint32(format),
		uint32(ty),
		gl.Ptr(pix),
	)
}

func (t Texture) SubImage2D(level int, x, y, width, height int, format TextureFormat, ty Type, pix []byte) {
	if t != currentTexture {
		panic("texture not bound")
	}
	gl.TexSubImage2D(
		uint32(currentTextureTarget),
		int32(level),
		int32(x),
		int32(y),
		int32(width),
		int32(height),
		uint32(format),
		uint32(ty),
		gl.Ptr(pix),
	)
}

func (t Texture) Parameter(param TextureParameter, val TextureValue) {
	if t != currentTexture {
		panic("texture not bound")
	}
	gl.TexParameteri(uint32(currentTextureTarget), uint32(param), int32(val))
}
