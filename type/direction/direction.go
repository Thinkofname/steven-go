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

package direction

import (
	"fmt"

	"github.com/thinkofdeath/steven/type/vmath"
)

// Type is a direction in the minecraft world.
type Type uint

// Possible direction values.
const (
	Up Type = iota
	Down
	North
	South
	West
	East
	Invalid
)

// Values is all valid directions.
var Values = []Type{
	Up,
	Down,
	North,
	South,
	West,
	East,
}

// FromString returns the direction that matches the passed
// string if possible otherwise it will return Invalid.
func FromString(str string) Type {
	switch str {
	case "up":
		return Up
	case "down":
		return Down
	case "north":
		return North
	case "south":
		return South
	case "west":
		return West
	case "east":
		return East
	}
	// ¯\_(ツ)_/¯
	return Invalid
}

// Offset returns the x, y and z offset this direction points in.
func (d Type) Offset() (x, y, z int) {
	switch d {
	case Up:
		return 0, 1, 0
	case Down:
		return 0, -1, 0
	case North:
		return 0, 0, -1
	case South:
		return 0, 0, 1
	case West:
		return -1, 0, 0
	case East:
		return 1, 0, 0
	}
	return 0, 0, 0

}

// AsVector returns a vector of the direction's offset.
func (d Type) AsVector() vmath.Vector3 {
	x, y, z := d.Offset()
	return vmath.Vector3{float32(x), float32(y), float32(z)}
}

// Opposite returns the direction directly opposite to this
// direction.
func (d Type) Opposite() Type {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case East:
		return West
	case West:
		return East
	case North:
		return South
	case South:
		return North
	}
	return Invalid
}

// String returns a string representation of the direction.
func (d Type) String() string {
	switch d {
	case Up:
		return "up"
	case Down:
		return "down"
	case North:
		return "north"
	case South:
		return "south"
	case West:
		return "west"
	case East:
		return "east"
	case Invalid:
		return "invalid"
	}
	return fmt.Sprintf("direction.Type(%d)", d)
}
