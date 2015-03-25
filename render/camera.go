package render

import "math"

// Camera is the main camera for the renderer
var Camera = &ClientCamera{}

func init() {
	Camera.Pitch = math.Pi
}

// ClientCamera is a camera within the world
type ClientCamera struct {
	X, Y, Z    float64
	Yaw, Pitch float64
}
