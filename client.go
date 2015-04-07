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
	"math"
	"sync"

	"github.com/thinkofdeath/steven/render"
)

const (
	playerHeight = 1.62
)

var Client ClientState

type ClientState struct {
	sync.Mutex
	X, Y, Z    float64
	Yaw, Pitch float64
}

func (c *ClientState) renderTick(delta float64) {
	c.Lock()
	defer c.Unlock()
	c.X += mf * math.Cos(c.Yaw-math.Pi/2) * -math.Cos(c.Pitch) * delta * 0.2
	c.Z -= mf * math.Sin(c.Yaw-math.Pi/2) * -math.Cos(c.Pitch) * delta * 0.2
	c.Y -= mf * math.Sin(c.Pitch) * delta * 0.2

	// Copy to the camera
	render.Camera.X = c.X
	render.Camera.Y = c.Y + playerHeight
	render.Camera.Z = c.Z
	render.Camera.Yaw = c.Yaw
	render.Camera.Pitch = c.Pitch
}

func (c *ClientState) tick() {
	c.Lock()
	defer c.Unlock()
}
