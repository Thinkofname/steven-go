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
	"github.com/thinkofdeath/steven/entitysys"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/type/direction"
)

func (ce *clientEntities) registerModels() {
	ce.container.AddSystem(entitysys.Add, esPlayerModelAdd)
	ce.container.AddSystem(entitysys.Tick, esPlayerModelTick)

	// Generic removal
	ce.container.AddSystem(entitysys.Remove, esModelRemove)
}

func appendBox(verts []*render.StaticVertex, x, y, z, w, h, d float32, textures [6]*render.TextureInfo) []*render.StaticVertex {
	for i, face := range faceVertices {
		tex := textures[i]
		if tex == nil {
			continue
		}
		for _, v := range face.verts {
			vert := &render.StaticVertex{
				X:        float32(v.X)*w + x,
				Y:        float32(v.Y)*h + y,
				Z:        float32(v.Z)*d + z,
				TOffsetX: v.TOffsetX * 16 * int16(tex.Width),
				TOffsetY: v.TOffsetY * 16 * int16(tex.Height),
				R:        255,
				G:        255,
				B:        255,
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
	return verts
}

// Player

type playerModelComponent struct {
	model *render.StaticModel
}

func (p *playerModelComponent) SetModel(m *render.StaticModel) { p.model = m }
func (p *playerModelComponent) Model() *render.StaticModel     { return p.model }

// Marker method
func (*playerModelComponent) playerModel() {}

type PlayerModelComponent interface {
	SetModel(m *render.StaticModel)
	Model() *render.StaticModel

	playerModel()
}

const (
	playerModelHead = iota
	playerModelBody
)

func esPlayerModelAdd(p PlayerModelComponent, pl PlayerComponent) {

	uuid := pl.UUID()
	info := Client.playerList.info[uuid]
	if info == nil {
		panic("missing player info")
	}
	skin := info.skin

	var hverts []*render.StaticVertex
	hverts = appendBox(hverts, -4/16.0, 0, -4/16.0, 8/16.0, 8/16.0, 8/16.0, [6]*render.TextureInfo{
		direction.North: skin.Sub(8, 8, 8, 8),
		direction.South: skin.Sub(24, 8, 8, 8),
		direction.East:  skin.Sub(0, 8, 8, 8),
		direction.West:  skin.Sub(16, 8, 8, 8),
		direction.Up:    skin.Sub(8, 0, 8, 8),
		direction.Down:  skin.Sub(16, 0, 8, 8),
	})
	hverts = appendBox(hverts, -4.5/16.0, -.5/16.0, -4.5/16.0, 9/16.0, 9/16.0, 9/16.0, [6]*render.TextureInfo{
		direction.North: skin.Sub(8+32, 8, 8, 8),
		direction.South: skin.Sub(24+32, 8, 8, 8),
		direction.East:  skin.Sub(0+32, 8, 8, 8),
		direction.West:  skin.Sub(16+32, 8, 8, 8),
		direction.Up:    skin.Sub(8+32, 0, 8, 8),
		direction.Down:  skin.Sub(16+32, 0, 8, 8),
	})

	var bverts []*render.StaticVertex
	bverts = appendBox(bverts, -4/16.0, -6/16.0, -2/16.0, 8/16.0, 12/16.0, 4/16.0, [6]*render.TextureInfo{
		direction.North: skin.Sub(20, 20, 8, 12),
		direction.South: skin.Sub(32, 20, 8, 12),
		direction.East:  skin.Sub(16, 20, 4, 12),
		direction.West:  skin.Sub(28, 20, 4, 12),
		direction.Up:    skin.Sub(20, 16, 8, 4),
		direction.Down:  skin.Sub(28, 16, 8, 4),
	})
	bverts = appendBox(bverts, -4.5/16.0, -6.5/16.0, -2.5/16.0, 9/16.0, 13/16.0, 5/16.0, [6]*render.TextureInfo{
		direction.North: skin.Sub(20, 20+16, 8, 12),
		direction.South: skin.Sub(32, 20+16, 8, 12),
		direction.East:  skin.Sub(16, 20+16, 4, 12),
		direction.West:  skin.Sub(28, 20+16, 4, 12),
		direction.Up:    skin.Sub(20, 16+16, 8, 4),
		direction.Down:  skin.Sub(28, 16+16, 8, 4),
	})

	model := render.NewStaticModel([][]*render.StaticVertex{
		playerModelHead: hverts,
		playerModelBody: bverts,
	})
	p.SetModel(model)
}

func esModelRemove(p interface {
	Model() *render.StaticModel
}) {
	p.Model().Free()
}

func esPlayerModelTick(p PlayerModelComponent, pos PositionComponent, r RotationComponent) {
	x, y, z := pos.Position()
	offMat := mgl32.Translate3D(float32(x), -float32(y), float32(z)).
		Mul4(mgl32.Rotate3DY(math.Pi - float32(r.Yaw())).Mat4())

	model := p.Model()
	model.Matrix[playerModelHead] = offMat.Mul4(mgl32.Translate3D(0, -1.62+(4/16.0), 0)).
		Mul4(mgl32.Rotate3DX(float32(r.Pitch())).Mat4())
	model.Matrix[playerModelBody] = offMat.Mul4(mgl32.Translate3D(0, -1.62+(10/16.0), 0))
}
