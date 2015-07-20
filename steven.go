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
	"github.com/thinkofdeath/steven/console"
	"github.com/thinkofdeath/steven/format"
	"github.com/thinkofdeath/steven/protocol/mojang"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/resource"
	"github.com/thinkofdeath/steven/ui"
)

var (
	connected bool
	skipLogin bool

	stevenBuildVersion string = "dev"

	clientUsername = console.NewStringVar("cl_username", "", console.Serializable).Doc(`
cl_username is the username that the client will use to connect
to servers.
`)
	clientUUID = console.NewStringVar("cl_uuid", "", console.Serializable).Doc(`
cl_uuid is the uuid of the client. This is unique to a player 
unlike their username.
`)
	clientAccessToken = console.NewStringVar("auth_token", "", console.Serializable).Doc(`
auth_token is the token used for this session to auth to servers
or relogin to this account.
`)
	clientToken = console.NewStringVar("auth_client_token", "", console.Serializable).Doc(`
auth_client_token is a token that stays static between sessions.
Used to identify this client vs others.
`)
)

func init() {
	console.NewStringVar("cl_version", stevenBuildVersion).Doc(`
cl_version returns the build information compiled into steven.
Returns dev if not set.
`)
	console.NewStringVar("cl_mc_version", resource.ResourcesVersion).Doc(`
cl_mc_version returns the minecraft version steven currently supports.
`)
	console.Register("connect %", connect)
}

func stevenVersion() string {
	return fmt.Sprintf("%s-%s", resource.ResourcesVersion, stevenBuildVersion)
}

func Main(username, uuid, accessToken string) {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	console.ExecConf("conf.cfg")

	if username != "" {
		clientUsername.SetValue(username)
	}
	if uuid != "" {
		clientUUID.SetValue(uuid)
	}
	if accessToken != "" {
		clientAccessToken.SetValue(accessToken)
		skipLogin = true
	}

	initResources()

	setUIScale()

	startWindow()
}

func getProfile() mojang.Profile {
	return mojang.Profile{
		Username:    clientUsername.Value(),
		ID:          clientUUID.Value(),
		AccessToken: clientAccessToken.Value(),
	}
}

func connect(server string) {
	setScreen(nil)
	connected = true
	initClient()
	disconnectReason.Value = nil
	Client.network.Connect(getProfile(), server)
}

func start() {
	render.LoadTextures()
	initBlocks()
	con.init()

	initClient()
	fakeGen()

	if skipLogin {
		setScreen(newServerList())
		console.ExecConf("autoexec.cfg")
	} else {
		setScreen(newLoginScreen())
	}
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
			console.Text("Disconnected: %s", err)
			// Reset the ready state to stop packets from being
			// sent.
			ready = false
			if err != errManualDisconnect && disconnectReason.Value == nil {
				txt := &format.TextComponent{Text: err.Error()}
				txt.Color = format.Red
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
	con.tick(delta)
	ui.Draw(width, height, delta)

	tickAudio()

	tickClouds(delta)

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
