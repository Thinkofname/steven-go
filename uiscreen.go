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

import "github.com/go-gl/glfw/v3.1/glfw"

var currentScreen screen

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
