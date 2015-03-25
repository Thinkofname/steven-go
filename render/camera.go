package render

import "math"

var Camera = &ClientCamera{}

func init() {
	Camera.Pitch = math.Pi
}

type ClientCamera struct {
	X, Y, Z    float64
	Yaw, Pitch float64
}
