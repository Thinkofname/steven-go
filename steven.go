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
	"github.com/thinkofdeath/steven/protocol/mojang"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/resource"
	"github.com/thinkofdeath/steven/ui"
)

var (
	profile       mojang.Profile
	server        string
	loadChan      = make(chan struct{})
	currentScreen screen
	connected     bool

	stevenBuildVersion string = "dev"
)

func stevenVersion() string {
	return fmt.Sprintf("%s-%s", resource.ResourcesVersion, stevenBuildVersion)
}

type screen interface {
	init()
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
		Client.scene.Hide()
		Client.hotbarScene.Hide()
		lockMouse = false
		window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		for i := range Client.KeyState {
			Client.KeyState[i] = false
		}
		s.init()
	} else {
		Client.scene.Show()
		Client.hotbarScene.Show()
	}
}

func Main(username, uuid, accessToken, s string) {
	profile = mojang.Profile{
		Username:    username,
		ID:          uuid,
		AccessToken: accessToken,
	}
	server = s

	for _, pck := range Config.Game.ResourcePacks {
		resource.LoadZip(pck)
	}
	loadBiomes()

	// Done on its own goroutine so the connection
	// + window opening, can be done in parallel
	render.LoadTextures()
	go func() {
		initBlocks()
		loadChan <- struct{}{}
	}()

	if profile.IsComplete() && server != "" {
		// Start connecting whilst starting the renderer
		connect()
	}

	setUIScale()

	startWindow()
}

func connect() {
	initClient()
	connected = true
	disconnectReason.Value = nil
	Client.network.Connect(profile, server)
	server = ""
}

func setUIScale() {
	switch Config.Game.UIScale {
	case uiAuto:
		ui.DrawMode = ui.Scaled
		ui.Scale = 1.0
	case uiSmall:
		ui.DrawMode = ui.Unscaled
		ui.Scale = 0.4
	case uiMedium:
		ui.DrawMode = ui.Unscaled
		ui.Scale = 0.6
	case uiLarge:
		ui.DrawMode = ui.Unscaled
		ui.Scale = 1.0
	}
	ui.ForceDraw()
}

func start() {
	<-loadChan
	if Client == nil {
		initClient()
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
	if Client.Pitch < (math.Pi/2)+0.01 {
		Client.Pitch = (math.Pi / 2) + 0.01
	}
	if Client.Pitch > (math.Pi/2)*3-0.01 {
		Client.Pitch = (math.Pi/2)*3 - 0.01
	}
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
		case err := <-Client.network.Error():
			if !connected {
				continue
			}
			connected = false

			Client.network.Close()
			fmt.Printf("Disconnected: %s\n", err)
			// Reset the ready state to stop packets from being
			// sent.
			ready = false
			if err != errManualDisconnect && disconnectReason.Value == nil {
				txt := &chat.TextComponent{Text: err.Error()}
				txt.Color = chat.Red
				disconnectReason.Value = txt
			}

			if Client.entity != nil && Client.entityAdded {
				Client.entityAdded = false
				Client.entities.container.RemoveEntity(Client.entity)
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
handle:
	for {
		select {
		case packet := <-Client.network.Read():
			defaultHandler.Handle(packet)
		case pos := <-completeBuilders:
			freeBuilders++
			if c := chunkMap[chunkPosition{pos.X, pos.Z}]; c != nil {
				if s := c.Sections[pos.Y]; s != nil {
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

	if ready && Client != nil {
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
