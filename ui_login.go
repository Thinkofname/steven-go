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
	crand "crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/protocol/mojang"
	"github.com/thinkofdeath/steven/resource"
	"github.com/thinkofdeath/steven/ui"
	"github.com/thinkofdeath/steven/ui/scene"
)

type loginScreen struct {
	scene *scene.Type
	logo  uiLogo

	username    *ui.Text
	password    *ui.Text
	focused     *ui.Text
	lastFocused *ui.Text

	inputUsername string
	inputPassword string

	loginBtn   *ui.Button
	loginTxt   *ui.Text
	loginError *ui.Text
}

func newLoginScreen() *loginScreen {
	ls := &loginScreen{
		scene: scene.New(true),
	}
	if Config.ClientToken == "" {
		data := make([]byte, 16)
		crand.Read(data)
		Config.ClientToken = hex.EncodeToString(data)
		saveConfig()
	}

	window.SetKeyCallback(ls.handleKey)
	window.SetCharCallback(ls.handleChar)
	Client.scene.Hide()
	ls.logo.init(ls.scene)

	btn := ui.NewButton(0, -20, 400, 40).Attach(ui.Middle, ui.Center)
	btn.Disabled = true
	ls.scene.AddDrawable(btn)
	label := ui.NewText("Username/Email:", 0, -18, 255, 255, 255).Attach(ui.Top, ui.Left)
	label.Parent = btn
	ls.scene.AddDrawable(label)
	text := ui.NewText("", 5, 0, 255, 255, 255).Attach(ui.Middle, ui.Left)
	text.Parent = btn
	ls.scene.AddDrawable(text)
	ls.username = text
	btn.ClickFunc = func() { ls.focused = ls.username }

	btn = ui.NewButton(0, 40, 400, 40).Attach(ui.Middle, ui.Center)
	btn.Disabled = true
	ls.scene.AddDrawable(btn)
	label = ui.NewText("Password:", 0, -18, 255, 255, 255).Attach(ui.Top, ui.Left)
	label.Parent = btn
	ls.scene.AddDrawable(label)
	text = ui.NewText("", 5, 0, 255, 255, 255).Attach(ui.Middle, ui.Left)
	text.Parent = btn
	ls.scene.AddDrawable(text)
	ls.password = text
	btn.ClickFunc = func() { ls.focused = ls.password }

	ls.loginBtn, ls.loginTxt = newButtonText("Login", 0, 100, 400, 40)
	ls.loginBtn.Attach(ui.Middle, ui.Center)
	ls.scene.AddDrawable(ls.loginBtn)
	ls.scene.AddDrawable(ls.loginTxt)
	ls.loginBtn.ClickFunc = ls.login

	ls.scene.AddDrawable(
		ui.NewText("Steven - "+resource.ResourcesVersion, 5, 5, 255, 255, 255).Attach(ui.Bottom, ui.Left),
	)
	ls.scene.AddDrawable(
		ui.NewText("Not affiliated with Mojang/Minecraft", 5, 5, 255, 200, 200).Attach(ui.Bottom, ui.Right),
	)

	ls.loginError = ui.NewText("", 0, 150, 255, 50, 50).Attach(ui.Center, ui.Middle)
	ls.scene.AddDrawable(ls.loginError)

	if Config.Profile.IsComplete() {
		ls.refresh()
	}

	return ls
}
func (ls *loginScreen) hover(x, y float64, w, h int) {
	ui.Hover(x, y, w, h)
}
func (ls *loginScreen) click(x, y float64, w, h int) {
	ui.Click(x, y, w, h)
}

func (ls *loginScreen) postLogin(p mojang.Profile, err error) {
	if err != nil {
		if me, ok := err.(mojang.Error); ok {
			ls.loginError.Update(me.Message)
		} else {
			ls.loginError.Update(err.Error())
		}
		ls.loginBtn.Disabled = false
		ls.loginTxt.Update("Login")
		return
	}
	profile = p
	Config.Profile = p
	saveConfig()
	if server == "" {
		setScreen(newServerList())
	} else {
		initClient()
		Client.init()
		connect()
		setScreen(nil)
	}
}

func (ls *loginScreen) refresh() {
	ls.loginError.Update("")
	ls.loginBtn.Disabled = true
	ls.loginTxt.Update("Logging in...")
	go func() {
		p, err := mojang.Refresh(Config.Profile, Config.ClientToken)
		syncChan <- func() { ls.postLogin(p, err) }
	}()
}

func (ls *loginScreen) login() {
	ls.loginError.Update("")
	ls.loginBtn.Disabled = true
	ls.loginTxt.Update("Logging in...")
	go func() {
		p, err := mojang.Login(ls.inputUsername, ls.inputPassword, Config.ClientToken)
		syncChan <- func() { ls.postLogin(p, err) }
	}()
}

func (ls *loginScreen) tick(delta float64) {
	if ls.loginBtn.Disabled {
		ls.focused = nil
	}
	if ls.focused == ls.username {
		ls.username.Update(ls.inputUsername + "|")
	} else if ls.focused == ls.password {
		ls.password.Update(strings.Repeat("*", len(ls.inputPassword)) + "|")
	}
	if ls.lastFocused != nil && ls.lastFocused != ls.focused {
		if ls.lastFocused == ls.username {
			ls.username.Update(ls.inputUsername)
		} else if ls.lastFocused == ls.password {
			ls.password.Update(strings.Repeat("*", len(ls.inputPassword)))
		}
	}
	ls.lastFocused = ls.focused

	ls.logo.tick(delta)
}

func (ls *loginScreen) handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if ls.focused == nil {
		return
	}

	if (key == glfw.KeyEnter || key == glfw.KeyTab) && action == glfw.Release {
		if ls.focused == ls.username {
			ls.focused = ls.password
		} else if ls.focused == ls.password {
			ls.focused = nil
			ls.login()
		}
	}

	if key == glfw.KeyEscape && action == glfw.Release {
		ls.focused = nil
	}

	if key == glfw.KeyBackspace && action != glfw.Release {
		if ls.focused == ls.username {
			if len(ls.inputUsername) > 0 {
				ls.inputUsername = ls.inputUsername[:len(ls.inputUsername)-1]
			}
			return
		}
		if len(ls.inputPassword) > 0 {
			ls.inputPassword = ls.inputPassword[:len(ls.inputPassword)-1]
		}
	}
}

func (ls *loginScreen) handleChar(w *glfw.Window, char rune) {
	if ls.focused == nil {
		return
	}
	if ls.focused == ls.username {
		ls.inputUsername += string(char)
		return
	}
	ls.inputPassword += string(char)
}

func (ls *loginScreen) remove() {
	ls.scene.Hide()
	window.SetKeyCallback(onKey)
	window.SetCharCallback(nil)
}
