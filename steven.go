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

package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"time"

	"github.com/thinkofdeath/steven/platform"
	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/protocol/mojang"
	"github.com/thinkofdeath/steven/render"
)

var loadChan = make(chan struct{})

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	if len(os.Args) == 0 {
		fmt.Println("steven must be run via the mojang launcher")
		return
	}

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

	go func() {
		render.LoadTextures()
		initBlocks()
		loadChan <- struct{}{}
	}()

	platform.Init(platform.Handler{
		Start:  start,
		Draw:   draw,
		Move:   move,
		Rotate: rotate,
		Action: action,
	})
}

func start() {
	<-loadChan
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

func action(action platform.Action) {
	fmt.Println("action:", action)
	switch action {
	case platform.Debug:
	}
}

var maxBuilders = runtime.NumCPU() * 2

var (
	ready            bool
	freeBuilders     = maxBuilders
	completeBuilders = make(chan buildPos, maxBuilders)
	syncChan         = make(chan func(), 200)
	ticker           = time.NewTicker(time.Second / 20)
)

func draw() {
handle:
	for {
		select {
		case err := <-errorChan:
			panic(err)
		case packet := <-readChan:
			defaultHandler.Handle(packet)
		case pos := <-completeBuilders:
			c := chunkMap[chunkPosition{pos.X, pos.Z}]
			freeBuilders++
			if c != nil {
				s := c.Sections[pos.Y]
				if s != nil {
					s.building = false
				}
			}
		case f := <-syncChan:
			f()
		default:
			break handle
		}
	}

	// TODO(Think) Tidy up
	render.Camera.X += mf * math.Cos(render.Camera.Yaw-math.Pi/2) * -math.Cos(render.Camera.Pitch) * (1.0 / 7.0)
	render.Camera.Z -= mf * math.Sin(render.Camera.Yaw-math.Pi/2) * -math.Cos(render.Camera.Pitch) * (1.0 / 7.0)
	render.Camera.Y -= mf * math.Sin(render.Camera.Pitch) * (1.0 / 7.0)
	if ready {
		select {
		case <-ticker.C:
			tick()
		default:
		}
	}

	// Search for 'dirty' chunk sections and start building
	// them if we have any builders free. To prevent race conditions
	// two flags are used, dirty and building, to allow a second
	// build to be requested whilst the chunk is still building
	// without either losing the change or having two builds
	// for the same section going on at once (where the second
	// could finish quicker causing the old version to be
	// displayed.
dirtyClean:
	for _, c := range sortedChunks() {
		for _, s := range c.Sections {
			if s == nil {
				continue
			}
			if freeBuilders <= 0 {
				break dirtyClean
			}
			if s.dirty && !s.building {
				freeBuilders--
				s.dirty = false
				s.building = true
				s.build(completeBuilders)
			}
		}
	}

	render.Draw()
}

// tick is called 20 times a second (bar any preformance issues).
// Minecraft is built around this fact so we have to follow it
// as well.
func tick() {
	// Now you may be wondering why we have to spam movement
	// packets (any of the Player* move/look packets) 20 times
	// a second instead of only sending when something changes.
	// This is because the server only ticks certain parts of
	// the player when a movement packet is recieved meaning
	// if we sent them any slower health regen would be slowed
	// down as well and various other things too (potions, speed
	// hack check). This also has issues if we send them too
	// fast as well since we will regen health at much faster
	// rates than normal players and some modded servers will
	// (correctly) detect this as cheating. Its Minecraft
	// what did you expect?
	// TODO(Think) Use the smaller packets when possible
	writeChan <- &protocol.PlayerPositionLook{
		X:     render.Camera.X,
		Y:     render.Camera.Y,
		Z:     render.Camera.Z,
		Yaw:   float32(-render.Camera.Yaw * (180 / math.Pi)),
		Pitch: float32((-render.Camera.Pitch - math.Pi) * (180 / math.Pi)),
	}
}
