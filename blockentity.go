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
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/thinkofdeath/steven/encoding/nbt"
	"github.com/thinkofdeath/steven/entitysys"
	"github.com/thinkofdeath/steven/render"
)

func init() {
	addSystem(entitysys.Add, esSkullAdd)
	addSystem(entitysys.Remove, esSkullRemove)
	addSystem(entitysys.Add, esSignAdd)
}

// updates the Colors of the model to fake lighting
func lightBlockModel(model *render.StaticModel, bp Position) {
	bx, by, bz := bp.X, bp.Y, bp.Z
	bl := float64(chunkMap.BlockLight(bx, by, bz)) / 16
	sl := float64(chunkMap.SkyLight(bx, by, bz)) / 16
	light := math.Max(bl, sl) + (1 / 16.0)
	for i := range model.Colors {
		model.Colors[i] = [4]float32{
			float32(light),
			float32(light),
			float32(light),
			1.0,
		}
	}
}

// BlockEntity is the interface for which all block entities
// must implement
type BlockEntity interface {
	BlockComponent
}

type blockComponent struct {
	Location Position
}

func (bc *blockComponent) Position() Position {
	return bc.Location
}

func (bc *blockComponent) SetPosition(p Position) {
	bc.Location = p
}

// BlockComponent is a component that defines the location
// of an entity when attached to a block.
type BlockComponent interface {
	Position() Position
	SetPosition(p Position)
}

// BlockNBTComponent is implemented by block entities that
// load information from nbt.
type BlockNBTComponent interface {
	Deserilize(tag *nbt.Compound)
	CanHandleAction(action int) bool
}

type blockBreakComponent struct {
	blockComponent
	stage int
	model *render.StaticModel
}

func (b *blockBreakComponent) SetStage(stage int)         { b.stage = stage }
func (b *blockBreakComponent) Stage() int                 { return b.stage }
func (b *blockBreakComponent) Model() *render.StaticModel { return b.model }
func (b *blockBreakComponent) Update() {
	if b.model != nil {
		b.model.Free()
	}
	bounds := chunkMap.Block(b.Location.X, b.Location.Y, b.Location.Z).CollisionBounds()
	tex := render.GetTexture(fmt.Sprintf("blocks/destroy_stage_%d", b.stage))

	var verts []*render.StaticVertex
	for _, bo := range bounds {
		// Slightly bigger than the block to prevent clipping
		bo = bo.Grow(0.01, 0.01, 0.01)
		verts = appendBox(verts,
			bo.Min.X(), bo.Min.Y(), bo.Min.Z(),
			bo.Max.X()-bo.Min.X(), bo.Max.Y()-bo.Min.Y(), bo.Max.Z()-bo.Min.Z(),
			[6]render.TextureInfo{
				tex, tex, tex, tex, tex, tex,
			})
	}
	b.model = render.NewStaticModel([][]*render.StaticVertex{verts})

	b.model.Matrix[0] = mgl32.Translate3D(
		float32(b.Location.X),
		-float32(b.Location.Y),
		float32(b.Location.Z),
	)
}

// BlockBreakComponent is implemented by the block break animation
// entity
type BlockBreakComponent interface {
	SetStage(stage int)
	Stage() int
	Update()
}

func newBlockBreakEntity() BlockEntity {
	type blockBreak struct {
		networkComponent
		blockBreakComponent
	}
	b := &blockBreak{}
	return b
}
