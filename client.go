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

package main

import (
	"fmt"
	"math"

	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/type/vmath"
)

const (
	playerHeight = 1.62
)

var Client = ClientState{
	Bounds: vmath.AABB{
		Min: vmath.Vector3{-0.3, 0, -0.3},
		Max: vmath.Vector3{0.3, 1.8, 0.3},
	},
}

type ClientState struct {
	X, Y, Z    float64
	Yaw, Pitch float64

	Jumping  bool
	VSpeed   float64
	OnGround bool

	GameMode gameMode
	HardCore bool

	Bounds vmath.AABB

	positionText *render.UIText

	chat ChatUI
}

// The render tick needs to remain pretty light so it
// doesn't hold the lock for too long.
func (c *ClientState) renderTick(delta float64) {
	if c.GameMode.Fly() {
		c.X += mf * math.Cos(c.Yaw-math.Pi/2) * -math.Cos(c.Pitch) * delta * 0.2
		c.Z -= mf * math.Sin(c.Yaw-math.Pi/2) * -math.Cos(c.Pitch) * delta * 0.2
		c.Y -= mf * math.Sin(c.Pitch) * delta * 0.2
	} else {
		c.X += mf * math.Cos(c.Yaw-math.Pi/2) * delta * 0.1
		c.Z -= mf * math.Sin(c.Yaw-math.Pi/2) * delta * 0.1
		if !c.OnGround {
			c.VSpeed -= 0.01 * delta
			if c.VSpeed < -0.3 {
				c.VSpeed = -0.3
			}
		} else if c.Jumping {
			c.VSpeed = 0.15
		} else {
			c.VSpeed = 0
		}
		c.Y += c.VSpeed * delta
	}

	if !c.GameMode.NoClip() && chunkMap[chunkPosition{int(c.X) >> 4, int(c.Z) >> 4}] != nil {
		cy := c.Y
		cz := c.Z
		c.Y = render.Camera.Y - playerHeight
		c.Z = render.Camera.Z

		// We handle each axis separately to allow for a sliding
		// effect when pushing up against walls.

		bounds, _ := c.checkCollisions(c.Bounds)
		c.X = bounds.Min.X + 0.3

		c.Z = cz
		bounds, _ = c.checkCollisions(c.Bounds)
		c.Z = bounds.Min.Z + 0.3

		c.Y = cy
		bounds, _ = c.checkCollisions(c.Bounds)
		c.Y = bounds.Min.Y

		ground := vmath.AABB{
			Min: vmath.Vector3{-0.3, -0.05, -0.3},
			Max: vmath.Vector3{0.3, 0, 0.3},
		}
		_, c.OnGround = c.checkCollisions(ground)
	}

	// Copy to the camera
	render.Camera.X = c.X
	render.Camera.Y = c.Y + playerHeight
	render.Camera.Z = c.Z
	render.Camera.Yaw = c.Yaw
	render.Camera.Pitch = c.Pitch

	if c.positionText != nil {
		c.positionText.Free()
	}
	c.positionText = render.AddUIText(
		fmt.Sprintf("X: %.2f, Y: %.2f, Z: %.2f", c.X, c.Y, c.Z),
		5, 5, 255, 255, 255,
	)

	c.chat.render(delta)
}

func (c *ClientState) checkCollisions(bounds vmath.AABB) (vmath.AABB, bool) {
	bounds.Shift(c.X, c.Y, c.Z)

	dir := &vmath.Vector3{
		X: -(render.Camera.X - c.X),
		Y: -(render.Camera.Y - playerHeight - c.Y),
		Z: -(render.Camera.Z - c.Z),
	}

	minX, minY, minZ := int(bounds.Min.X-1), int(bounds.Min.Y-1), int(bounds.Min.Z-1)
	maxX, maxY, maxZ := int(bounds.Max.X+1), int(bounds.Max.Y+1), int(bounds.Max.Z+1)

	hit := false
	for y := minY; y < maxY; y++ {
		for z := minZ; z < maxZ; z++ {
			for x := minX; x < maxX; x++ {
				b := chunkMap.Block(x, y, z)

				if b.Collidable() {
					for _, bb := range b.CollisionBounds() {
						bb.Shift(float64(x), float64(y), float64(z))
						if bb.Intersects(&bounds) {
							bounds.MoveOutOf(&bb, dir)
							hit = true
						}
					}
				}
			}
		}
	}
	return bounds, hit
}

func (c *ClientState) tick() {
}

type gameMode int

const (
	gmSurvival gameMode = iota
	gmCreative
	gmAdventure
	gmSpecator
)

func (g gameMode) Fly() bool {
	switch g {
	case gmCreative, gmSpecator:
		return true
	}
	return false
}

func (g gameMode) NoClip() bool {
	switch g {
	case gmSpecator:
		return true
	}
	return false
}

type teleportFlag byte

const (
	teleportRelX teleportFlag = 1 << iota
	teleportRelY
	teleportRelZ
	teleportRelYaw
	teleportRelPitch
)

func calculateTeleport(flag teleportFlag, flags byte, base, val float64) float64 {
	if flags&byte(flag) != 0 {
		return base + val
	}
	return val
}
