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
	"encoding/json"
	"os"

	"github.com/thinkofdeath/steven/protocol/mojang"
)

var Config ConfigData

type ConfigData struct {
	Profile     mojang.Profile
	ClientToken string

	Servers []ConfigServer

	Render struct {
		Samples int
		FOV     int
		VSync   bool
	}
	Game struct {
		MouseSensitivity int
		UIScale          string
	}
}

const (
	uiAuto   = "auto"
	uiSmall  = "small"
	uiMedium = "medium"
	uiLarge  = "large"
)

type ConfigServer struct {
	Name    string
	Address string
}

func init() {
	// Defaults
	Config.Render.FOV = 80
	Config.Render.VSync = true
	Config.Game.MouseSensitivity = 2000
	Config.Game.UIScale = "auto"

	f, err := os.Open("config.json")
	if err != nil {
		return
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&Config)
	if err != nil {
		panic(err)
	}
	saveConfig()
}

func saveConfig() {
	f, err := os.Create("config.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data, err := json.MarshalIndent(&Config, "", "    ")
	if err != nil {
		panic(err)
	}
	_, err = f.Write(data)
	if err != nil {
		panic(err)
	}
}
