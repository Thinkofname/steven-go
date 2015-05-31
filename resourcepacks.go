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

	"github.com/thinkofdeath/steven/console"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/resource"
	"github.com/thinkofdeath/steven/resource/locale"
	"github.com/thinkofdeath/steven/ui"
)

func initResources() {
	var progressBar *ui.Image
	var progressText *ui.Text
	resource.Init(func(progress float64, done bool) {
		fmt.Printf("Progress: %0.2f %t\n", progress, done)
		if !done {
			if progressBar == nil {
				progressBar = ui.NewImage(render.GetTexture("solid"), 0, 0, 854, 21, 0, 0, 1, 1, 0, 125, 0)
				ui.AddDrawable(progressBar.Attach(ui.Top, ui.Left))
				progressText = ui.NewText("", 1, 1, 255, 255, 255)
				ui.AddDrawable(progressText.Attach(ui.Top, ui.Left))
			}
			progressText.Update(fmt.Sprintf("Downloading: %d/100", int(100*progress)))
			width, _ := window.GetFramebufferSize()
			sw := 854 / float64(width)
			if ui.DrawMode == ui.Unscaled {
				sw = ui.Scale
				progressBar.SetWidth((854 / sw) * progress)
			} else {
				progressBar.SetWidth(float64(width) * progress)
			}
		} else {
			if progressBar != nil {
				progressBar.Remove()
				progressText.Remove()
			}
			reloadResources()
		}
	}, syncChan)
}

func AddPack(path string) {
	console.Text("Adding pack " + path)
	if err := resource.LoadZip(path); err != nil {
		fmt.Println("Failed to load pack", path)
		return
	}
	Config.Game.ResourcePacks = append(Config.Game.ResourcePacks, path)
	saveConfig()
	reloadResources()
}

func RemovePack(path string) {
	console.Text("Removing pack " + path)
	resource.RemovePack(path)
	for i, pck := range Config.Game.ResourcePacks {
		if pck == path {
			Config.Game.ResourcePacks = append(Config.Game.ResourcePacks[:i], Config.Game.ResourcePacks[i+1:]...)
			break
		}
	}
	saveConfig()
	reloadResources()
}

func reloadResources() {
	console.Text("Bringing everything to a stop")
	for freeBuilders < maxBuilders {
		select {
		case pos := <-completeBuilders:
			freeBuilders++
			if c := chunkMap[chunkPosition{pos.X, pos.Z}]; c != nil {
				if s := c.Sections[pos.Y]; s != nil {
					s.building = false
				}
			}
		}
	}
	locale.Clear()
	render.LoadSkinBuffer()
	modelCache = map[string]*model{}
	console.Text("Reloading textures")
	render.LoadTextures()
	console.Text("Reloading biomes")
	loadBiomes()
	ui.ForceDraw()
	console.Text("Reloading blocks")
	reinitBlocks()
	console.Text("Marking chunks for rebuild")
	for _, c := range chunkMap {
		for _, s := range c.Sections {
			if s != nil {
				s.dirty = true
			}
		}
	}
	console.Text("Rebuilding static models")
	render.RefreshStaticModels()
	console.Text("Reloading inventory")
	Client.playerInventory.Update()
}
