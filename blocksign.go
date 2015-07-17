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
	"github.com/go-gl/mathgl/mgl32"
	"github.com/thinkofdeath/steven/format"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/type/direction"
	"github.com/thinkofdeath/steven/ui"
)

type signComponent struct {
	model *render.StaticModel

	lines    [4]format.AnyComponent
	position Position

	ox, oy, oz float64
	rotation   float64
	hasStand   bool
}

type SignComponent interface {
	Update(lines [4]format.AnyComponent)
}

func (s *signComponent) Model() *render.StaticModel {
	return s.model
}

func (s *signComponent) Update(lines [4]format.AnyComponent) {
	s.free()
	s.lines = lines
	s.create()
}

func (s *signComponent) free() {
	if s.model != nil {
		s.model.Free()
	}
}

func (s *signComponent) create() {
	const yS = (6.0 / 16.0) / 4.0
	const xS = yS / 16.0

	var verts []*render.StaticVertex
	for i, line := range s.lines {
		if line.Value == nil {
			continue
		}
		format.ConvertLegacy(line)
		// Hijack ui.Formatted's component parsing to split
		// up components into ui.Text elements.
		// TODO(Think) Move this into some common place for
		// easier reuse in other places?
		wrap := &format.TextComponent{}
		wrap.Color = format.Black
		wrap.Extra = []format.AnyComponent{line}
		f := ui.NewFormatted(format.Wrap(wrap), 0, 0)
		offset := 0.0
		for _, txt := range f.Text {
			str := txt.Value()

			for _, r := range str {
				tex := render.CharacterTexture(r)
				if tex == nil {
					continue
				}
				s := render.SizeOfCharacter(r)
				if r == ' ' {
					offset += (s + 2) * xS
					continue
				}

				for _, v := range faceVertices[direction.North].verts {
					vert := &render.StaticVertex{
						X:        float32(v.X)*float32(s*xS) - float32(offset+s*xS) + float32(f.Width*xS*0.5),
						Y:        float32(v.Y)*yS - yS*float32(i-1),
						Z:        -.6 / 16.0,
						Texture:  tex,
						TextureX: float64(v.TOffsetX),
						TextureY: float64(v.TOffsetY),
						R:        byte(txt.R()),
						G:        byte(txt.G()),
						B:        byte(txt.B()),
						A:        255,
					}
					verts = append(verts, vert)
				}
				offset += (s + 2) * xS
			}
		}
	}
	wood := render.GetTexture("blocks/planks_oak")
	// The backboard
	verts = appendBoxExtra(verts, -0.5, -4/16.0, -0.5/16.0, 1.0, 8/16.0, 1/16.0, [6]render.TextureInfo{
		direction.Up:    wood.Sub(0, 0, 16, 2),
		direction.Down:  wood.Sub(0, 0, 16, 2),
		direction.East:  wood.Sub(0, 0, 2, 12),
		direction.West:  wood.Sub(0, 0, 2, 12),
		direction.North: wood.Sub(0, 4, 16, 12),
		direction.South: wood.Sub(0, 4, 16, 12),
	}, [6][2]float64{
		direction.Up:    {1.5, 1.0},
		direction.Down:  {1.5, 1.0},
		direction.East:  {1.0, 1.0},
		direction.West:  {1.0, 1.0},
		direction.North: {1.5, 1.0},
		direction.South: {1.5, 1.0},
	})
	if s.hasStand {
		// Stand
		log := render.GetTexture("blocks/log_oak")
		verts = appendBox(verts, -0.5/16.0, -0.25-9/16.0, -0.5/16.0, 1/16.0, 9/16.0, 1/16.0, [6]render.TextureInfo{
			direction.Up:    log.Sub(0, 0, 2, 2),
			direction.Down:  log.Sub(0, 0, 2, 2),
			direction.East:  log.Sub(0, 0, 2, 12),
			direction.West:  log.Sub(0, 0, 2, 12),
			direction.North: log.Sub(0, 0, 2, 12),
			direction.South: log.Sub(0, 0, 2, 12),
		})
	}
	s.model = render.NewStaticModel([][]*render.StaticVertex{
		verts,
	})
	s.model.Radius = 2
	x, y, z := s.position.X, s.position.Y, s.position.Z

	s.model.X, s.model.Y, s.model.Z = -float32(x)-0.5, -float32(y)-0.5, float32(z)+0.5
	s.model.Matrix[0] = mgl32.Translate3D(
		float32(x)+0.5,
		-float32(y)-0.5,
		float32(z)+0.5,
	).Mul4(mgl32.Rotate3DY(float32(s.rotation)).Mat4()).
		Mul4(mgl32.Translate3D(float32(s.ox), float32(-s.oy), float32(s.oz)))
}

func esSignAdd(s *signComponent, b BlockComponent) {
	s.position = b.Position()
	s.create()
}

func esSignRemove(s *signComponent) {
	s.free()
}
