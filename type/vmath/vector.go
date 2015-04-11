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
	"math"
)

// Vector3 is a 3 component vector
type Vector3 struct {
	X, Y, Z float64
}

// Dot returns the result of preforming the dot operation on this
// vector and the passed vector.
func (v Vector3) Dot(other Vector3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v *Vector3) Normalize() {
	l := v.Length()
	v.X /= l
	v.Y /= l
	v.Z /= l
}

func (v Vector3) DistanceSquared(o Vector3) float64 {
	dx := v.X - o.X
	dy := v.Y - o.Y
	dz := v.Z - o.Z
	return dx*dx + dy*dy + dz*dz
}

func (v Vector3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vector3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v *Vector3) Apply(m *Matrix4) float64 {
	x := v.X
	y := v.Y
	z := v.Z
	w := float64(m[3])*x + float64(m[7])*y + float64(m[11])*z + float64(m[15])
	if w == 0 {
		w = 1
	}

	v.X = (float64(m[0])*x + float64(m[4])*y + float64(m[8])*z + float64(m[12])) / w
	v.Y = (float64(m[1])*x + float64(m[5])*y + float64(m[9])*z + float64(m[13])) / w
	v.Z = (float64(m[2])*x + float64(m[6])*y + float64(m[10])*z + float64(m[14])) / w
	return w
}

func (v Vector3) AngleTo(other Vector3) float64 {
	return v.Dot(other) / (v.Length() + other.Length())
}

func (v Vector3) String() string {
	return fmt.Sprintf("(%f,%f,%f)", v.X, v.Y, v.Z)
}
