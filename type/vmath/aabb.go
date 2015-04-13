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

package vmath

import (
	"fmt"
)

type AABB struct {
	Min Vector3
	Max Vector3
}

func NewAABB(x1, y1, z1, x2, y2, z2 float64) *AABB {
	return &AABB{
		Min: Vector3{x1, y1, z1},
		Max: Vector3{x2, y2, z2},
	}
}

func (a *AABB) RotateX(an, ox, oy, oz float64) {
	a.Max.RotateX(an, ox, oy, oz)
	a.Min.RotateX(an, ox, oy, oz)
	a.fixBounds()
}

func (a *AABB) RotateY(an, ox, oy, oz float64) {
	a.Max.RotateY(an, ox, oy, oz)
	a.Min.RotateY(an, ox, oy, oz)
	a.fixBounds()
}

func (a *AABB) fixBounds() {
	if a.Max.X < a.Min.X || a.Min.X > a.Max.X {
		a.Max.X, a.Min.X = a.Min.X, a.Max.X
	}
	if a.Max.Y < a.Min.Y || a.Min.Y > a.Max.Y {
		a.Max.Y, a.Min.Y = a.Min.Y, a.Max.Y
	}
	if a.Max.Z < a.Min.Z || a.Min.Z > a.Max.Z {
		a.Max.Z, a.Min.Z = a.Min.Z, a.Max.Z
	}
}

func (a *AABB) Intersects(o *AABB) bool {
	return !(o.Min.X >= a.Max.X ||
		o.Max.X <= a.Min.X ||
		o.Min.Y >= a.Max.Y ||
		o.Max.Y <= a.Min.Y ||
		o.Min.Z >= a.Max.Z ||
		o.Max.Z <= a.Min.Z)
}

func (a *AABB) IntersectsLine(origin, dir Vector3) bool {
	const right, left, middle = 0, 1, 2
	var (
		quadrant       [3]int
		candidatePlane [3]float64
		maxT           = [3]float64{-1, -1, -1}
	)
	inside := true
	findC := func(i int, x, minX, maxX float64) {
		if x < minX {
			quadrant[i] = left
			candidatePlane[i] = minX
			inside = false
		} else if x > maxX {
			quadrant[i] = right
			candidatePlane[i] = maxX
			inside = false
		} else {
			quadrant[i] = middle
		}
	}
	findC(0, origin.X, a.Min.X, a.Max.X)
	findC(1, origin.Y, a.Min.Y, a.Max.Y)
	findC(2, origin.Z, a.Min.Z, a.Max.Z)
	if inside {
		return true
	}

	if quadrant[0] != middle && dir.X != 0 {
		maxT[0] = (candidatePlane[0] - origin.X) / dir.X
	}
	if quadrant[1] != middle && dir.Y != 0 {
		maxT[1] = (candidatePlane[1] - origin.Y) / dir.Y
	}
	if quadrant[2] != middle && dir.Z != 0 {
		maxT[2] = (candidatePlane[2] - origin.Z) / dir.Z
	}
	whichPlane := 0
	for i := 1; i < 3; i++ {
		if maxT[whichPlane] < maxT[i] {
			whichPlane = i
		}
	}
	if maxT[whichPlane] < 0 {
		return false
	}
	check := func(i int, oX, dX, min, max float64) bool {
		if whichPlane != i {
			coord := oX + maxT[whichPlane]*dX
			if coord < min || coord > max {
				return false
			}
		}
		return true
	}
	if !check(0, origin.X, dir.X, a.Min.X, a.Max.X) {
		return false
	}
	if !check(1, origin.Y, dir.Y, a.Min.Y, a.Max.Y) {
		return false
	}
	if !check(2, origin.Z, dir.Z, a.Min.Z, a.Max.Z) {
		return false
	}
	return true
}

func (a *AABB) Shift(x, y, z float64) {
	a.Min.X += x
	a.Max.X += x
	a.Min.Y += y
	a.Max.Y += y
	a.Min.Z += z
	a.Max.Z += z
}

func (a *AABB) MoveOutOf(o *AABB, dir *Vector3) {
	if dir.X != 0 {
		if dir.X > 0 {
			ox := a.Max.X
			a.Max.X = o.Min.X - 0.0001
			a.Min.X += a.Max.X - ox
		} else {
			ox := a.Min.X
			a.Min.X = o.Max.X + 0.0001
			a.Max.X += a.Min.X - ox
		}
	}
	if dir.Y != 0 {
		if dir.Y > 0 {
			oy := a.Max.Y
			a.Max.Y = o.Min.Y - 0.0001
			a.Min.Y += a.Max.Y - oy
		} else {
			oy := a.Min.Y
			a.Min.Y = o.Max.Y + 0.0001
			a.Max.Y += a.Min.Y - oy
		}
	}

	if dir.Z != 0 {
		if dir.Z > 0 {
			oz := a.Max.Z
			a.Max.Z = o.Min.Z - 0.0001
			a.Min.Z += a.Max.Z - oz
		} else {
			oz := a.Min.Z
			a.Min.Z = o.Max.Z + 0.0001
			a.Max.Z += a.Min.Z - oz
		}
	}
}

func (a AABB) String() string {
	return fmt.Sprintf("[%s->%s]", a.Min, a.Max)
}
