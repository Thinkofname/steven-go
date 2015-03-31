package direction

import (
	"fmt"

	"github.com/thinkofdeath/steven/type/vmath"
)

// Type is a direction in the minecraft world
type Type uint

const (
	Up Type = iota
	Down
	North
	South
	West
	East
	Invalid
)

var Values = []Type{
	Up,
	Down,
	North,
	South,
	West,
	East,
}

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

func (d Type) AsVector() vmath.Vector3 {
	x, y, z := d.Offset()
	return vmath.Vector3{float32(x), float32(y), float32(z)}
}

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
