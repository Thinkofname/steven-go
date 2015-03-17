package main

import (
	"github.com/thinkofdeath/steven/platform"
	"github.com/thinkofdeath/steven/platform/gl"
)

func main() {
	platform.Init(platform.Handler{
		Start: start,
		Draw:  draw,
	})
}

func start() {

}

var test = float32(0.0)

func draw() {
	gl.ClearColor(0.0, test, 0.0, 1.0)
	test += 0.005
	if test > 1.0 {
		test = 0.0
	}
	gl.Clear(gl.ColorBufferBit)
}
