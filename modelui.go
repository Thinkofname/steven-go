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
	"github.com/thinkofdeath/steven/type/direction"
	"github.com/thinkofdeath/steven/ui"
	"github.com/thinkofdeath/steven/world/biome"
)

func modelToUI(mdl *model, block Block) *ui.Model {
	mat := mgl32.Rotate3DX(math.Pi / 6).Mat4().
		Mul4(mgl32.Rotate3DY(math.Pi/4 + math.Pi).Mat4()).
		Mul4(mgl32.Scale3D(0.65, 0.65, 0.65))

	if gui, ok := mdl.display["gui"]; ok {
		if gui.Scale != nil {
			mat = mat.Mul4(mgl32.Scale3D(
				float32(gui.Scale[0]),
				float32(gui.Scale[1]),
				float32(gui.Scale[2]),
			))
		}
		if gui.Translation != nil {
			mat = mat.Mul4(mgl32.Translate3D(
				float32(gui.Translation[0]/16),
				float32(gui.Translation[1]/16),
				float32(gui.Translation[2]/16),
			))
		}
		if gui.Rotation != nil {
			mat = mat.Mul4(mgl32.Rotate3DY(float32(gui.Rotation[1]/180) * math.Pi).Mat4())
			mat = mat.Mul4(mgl32.Rotate3DX(float32(gui.Rotation[0]/180) * math.Pi).Mat4())
			mat = mat.Mul4(mgl32.Rotate3DZ(float32(gui.Rotation[2]/180) * math.Pi).Mat4())
		}
	}

	var verts []*ui.ModelVertex

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

		for _, vert := range f.vertices {
			vert := &ui.ModelVertex{
				X:        float32(vert.X) / 256,
				Y:        float32(vert.Y) / 256,
				Z:        float32(vert.Z) / 256,
				TOffsetX: vert.TOffsetX,
				TOffsetY: vert.TOffsetY,
				R:        cr,
				G:        cg,
				B:        cb,
				A:        255,
				TX:       vert.TX,
				TY:       vert.TY,
				TW:       vert.TW,
				TH:       vert.TH,
				TAtlas:   vert.TAtlas,
			}
			verts = append(verts, vert)
		}
	}
	return ui.NewModel(0, 0, verts, mat)
}
