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

	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/resource"
	"github.com/thinkofdeath/steven/ui"
)

func AddPack(path string) {
	fmt.Println("Adding pack " + path)
	if err := resource.LoadZip(path); err != nil {
		fmt.Println("Failed to load pack", path)
		return
	}
	Config.Game.ResourcePacks = append(Config.Game.ResourcePacks, path)
	saveConfig()
	reloadResources()
}

func RemovePack(path string) {
	fmt.Println("Removing pack " + path)
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
	fmt.Println("Bringing everything to a stop")
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
	modelCache = map[string]*model{}
	fmt.Println("Reloading textures")
	render.LoadTextures()
	fmt.Println("Reloading biomes")
	loadBiomes()
	ui.ForceDraw()
	fmt.Println("Reloading blocks")
	reinitBlocks()
	fmt.Println("Marking chunks for rebuild")
	for _, c := range chunkMap {
		for _, s := range c.Sections {
			if s != nil {
				s.dirty = true
			}
		}
	}
	fmt.Println("Rebuilding static models")
	render.RefreshStaticModels()
	fmt.Println("Reloading inventory")
	Client.playerInventory.Update()
}
