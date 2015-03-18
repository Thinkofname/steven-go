package main

import (
	"github.com/thinkofdeath/steven/platform"
	"github.com/thinkofdeath/steven/render"
)

func main() {
	platform.Init(platform.Handler{
		Start: start,
		Draw:  draw,
	})
}

func start() {
	render.Start()
}

func draw() {
	render.Draw()
}
