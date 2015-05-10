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

import "github.com/thinkofdeath/steven/render"

func (ce *clientEntities) register() {
	ce.container.AddSystem(es_moveToTarget)
	ce.container.AddSystem(es_drawOutline)
}

func es_drawOutline(p PositionComponent, s SizeComponent, d DebugComponent) {
	bounds := s.Bounds()
	x, y, z := p.Position()
	bounds.Shift(float32(x), float32(y), float32(z))

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

func es_moveToTarget(p PositionComponent, t TargetPositionComponent) {
	px, py, pz := p.Position()
	tx, ty, tz := t.TargetPosition()

	dx, dy, dz := tx-px, ty-py, tz-pz

	px += dx * 0.4 * Client.delta
	py += dy * 0.4 * Client.delta
	pz += dz * 0.4 * Client.delta
	p.SetPosition(px, py, pz)
}
