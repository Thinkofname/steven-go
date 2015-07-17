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

	"github.com/thinkofdeath/steven/console"
	"github.com/thinkofdeath/steven/protocol/mojang"
	"github.com/thinkofdeath/steven/ui"
	"github.com/thinkofdeath/steven/ui/scene"
)

type loginScreen struct {
	baseUI
	scene *scene.Type
	logo  uiLogo

	User       *ui.TextBox `uiID:"user"`
	Pass       *ui.TextBox `uiID:"pass"`
	LoginBtn   *ui.Button  `uiID:"btnLogin"`
	LoginTxt   *ui.Text    `uiID:"txtLogin"`
	LoginError *ui.Text    `uiID:"txtError"`
}

func newLoginScreen() *loginScreen {
	ls := &loginScreen{}
	if clientToken.Value() == "" {
		data := make([]byte, 16)
		crand.Read(data)
		clientToken.SetValue(hex.EncodeToString(data))
	}

	ls.scene = scene.LoadScene("steven", "login", ls)
	ls.logo.init(ls.scene)
	uiFooter(ls.scene)

	ls.scene.Show()
	if getProfile().IsComplete() {
		ls.refresh()
	}

	return ls
}

func (ls *loginScreen) postLogin(p mojang.Profile, err error) {
	if err != nil {
		if me, ok := err.(mojang.Error); ok {
			ls.LoginError.Update(me.Message)
		} else {
			ls.LoginError.Update(err.Error())
		}
		ls.LoginBtn.SetDisabled(false)
		ls.LoginTxt.Update("Login")
		return
	}
	clientUsername.SetValue(p.Username)
	clientUUID.SetValue(p.ID)
	clientAccessToken.SetValue(p.AccessToken)

	setScreen(newServerList())
	console.ExecConf("autoexec.cfg")
}

func (ls *loginScreen) refresh() {
	ls.LoginError.Update("")
	ls.LoginBtn.SetDisabled(true)
	ls.LoginTxt.Update("Logging in...")
	go func() {
		p, err := mojang.Refresh(getProfile(), clientToken.Value())
		syncChan <- func() { ls.postLogin(p, err) }
	}()
}

func (ls *loginScreen) Login() {
	ls.LoginError.Update("")
	ls.LoginBtn.SetDisabled(true)
	ls.LoginTxt.Update("Logging in...")
	go func() {
		p, err := mojang.Login(ls.User.Value(), ls.Pass.Value(), clientToken.Value())
		syncChan <- func() { ls.postLogin(p, err) }
	}()
}

func (ls *loginScreen) tick(delta float64) {
	ls.logo.tick(delta)
}

func (ls *loginScreen) remove() {
	ls.scene.Hide()
}
