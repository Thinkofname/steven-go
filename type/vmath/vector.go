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
	"math"
)

// Vector3 is a 3 component vector
type Vector3 struct {
	X, Y, Z float32
}

// Dot returns the result of preforming the dot operation on this
// vector and the passed vector.
func (v Vector3) Dot(other Vector3) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v *Vector3) Normalize() {
	l := float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
	v.X /= l
	v.Y /= l
	v.Z /= l
}

func (v *Vector3) LengthSquared() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v *Vector3) Apply(m *Matrix4) float32 {
	x := v.X
	y := v.Y
	z := v.Z
	w := m[3]*x + m[7]*y + m[11]*z + m[15]
	if w == 0 {
		w = 1
	}

	v.X = (m[0]*x + m[4]*y + m[8]*z + m[12]) / w
	v.Y = (m[1]*x + m[5]*y + m[9]*z + m[13]) / w
	v.Z = (m[2]*x + m[6]*y + m[10]*z + m[14]) / w
	return w
}
