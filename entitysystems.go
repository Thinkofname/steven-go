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

	"github.com/thinkofdeath/steven/entitysys"
	"github.com/thinkofdeath/steven/render"
)

func init() {
	addSystem(entitysys.Tick, esMoveToTarget)
	addSystem(entitysys.Tick, esRotateToTarget)
	addSystem(entitysys.Tick, esDrawOutline)
	addSystem(entitysys.Tick, esLightModel)
	addSystem(entitysys.Tick, esMoveChunk)
}

func esDrawOutline(p PositionComponent, s SizeComponent, d DebugComponent) {
	x, y, z := p.Position()
	bounds := s.Bounds().Shift(float32(x), float32(y), float32(z))

	r, g, b := d.DebugColor()
	render.DrawBox(
		float64(bounds.Min.X()),
		float64(bounds.Min.Y()),
		float64(bounds.Min.Z()),
		float64(bounds.Max.X()),
		float64(bounds.Max.Y()),
		float64(bounds.Max.Z()),
		r, g, b, 255,
	)
}

// updates the Colors of the model to fake lighting
func esLightModel(p PositionComponent, m interface {
	Model() *render.StaticModel
}) {
	if m.Model() == nil {
		return
	}
	x, y, z := p.Position()
	bx, by, bz := int(math.Floor(x)), int(math.Floor(y)), int(math.Floor(z))
	bl := float64(chunkMap.BlockLight(bx, by, bz)) / 16
	sl := float64(chunkMap.SkyLight(bx, by, bz)) / 16
	light := math.Max(bl, sl) + (1 / 16.0)
	model := m.Model()
	for i := range model.Colors {
		model.Colors[i] = [4]float32{
			float32(light),
			float32(light),
			float32(light),
			1.0,
		}
	}
}

// Moves the entity from the previous chunk to its
// new chunk. Allows for optimized lookups
func esMoveChunk(e Entity, p *positionComponent) {
	cx, cz := int(p.X)>>4, int(p.Z)>>4
	if cx != p.CX || cz != p.CZ {
		oc := chunkMap[chunkPosition{p.CX, p.CZ}]
		if oc != nil {
			oc.removeEntity(e)
		}
		c := chunkMap[chunkPosition{cx, cz}]
		if c != nil {
			c.addEntity(e)
			p.CX, p.CZ = cx, cz
		}
	}
}

// Smoothly moves the entity from its current position to the target
// location
func esMoveToTarget(p PositionComponent, t *targetPositionComponent) {
	px, py, pz := p.Position()
	tx, ty, tz := t.TargetPosition()

	if t.pX != tx || t.pY != ty || t.pZ != tz || t.time >= 4 {
		t.sX, t.sY, t.sZ = px, py, pz
		t.time = 0
		t.pX = tx
		t.pY = ty
		t.pZ = tz
	}
	sx, sy, sz := t.sX, t.sY, t.sZ

	dx, dy, dz := tx-sx, ty-sy, tz-sz

	t.time = math.Min(4.0, t.time+Client.delta)

	px = sx + dx*(1/4.0)*t.time
	py = sy + dy*(1/4.0)*t.time
	pz = sz + dz*(1/4.0)*t.time
	p.SetPosition(px, py, pz)
}

// Smoothly rotates the entity from its current rotation to the target
// rotation
func esRotateToTarget(r RotationComponent, t *targetRotationComponent) {
	py, pp := r.Yaw(), r.Pitch()
	ty, tp := t.TargetYaw(), t.TargetPitch()

	if t.pPitch != tp || t.pYaw != ty || t.time >= 4 {
		t.sYaw, t.sPitch = py, pp
		t.time = 0
		t.pPitch = tp
		t.pYaw = ty
	}
	sy, sp := t.sYaw, t.sPitch

	dy, dp := ty-sy, tp-sp
	// Make sure we go for the shortest route.
	// e.g. (in degrees) 1 to 359 is quicker
	// to decrease to wrap around than it is
	// to increase all the way around
	if dy > math.Pi || dy < -math.Pi {
		sy += math.Copysign(math.Pi*2, dy)
		dy = ty - sy
	}
	if dp > math.Pi || dp < -math.Pi {
		sp += math.Copysign(math.Pi*2, dp)
		dp = tp - sp
	}

	t.time = math.Min(4.0, t.time+Client.delta)

	py = sy + dy*(1/4.0)*t.time
	pp = sp + dp*(1/4.0)*t.time
	r.SetPitch(pp)
	r.SetYaw(py)
}
