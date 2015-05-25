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
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/thinkofdeath/steven/render"
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
			vX, vY, vZ := 1.0-(float32(vert.X)/256), float32(vert.Y)/256, float32(vert.Z)/256
			/*vec := mgl32.Vec3{vX - 0.5, vY - 0.5, vZ - 0.5}
			vec = mat.Mul4x1(vec.Vec4(1)).Vec3().
				Add(mgl32.Vec3{0.5, 0.5, 0.5})
			vX, vY, vZ = vec[0], 1.0-vec[1], vec[2]*/
			tex := f.verticesTexture[i]
			rect := tex.Rect()
			vert := &render.StaticVertex{
				X:        0.5 - vX,
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
