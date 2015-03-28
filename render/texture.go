package render

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"sync"

	"github.com/thinkofdeath/steven/platform/gl"
	"github.com/thinkofdeath/steven/render/atlas"
	"github.com/thinkofdeath/steven/resource"
)

var (
	textures    []*atlas.Type
	textureMap  = map[string]TextureInfo{}
	textureLock sync.RWMutex
)

const atlasSize = 1024

// TextureInfo returns information about a texture in an atlas
type TextureInfo struct {
	Atlas int
	*atlas.Rect
}

// GetTexture returns the related TextureInfo for the requested texture.
// If the texture isn't found a placeholder is returned instead.
func GetTexture(name string) TextureInfo {
	textureLock.RLock()
	defer textureLock.RUnlock()
	t, ok := textureMap[name]
	if !ok {
		return textureMap["missing_texture"]
	}
	return t
}

// TODO(Think) better error handling (if possible to recover?)
// TODO(Think) Store textures
func loadTextures() {
	textureLock.Lock()
	defer textureLock.Unlock()

	// Clear existing
	textures = nil
	textureMap = map[string]TextureInfo{}

	for _, file := range resource.Search("minecraft", "textures/blocks/", ".png") {
		r, err := resource.Open("minecraft", file)
		if err != nil {
			panic(err)
		}
		img, err := png.Decode(r)
		if err != nil {
			panic(err)
		}
		width, height := img.Bounds().Dx(), img.Bounds().Dy()
		if width != height {
			fmt.Printf("Skipping %s for now...\n", file)
			continue
		}
		var pix []byte
		switch img := img.(type) {
		case *image.NRGBA:
			pix = img.Pix
		case *image.RGBA:
			pix = img.Pix
		default:
			panic(fmt.Sprintf("unsupported image type %T", img))
		}
		name := file[len("textures/blocks/") : len(file)-4]
		at, rect := addTexture(pix, width, height)
		textureMap[name] = TextureInfo{
			Rect:  rect,
			Atlas: at,
		}
	}

	at, rect := addTexture([]byte{
		0, 0, 0, 255,
		255, 0, 255, 255,
		255, 0, 255, 255,
		0, 0, 0, 255,
	}, 2, 2)
	textureMap["missing_texture"] = TextureInfo{
		Rect:  rect,
		Atlas: at,
	}

	img := image.NewRGBA(image.Rect(0, 0, 1024, 1024))
	img.Pix = textures[0].Buffer
	var buf bytes.Buffer
	png.Encode(&buf, img)
	ioutil.WriteFile(".steven/atlas.png", buf.Bytes(), 0777)
}

func addTexture(pix []byte, width, height int) (int, *atlas.Rect) {
	for i, a := range textures {
		rect, err := a.Add(pix, width, height)
		if err == nil {
			return i, rect
		}
	}
	a := atlas.New(atlasSize, atlasSize, 4)
	textures = append(textures, a)
	rect, err := a.Add(pix, width, height)
	if err != nil {
		panic(err)
	}
	return len(textures) - 1, rect
}

type glTexture struct {
	Data          []byte
	Width, Height int
	Format        gl.TextureFormat
	Type          gl.Type
	Filter        gl.TextureValue
	MinFilter     gl.TextureValue
	Wrap          gl.TextureValue
}

func createTexture(t glTexture) gl.Texture {
	if t.Format == 0 {
		t.Format = gl.RGB
	}
	if t.Type == 0 {
		t.Type = gl.UnsignedByte
	}
	if t.Filter == 0 {
		t.Filter = gl.Nearest
	}
	if t.MinFilter == 0 {
		t.MinFilter = t.Filter
	}
	if t.Wrap == 0 {
		t.Wrap = gl.ClampToEdge
	}

	texture := gl.CreateTexture()
	texture.Bind(gl.Texture2D)
	texture.Image2D(0, t.Width, t.Height, t.Format, t.Type, t.Data)
	texture.Parameter(gl.TextureMagFilter, t.Filter)
	texture.Parameter(gl.TextureMinFilter, t.MinFilter)
	texture.Parameter(gl.TextureWrapS, t.Wrap)
	texture.Parameter(gl.TextureWrapT, t.Wrap)
	return texture
}
