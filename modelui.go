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
)

func modelToUI(mdl *model) *ui.Model {
	mat := mgl32.Rotate3DX(math.Pi / 6).Mat4().
		Mul4(mgl32.Rotate3DY(math.Pi / 4).Mat4())
	var verts []*ui.ModelVertex
	for i := range mdl.elements {
		e := mdl.elements[len(mdl.elements)-1-i]
		for fi, f := range e.faces {
			if f == nil {
				continue
			}

			face := faceVertices[fi]

			x, w := e.from[0]/16, (e.to[0]-e.from[0])/16
			y, h := e.from[1]/16, (e.to[1]-e.from[1])/16
			z, d := e.from[2]/16, (e.to[2]-e.from[2])/16
			tex := mdl.lookupTexture(f.texture)
			for _, v := range face.verts {
				var rr, gg, bb byte = 255, 255, 255
				if direction.Type(fi) == direction.West || direction.Type(fi) == direction.East {
					rr = byte(255 * 0.8)
					gg = byte(255 * 0.8)
					bb = byte(255 * 0.8)
				}

				vert := &ui.ModelVertex{
					X:        float32(float64(v.X)*w + x),
					Y:        float32(float64(v.Y)*h + y),
					Z:        float32(float64(v.Z)*d + z),
					TOffsetX: int16(float64(v.TOffsetX*16*int16(tex.Width)) * 1),
					TOffsetY: int16(float64(v.TOffsetY*16*int16(tex.Height)) * 1),
					R:        rr,
					G:        gg,
					B:        bb,
					A:        255,
					TX:       uint16(tex.X),
					TY:       uint16(tex.Y),
					TW:       uint16(tex.Width),
					TH:       uint16(tex.Height),
					TAtlas:   int16(tex.Atlas),
				}
				verts = append(verts, vert)
			}
		}
	}
	return ui.NewModel(0, 0, verts, mat)
}
