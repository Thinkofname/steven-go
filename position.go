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
	"github.com/go-gl/mathgl/mgl32"
	"github.com/thinkofdeath/steven/type/direction"
)

type Position struct {
	X, Y, Z int
}

func (p Position) Shift(x, y, z int) Position {
	return Position{X: p.X + x, Y: p.Y + y, Z: p.Z + z}
}

func (p Position) ShiftDir(d direction.Type) Position {
	return p.Shift(d.Offset())
}

func (p Position) Get() (int, int, int) {
	return p.X, p.Y, p.Z
}

func (p Position) Vec() mgl32.Vec3 {
	return mgl32.Vec3{float32(p.X), float32(p.Y), float32(p.Z)}
}
