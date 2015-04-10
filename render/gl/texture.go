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

// TextureTarget is a target were a texture can be bound to.
type TextureTarget uint32

// Valid texture targets.
const (
	Texture2D      TextureTarget = gl.TEXTURE_2D
	Texture2DArray TextureTarget = gl.TEXTURE_2D_ARRAY
)

// TextureFormat is the format of a texture either internally or
// to be uploaded.
type TextureFormat uint32

// Valid texture formats.
const (
	RGB   TextureFormat = gl.RGB
	RGBA  TextureFormat = gl.RGBA
	RGBA8 TextureFormat = gl.RGBA8
)

// TextureParameter is a parameter that can be read or set on a texture.
type TextureParameter uint32

// Valid texture parameters.
const (
	TextureMinFilter TextureParameter = gl.TEXTURE_MIN_FILTER
	TextureMagFilter TextureParameter = gl.TEXTURE_MAG_FILTER
	TextureWrapS     TextureParameter = gl.TEXTURE_WRAP_S
	TextureWrapT     TextureParameter = gl.TEXTURE_WRAP_T
	TextureMaxLevel  TextureParameter = gl.TEXTURE_MAX_LEVEL
)

// TextureValue is a value that be set on a texture's parameter.
type TextureValue int32

// Valid texture values.
const (
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

// Texture is a buffer of data used by fragment shaders to produce color.
type Texture struct {
	internal uint32
}

// CreateTexture allocates a new texture.
func CreateTexture() Texture {
	var texture Texture
	gl.GenTextures(1, &texture.internal)
	return texture
}

// Bind binds the texture to the passed target, if the texture is already bound
// then this does nothing.
func (t Texture) Bind(target TextureTarget) {
	if currentTexture == t && currentTextureTarget == target {
		return
	}
	gl.BindTexture(uint32(target), t.internal)
	currentTexture = t
	currentTextureTarget = target
}

// Image3D uploads a 3D texture to the GPU.
func (t Texture) Image3D(level, width, height, depth int, format TextureFormat, ty Type, pix []byte) {
	if t != currentTexture {
		panic("texture not bound")
	}
	gl.TexImage3D(uint32(currentTextureTarget), int32(level), int32(format), int32(width), int32(height), int32(depth), int32(0), uint32(format), uint32(ty), gl.Ptr(pix))
}

// SubImage3D updates a region of a 3D texture.
func (t Texture) SubImage3D(level, x, y, z, width, height, depth int, format TextureFormat, ty Type, pix []byte) {
	if t != currentTexture {
		panic("texture not bound")
	}
	gl.TexSubImage3D(uint32(currentTextureTarget), int32(level), int32(x), int32(y), int32(z), int32(width), int32(height), int32(depth), uint32(format), uint32(ty), gl.Ptr(pix))
}

// Image2D uploads a 2D texture to the GPU.
func (t Texture) Image2D(level, width, height int, format TextureFormat, ty Type, pix []byte) {
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

// SubImage2D updates a region of a 2D texture.
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

// Parameter sets a parameter on the texture to passed value.
func (t Texture) Parameter(param TextureParameter, val TextureValue) {
	if t != currentTexture {
		panic("texture not bound")
	}
	gl.TexParameteri(uint32(currentTextureTarget), uint32(param), int32(val))
}
