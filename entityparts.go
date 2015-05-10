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

import "github.com/thinkofdeath/steven/type/vmath"

// Networkable

type networkComponent struct {
	networkID int
}

func (n *networkComponent) NetworkID() int { return n.networkID }

type NetworkComponent interface {
	NetworkID() int
}

// Positionable

type positionComponent struct {
	X, Y, Z float64
}

func (p *positionComponent) Position() (x, y, z float64) {
	return p.X, p.Y, p.Z
}

func (p *positionComponent) SetPosition(x, y, z float64) {
	p.X, p.Y, p.Z = x, y, z
}

type PositionComponent interface {
	Position() (x, y, z float64)
	SetPosition(x, y, z float64)
}

// Target Positionable

type targetPositionComponent struct {
	X, Y, Z float64
}

func (p *targetPositionComponent) TargetPosition() (x, y, z float64) {
	return p.X, p.Y, p.Z
}

func (p *targetPositionComponent) SetTargetPosition(x, y, z float64) {
	p.X, p.Y, p.Z = x, y, z
}

type TargetPositionComponent interface {
	TargetPosition() (x, y, z float64)
	SetTargetPosition(x, y, z float64)
}

// Sized

type sizeComponent struct {
	bounds vmath.AABB
}

func (s sizeComponent) Bounds() vmath.AABB { return s.bounds }

type SizeComponent interface {
	Bounds() vmath.AABB
}

// Debug

type debugComponent struct {
	R, G, B byte
}

func (d debugComponent) DebugColor() (r, g, b byte) {
	return d.R, d.G, d.B
}

type DebugComponent interface {
	DebugColor() (r, g, b byte)
}
