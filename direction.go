package main

import (
	"fmt"
)

// Direction is a direction in the minecraft world
type Direction int

const (
	DirUp Direction = iota
	DirDown
	DirNorth
	DirSouth
	DirWest
	DirEast
)

var Directions = []Direction{
	DirUp,
	DirDown,
	DirNorth,
	DirSouth,
	DirWest,
	DirEast,
}

func DirectionFromString(str string) Direction {
	switch str {
	case "up":
		return DirUp
	case "down":
		return DirDown
	case "north":
		return DirNorth
	case "south":
		return DirSouth
	case "west":
		return DirWest
	case "east":
		return DirEast
	}
	// ¯\_(ツ)_/¯
	return -1
}

func (d Direction) Offset() (x, y, z int) {
	switch d {
	case DirUp:
		return 0, 1, 0
	case DirDown:
		return 0, -1, 0
	case DirNorth:
		return 0, 0, -1
	case DirSouth:
		return 0, 0, 1
	case DirWest:
		return -1, 0, 0
	case DirEast:
		return 1, 0, 0
	}
	return 0, 0, 0

}

func (d Direction) String() string {
	switch d {
	case DirUp:
		return "up"
	case DirDown:
		return "down"
	case DirNorth:
		return "north"
	case DirSouth:
		return "south"
	case DirWest:
		return "west"
	case DirEast:
		return "east"
	}
	return fmt.Sprintf("Direction(%d)", d)
}
