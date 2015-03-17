// +build mobile

package gl

import (
	"golang.org/x/mobile/gl"
)

const (
	Texture2D TextureTarget = gl.TEXTURE_2D

	RGB TextureFormat = gl.RGB

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

type Texture gl.Texture

func CreateTexture() Texture {
	return Texture(gl.GenTexture())
}

func (t Texture) Bind(target TextureTarget) {
	if currentTexture == t && currentTextureTarget == target {
		return
	}
	gl.BindTexture(gl.Enum(target), gl.Texture(t))
	currentTexture = t
	currentTextureTarget = target
}

func (t Texture) Image2D(level int, width, height int, format TextureFormat, ty Type, pix []byte) {
	if t != currentTexture {
		panic("texture not bound")
	}
	gl.TexImage2D(
		gl.Enum(currentTextureTarget),
		level,
		width,
		height,
		gl.Enum(format),
		gl.Enum(ty),
		pix,
	)
}

func (t Texture) Parameter(param TextureParameter, val TextureValue) {
	if t != currentTexture {
		panic("texture not bound")
	}
	gl.TexParameteri(gl.Enum(currentTextureTarget), gl.Enum(param), int(val))
}
