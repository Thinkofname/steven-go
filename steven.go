package main

import (
	"math"
	"os"
	"runtime"

	"github.com/thinkofdeath/steven/platform"
	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/protocol/mojang"
	"github.com/thinkofdeath/steven/render"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	// Can't use flags as we need to support a weird flag
	// format
	var username, uuid, accessToken, server string

	for i, arg := range os.Args {
		switch arg {
		case "--username":
			username = os.Args[i+1]
		case "--uuid":
			uuid = os.Args[i+1]
		case "--accessToken":
			accessToken = os.Args[i+1]
		case "--server":
			server = os.Args[i+1]
		}
	}

	// Start connecting whilst starting the renderer
	go startConnection(mojang.Profile{
		Username:    username,
		ID:          uuid,
		AccessToken: accessToken,
	}, server)

	platform.Init(platform.Handler{
		Start:  start,
		Draw:   draw,
		Move:   move,
		Rotate: rotate,
	})
}

func start() {
	render.Start()
}

func rotate(x, y float64) {
	render.Camera.Yaw -= x
	render.Camera.Pitch -= y
}

var mf, ms float64

func move(f, s float64) {
	mf, ms = f, s
}

var ready bool
var i int

func draw() {
handle:
	for {
		select {
		case err := <-errorChan:
			panic(err)
		case packet := <-readChan:
			defaultHandler.Handle(packet)
		default:
			break handle
		}
	}

	render.Camera.X += mf * math.Cos(render.Camera.Yaw-math.Pi/2) * -math.Cos(render.Camera.Pitch) * (1.0 / 10.0)
	render.Camera.Z -= mf * math.Sin(render.Camera.Yaw-math.Pi/2) * -math.Cos(render.Camera.Pitch) * (1.0 / 10.0)
	render.Camera.Y -= mf * math.Sin(render.Camera.Pitch) * (1.0 / 10.0)
	i++
	if ready && i%3 == 0 {
		writeChan <- &protocol.PlayerPositionLook{
			X:     render.Camera.X,
			Y:     render.Camera.Y,
			Z:     render.Camera.Z,
			Yaw:   float32(-render.Camera.Yaw * (180 / math.Pi)),
			Pitch: float32((-render.Camera.Pitch - math.Pi) * (180 / math.Pi)),
		}
	}

	render.Draw()
}
