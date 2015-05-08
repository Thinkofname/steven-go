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

import "github.com/go-gl/mathgl/mgl32"

type Frustum struct {
	planes [6]fPlane
}

type fPlane struct {
	N Vector3
	D float64
}

func (f *Frustum) FromMatrix(m mgl32.Mat4) {
	for i := range f.planes {
		off := i >> 1
		f.planes[i] = fPlane{
			N: Vector3{
				X: float64(m.At(0, 3) - m.At(0, off)),
				Y: float64(m.At(1, 3) - m.At(1, off)),
				Z: float64(m.At(2, 3) - m.At(2, off)),
			},
			D: float64(m.At(3, 3) - m.At(3, off)),
		}
	}

	for i := range f.planes {
		f.planes[i].N.Normalize()
	}
}

func (f *Frustum) IsSphereInside(x, y, z, radius float64) bool {
	center := Vector3{x, y, z}
	for i := 0; i < 6; i++ {
		if center.Dot(f.planes[i].N)+f.planes[i].D+radius <= 0 {
			return false
		}
	}
	return true
}
