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

package steven

import (
	"image/png"
	"math"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/resource"
	"github.com/thinkofdeath/steven/type/direction"
	"github.com/thinkofdeath/steven/world/biome"
)

func staticModelFromItem(mdl *model, block Block, mode string) (out []*render.StaticVertex, mat mgl32.Mat4) {
	mat = mgl32.Ident4()
	if gui, ok := mdl.display[mode]; ok {
		if gui.Scale != nil {
			mat = mat.Mul4(mgl32.Scale3D(
				float32(gui.Scale[0]),
				float32(gui.Scale[1]),
				float32(gui.Scale[2]),
			))
		}
		if gui.Translation != nil {
			mat = mat.Mul4(mgl32.Translate3D(
				float32(gui.Translation[0]/32),
				float32(-gui.Translation[1]/32),
				float32(-gui.Translation[2]/32),
			))
		}
		if gui.Rotation != nil {
			mat = mat.Mul4(mgl32.Rotate3DY(float32(gui.Rotation[1]/180) * math.Pi).Mat4())
			mat = mat.Mul4(mgl32.Rotate3DX(float32(gui.Rotation[0]/180) * math.Pi).Mat4())
			mat = mat.Mul4(mgl32.Rotate3DZ(float32(gui.Rotation[2]/180) * math.Pi).Mat4())
		}
	}

	p := precomputeModel(mdl)
	for fi := range p.faces {
		f := p.faces[len(p.faces)-1-fi]
		var cr, cg, cb byte
		cr = 255
		cg = 255
		cb = 255
		if block != nil && block.TintImage() != nil {
			switch f.tintIndex {
			case 0:
				bi := biome.Plains
				ix := bi.ColorIndex & 0xFF
				iy := bi.ColorIndex >> 8
				col := block.TintImage().NRGBAAt(ix, iy)
				cr = byte(col.R)
				cg = byte(col.G)
				cb = byte(col.B)
			}
		}
		if f.facing == direction.East || f.facing == direction.West {
			cr = byte(float64(cr) * 0.8)
			cg = byte(float64(cg) * 0.8)
			cb = byte(float64(cb) * 0.8)
		}
		if f.facing == direction.North || f.facing == direction.South {
			cr = byte(float64(cr) * 0.6)
			cg = byte(float64(cg) * 0.6)
			cb = byte(float64(cb) * 0.6)
		}

		for i, vert := range f.vertices {
			vX, vY, vZ := float32(vert.X)/256, float32(vert.Y)/256, float32(vert.Z)/256
			tex := f.verticesTexture[i]
			rect := tex.Rect()
			vert := &render.StaticVertex{
				X:        vX - 0.5,
				Y:        vY - 0.5,
				Z:        vZ - 0.5,
				Texture:  tex,
				TextureX: float64(vert.TOffsetX) / float64(16*rect.Width),
				TextureY: float64(vert.TOffsetY) / float64(16*rect.Height),
				R:        cr,
				G:        cg,
				B:        cb,
				A:        255,
			}
			out = append(out, vert)
		}
	}
	return
}

func genStaticModelFromItem(mdl *model, block Block, mode string) (out []*render.StaticVertex, mat mgl32.Mat4) {
	if mode == "thirdperson" {
		mat = mgl32.Translate3D(0, 0, 2/16.0).
			Mul4(mgl32.Rotate3DY(math.Pi).Mat4()).
			Mul4(mgl32.Rotate3DZ(math.Pi).Mat4())
	} else {
		mat = mgl32.Translate3D(0, -8/16.0, 0).
			Mul4(mgl32.Rotate3DX(math.Pi).Mat4()).
			Mul4(mgl32.Rotate3DY(math.Pi).Mat4()).
			Mul4(mgl32.Rotate3DZ(-0.6).Mat4()).
			Mul4(mgl32.Scale3D(1.1, 1.1, 1.1))
	}
	if gui, ok := mdl.display[mode]; ok {
		if gui.Scale != nil {
			mat = mat.Mul4(mgl32.Scale3D(
				float32(gui.Scale[0]),
				float32(gui.Scale[1]),
				float32(gui.Scale[2]),
			))
		}
		if gui.Translation != nil {
			mat = mat.Mul4(mgl32.Translate3D(
				float32(gui.Translation[0]/32),
				float32(-gui.Translation[1]/32),
				float32(-gui.Translation[2]/32),
			))
		}
		if gui.Rotation != nil {
			mat = mat.Mul4(mgl32.Rotate3DY(float32(gui.Rotation[1]/180) * math.Pi).Mat4())
			mat = mat.Mul4(mgl32.Rotate3DX(float32(gui.Rotation[0]/180) * math.Pi).Mat4())
			mat = mat.Mul4(mgl32.Rotate3DZ(float32(gui.Rotation[2]/180) * math.Pi).Mat4())
		}
	}

	tex := render.GetTexture("solid")
	rect := tex.Rect()

	tName, plugin := mdl.textureVars["layer0"], "minecraft"
	if pos := strings.IndexRune(tName, ':'); pos != -1 {
		plugin = tName[:pos]
		tName = tName[pos:]
	}
	f, err := resource.Open(plugin, "textures/"+tName+".png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}

	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	sx := 1 / float32(w)
	sy := 1 / float32(h)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			col := img.At(x, y)
			rr, gg, bb, aa := col.RGBA()
			if aa == 0 {
				continue
			}

			for i, f := range faceVertices {
				var cr, cg, cb byte
				cr = byte(rr >> 8)
				cg = byte(gg >> 8)
				cb = byte(bb >> 8)
				facing := direction.Type(i)
				if facing == direction.East || facing == direction.West {
					cr = byte(float64(cr) * 0.8)
					cg = byte(float64(cg) * 0.8)
					cb = byte(float64(cb) * 0.8)
				}
				if facing == direction.North || facing == direction.South {
					cr = byte(float64(cr) * 0.6)
					cg = byte(float64(cg) * 0.6)
					cb = byte(float64(cb) * 0.6)
				}

				for _, vert := range f.verts {
					vX, vY, vZ := float32(vert.X), float32(vert.Y), float32(vert.Z)
					vert := &render.StaticVertex{
						Y:        vY*sy - 0.5 + sy*float32(y),
						X:        vX*sx - 0.5 + sx*float32(x),
						Z:        (vZ - 0.5) * (1.0 / 16.0),
						Texture:  tex,
						TextureX: float64(vert.TOffsetX) / float64(16*rect.Width),
						TextureY: float64(vert.TOffsetY) / float64(16*rect.Height),
						R:        cr,
						G:        cg,
						B:        cb,
						A:        byte(aa >> 8),
					}
					out = append(out, vert)
				}
			}
		}
	}
	return
}
