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

	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/type/vmath"
)

// Network

type networkComponent struct {
	networkID int
}

func (n *networkComponent) NetworkID() int { return n.networkID }

type NetworkComponent interface {
	NetworkID() int
}

// Position

type positionComponent struct {
	X, Y, Z    float64
	LX, LY, LZ float64
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

// Target Position

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

// Rotation

type rotationComponent struct {
	yaw, pitch float64
}

func (r *rotationComponent) Yaw() float64 { return r.yaw }
func (r *rotationComponent) SetYaw(y float64) {
	r.yaw = math.Mod(math.Pi*2+y, math.Pi*2)
}
func (r *rotationComponent) Pitch() float64 { return r.pitch }
func (r *rotationComponent) SetPitch(p float64) {
	r.pitch = math.Mod(math.Pi*2+p, math.Pi*2)
}

type RotationComponent interface {
	Yaw() float64
	SetYaw(y float64)
	Pitch() float64
	SetPitch(p float64)
}

// Target Rotation

type targetRotationComponent struct {
	yaw, pitch float64
}

func (r *targetRotationComponent) TargetYaw() float64 { return r.yaw }
func (r *targetRotationComponent) SetTargetYaw(y float64) {
	r.yaw = math.Mod(math.Pi*2+y, math.Pi*2)
}
func (r *targetRotationComponent) TargetPitch() float64 { return r.pitch }
func (r *targetRotationComponent) SetTargetPitch(p float64) {
	r.pitch = math.Mod(math.Pi*2+p, math.Pi*2)
}

type TargetRotationComponent interface {
	TargetYaw() float64
	SetTargetYaw(y float64)
	TargetPitch() float64
	SetTargetPitch(p float64)
}

// Size

type sizeComponent struct {
	bounds vmath.AABB
}

func (s sizeComponent) Bounds() vmath.AABB { return s.bounds }

type SizeComponent interface {
	Bounds() vmath.AABB
}

// Player

type playerComponent struct {
	uuid protocol.UUID
}

func (p *playerComponent) SetUUID(u protocol.UUID) {
	p.uuid = u
}
func (p *playerComponent) UUID() protocol.UUID {
	return p.uuid
}

type PlayerComponent interface {
	SetUUID(protocol.UUID)
	UUID() protocol.UUID
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
