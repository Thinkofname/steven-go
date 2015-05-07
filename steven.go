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
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/chat"
	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/protocol/mojang"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/ui"
)

var (
	profile       mojang.Profile
	server        string
	loadChan      = make(chan struct{})
	currentScreen screen
	connected     bool
)

type screen interface {
	tick(delta float64)
	hover(x, y float64, w, h int)
	click(down bool, x, y float64, w, h int)
	remove()
}

func setScreen(s screen) {
	if currentScreen != nil {
		currentScreen.remove()
	}
	currentScreen = s
	if s != nil {
		lockMouse = false
		window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		for i := range Client.KeyState {
			Client.KeyState[i] = false
		}
	} else {
		Client.scene.Show()
	}
}

func Main(username, uuid, accessToken, s string) {
	profile = mojang.Profile{
		Username:    username,
		ID:          uuid,
		AccessToken: accessToken,
	}
	server = s

	go func() {
		render.LoadTextures()
		initBlocks()
		loadChan <- struct{}{}
	}()

	if profile.IsComplete() && server != "" {
		// Start connecting whilst starting the renderer
		connect()
	} else {
		Client.valid = false
	}

	startWindow()
}

func connect() {
	connected = true
	disconnectReason.Value = nil
	go startConnection(profile, server)
	server = ""
}

func start() {
	<-loadChan
	if Client.valid {
		Client.init()
	} else {
		fakeGen()
		if !profile.IsComplete() {
			setScreen(newLoginScreen())
		} else {
			setScreen(newServerList())
		}
	}
	render.FOV = Config.Render.FOV
	render.Start()
}

func rotate(x, y float64) {
	Client.Yaw -= x
	Client.Pitch -= y
}

var maxBuilders = runtime.NumCPU() * 2

var (
	ready            bool
	freeBuilders     = maxBuilders
	completeBuilders = make(chan buildPos, maxBuilders)
	syncChan         = make(chan func(), 200)
	ticker           = time.NewTicker(time.Second / 20)
	lastFrame        = time.Now()
)

func handleErrors() {
handle:
	for {
		select {
		case err := <-errorChan:
			if !connected {
				continue
			}
			connected = false

			// Replace the writer ready for the next connection
			close(writeChan)
			writeChan = make(chan protocol.Packet, 200)

			conn.Close()
			fmt.Printf("Disconnected: %s\n", err)
			// Reset the ready state to stop packets from being
			// sent.
			ready = false
			if err != errManualDisconnect && disconnectReason.Value == nil {
				txt := &chat.TextComponent{Text: err.Error()}
				txt.Color = chat.Red
				disconnectReason.Value = txt
			}
			setScreen(newServerList())
		default:
			break handle
		}
	}
}

func draw() {
	now := time.Now()
	diff := now.Sub(lastFrame)
	lastFrame = now
	delta := float64(diff.Nanoseconds()) / (float64(time.Second) / 60)
	delta = math.Min(delta, 1.6)
handle:
	for {
		select {
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
	handleErrors()

	width, height := window.GetFramebufferSize()

	if currentScreen != nil {
		currentScreen.tick(delta)
	}

	if ready && Client.valid {
		Client.renderTick(delta)
		select {
		case <-ticker.C:
			tick()
		default:
		}
	} else {
		render.Camera.Yaw += 0.005 * delta
		if render.Camera.Yaw > math.Pi*2 {
			render.Camera.Yaw = 0
		}
	}
	ui.Draw(width, height, delta)

	render.Draw(width, height, delta)
	chunks := sortedChunks()

	// Search for 'dirty' chunk sections and start building
	// them if we have any builders free. To prevent race conditions
	// two flags are used, dirty and building, to allow a second
	// build to be requested whilst the chunk is still building
	// without either losing the change or having two builds
	// for the same section going on at once (where the second
	// could finish quicker causing the old version to be
	// displayed.
dirtyClean:
	for _, c := range chunks {
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
}

// tick is called 20 times a second (bar any preformance issues).
// Minecraft is built around this fact so we have to follow it
// as well.
func tick() {
	Client.tick()
}
