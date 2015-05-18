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

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/protocol/mojang"
	"github.com/thinkofdeath/steven/resource"
	"github.com/thinkofdeath/steven/ui"
	"github.com/thinkofdeath/steven/ui/scene"
)

type loginScreen struct {
	baseUI
	scene *scene.Type
	logo  uiLogo

	user    *textBox
	pass    *textBox
	focused *textBox

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

	ls.user = newTextBox(0, -20, 400, 40)
	ls.user.back.Attach(ui.Middle, ui.Center)
	ls.user.add(ls.scene)
	label := ui.NewText("Username/Email:", 0, -18, 255, 255, 255).Attach(ui.Top, ui.Left)
	label.AttachTo(ls.user.back)
	ls.scene.AddDrawable(label)
	ls.user.back.ClickFunc = func() {
		if ls.focused != nil {
			ls.focused.Focused = false
		}
		ls.user.Focused = true
		ls.focused = ls.user
	}

	ls.pass = newTextBox(0, 40, 400, 40)
	ls.pass.back.Attach(ui.Middle, ui.Center)
	ls.pass.add(ls.scene)
	label = ui.NewText("Password:", 0, -18, 255, 255, 255).Attach(ui.Top, ui.Left)
	label.AttachTo(ls.pass.back)
	ls.scene.AddDrawable(label)
	ls.pass.back.ClickFunc = func() {
		if ls.focused != nil {
			ls.focused.Focused = false
		}
		ls.pass.Focused = true
		ls.focused = ls.pass
	}
	ls.pass.Password = true

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

func (ls *loginScreen) postLogin(p mojang.Profile, err error) {
	if err != nil {
		if me, ok := err.(mojang.Error); ok {
			ls.loginError.Update(me.Message)
		} else {
			ls.loginError.Update(err.Error())
		}
		ls.loginBtn.SetDisabled(false)
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
	ls.loginBtn.SetDisabled(true)
	ls.loginTxt.Update("Logging in...")
	go func() {
		p, err := mojang.Refresh(Config.Profile, Config.ClientToken)
		syncChan <- func() { ls.postLogin(p, err) }
	}()
}

func (ls *loginScreen) login() {
	ls.loginError.Update("")
	ls.loginBtn.SetDisabled(true)
	ls.loginTxt.Update("Logging in...")
	go func() {
		p, err := mojang.Login(ls.user.input, ls.pass.input, Config.ClientToken)
		syncChan <- func() { ls.postLogin(p, err) }
	}()
}

func (ls *loginScreen) tick(delta float64) {
	if ls.loginBtn.Disabled() {
		if ls.focused != nil {
			ls.focused.Focused = false
			ls.focused = nil
		}
	}
	ls.user.tick(delta)
	ls.pass.tick(delta)
	ls.logo.tick(delta)
}

func (ls *loginScreen) handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if ls.focused == nil {
		return
	}

	if (key == glfw.KeyEnter || key == glfw.KeyTab) && action == glfw.Release {
		if ls.focused == ls.user {
			ls.user.Focused = false
			ls.focused = ls.pass
			ls.pass.Focused = true
		} else if ls.focused == ls.pass {
			ls.pass.Focused = false
			ls.focused = nil
			ls.login()
		}
		return
	}

	if key == glfw.KeyEscape && action == glfw.Release {
		ls.focused.Focused = false
		ls.focused = nil
	}

	ls.focused.handleKey(w, key, scancode, action, mods)
}

func (ls *loginScreen) handleChar(w *glfw.Window, char rune) {
	if ls.focused == nil {
		return
	}
	ls.focused.handleChar(w, char)
}

func (ls *loginScreen) remove() {
	ls.scene.Hide()
	window.SetKeyCallback(onKey)
	window.SetCharCallback(nil)
}
