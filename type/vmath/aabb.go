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

	"github.com/go-gl/mathgl/mgl32"
)

type AABB struct {
	Min mgl32.Vec3
	Max mgl32.Vec3
}

func NewAABB(x1, y1, z1, x2, y2, z2 float32) AABB {
	return AABB{
		Min: mgl32.Vec3{x1, y1, z1},
		Max: mgl32.Vec3{x2, y2, z2},
	}
}

func (a AABB) RotateX(an, ox, oy, oz float32) AABB {
	mat := mgl32.Rotate3DX(an)
	o := mgl32.Vec3{ox, oy, oz}
	a.Max = mat.Mul3x1(a.Max.Sub(o)).Add(o)
	a.Min = mat.Mul3x1(a.Min.Sub(o)).Add(o)
	a.fixBounds()
	return a
}

func (a AABB) RotateY(an, ox, oy, oz float32) AABB {
	mat := mgl32.Rotate3DY(an)
	o := mgl32.Vec3{ox, oy, oz}
	a.Max = mat.Mul3x1(a.Max.Sub(o)).Add(o)
	a.Min = mat.Mul3x1(a.Min.Sub(o)).Add(o)
	a.fixBounds()
	return a
}

func (a *AABB) fixBounds() {
	for i := range a.Min {
		if a.Max[i] < a.Min[i] || a.Min[i] > a.Max[i] {
			a.Max[i], a.Min[i] = a.Min[i], a.Max[i]
		}
	}
}

func (a AABB) Intersects(o AABB) bool {
	return !(o.Min.X() >= a.Max.X() ||
		o.Max.X() <= a.Min.X() ||
		o.Min.Y() >= a.Max.Y() ||
		o.Max.Y() <= a.Min.Y() ||
		o.Min.Z() >= a.Max.Z() ||
		o.Max.Z() <= a.Min.Z())
}

func (a AABB) IntersectsLine(origin, dir mgl32.Vec3) (mgl32.Vec3, bool) {
	const right, left, middle = 0, 1, 2
	var (
		quadrant       [3]int
		candidatePlane [3]float32
		maxT           = [3]float32{-1, -1, -1}
	)
	inside := true
	for i := range origin {
		if origin[i] < a.Min[i] {
			quadrant[i] = left
			candidatePlane[i] = a.Min[i]
			inside = false
		} else if origin[i] > a.Max[i] {
			quadrant[i] = right
			candidatePlane[i] = a.Max[i]
			inside = false
		} else {
			quadrant[i] = middle
		}
	}
	if inside {
		return origin, true
	}

	for i := range dir {
		if quadrant[i] != middle && dir[i] != 0 {
			maxT[i] = (candidatePlane[i] - origin[i]) / dir[i]
		}
	}
	whichPlane := 0
	for i := 1; i < 3; i++ {
		if maxT[whichPlane] < maxT[i] {
			whichPlane = i
		}
	}
	if maxT[whichPlane] < 0 {
		return origin, false
	}

	var coord mgl32.Vec3
	for i := range origin {
		if whichPlane != i {
			coord[i] = origin[i] + maxT[whichPlane]*dir[i]
			if coord[i] < a.Min[i] || coord[i] > a.Max[i] {
				return origin, false
			}
		} else {
			coord[i] = candidatePlane[i]
		}
	}
	return coord, true
}

func (a AABB) Shift(x, y, z float32) AABB {
	a.Min[0] += x
	a.Max[0] += x
	a.Min[1] += y
	a.Max[1] += y
	a.Min[2] += z
	a.Max[2] += z
	return a
}

func (a AABB) Grow(x, y, z float32) AABB {
	a.Min[0] -= x
	a.Max[0] += x
	a.Min[1] -= y
	a.Max[1] += y
	a.Min[2] -= z
	a.Max[2] += z
	return a
}

func (a AABB) MoveOutOf(o AABB, dir mgl32.Vec3) AABB {
	if dir.X() != 0 {
		if dir.X() > 0 {
			ox := a.Max.X()
			a.Max[0] = o.Min.X() - 0.0001
			a.Min[0] += a.Max.X() - ox
		} else {
			ox := a.Min.X()
			a.Min[0] = o.Max.X() + 0.0001
			a.Max[0] += a.Min.X() - ox
		}
	}
	if dir.Y() != 0 {
		if dir.Y() > 0 {
			oy := a.Max.Y()
			a.Max[1] = o.Min.Y() - 0.0001
			a.Min[1] += a.Max.Y() - oy
		} else {
			oy := a.Min.Y()
			a.Min[1] = o.Max.Y() + 0.0001
			a.Max[1] += a.Min.Y() - oy
		}
	}

	if dir.Z() != 0 {
		if dir.Z() > 0 {
			oz := a.Max.Z()
			a.Max[2] = o.Min.Z() - 0.0001
			a.Min[2] += a.Max.Z() - oz
		} else {
			oz := a.Min.Z()
			a.Min[2] = o.Max.Z() + 0.0001
			a.Max[2] += a.Min.Z() - oz
		}
	}
	return a
}

func (a AABB) String() string {
	return fmt.Sprintf("[%v->%v]", a.Min, a.Max)
}
