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

func (a *AABB) Intersects(o *AABB) bool {
	return !(o.Min.X > a.Max.X ||
		o.Max.X < a.Min.X ||
		o.Min.Y > a.Max.Y ||
		o.Max.Y < a.Min.Y ||
		o.Min.Z > a.Max.Z ||
		o.Max.Z < a.Min.Z)
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
