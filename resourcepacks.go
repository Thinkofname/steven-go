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
	"bytes"
	"fmt"
	"strings"

	"github.com/thinkofdeath/steven/console"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/resource"
	"github.com/thinkofdeath/steven/resource/locale"
	"github.com/thinkofdeath/steven/ui"
)

var resourcePacks = console.NewStringVar("cl_resource_packs", "", console.Serializable).Doc(`
cl_resource_packs is a comma seperated list of resource packs 
that are currently enabled.
`)

func initResources() {
	var pBar *progressBar
	resource.Init(func(progress float64, done bool) {
		if !done {
			if pBar == nil {
				pBar = newProgressBar()
			}
			pBar.update(progress, fmt.Sprintf("Downloading textures: %v/100", int(100*progress)))
		} else {
			if pBar != nil {
				pBar.remove()
			}
			reloadResources()
		}
	}, syncChan)

	for _, pck := range strings.Split(resourcePacks.Value(), ",") {
		if pck == "" {
			continue
		}
		resource.LoadZip(pck)
	}
	locale.Clear()
	loadBiomes()
	render.LoadSkinBuffer()
}

func AddPack(path string) {
	console.Text("Adding pack " + path)
	if err := resource.LoadZip(path); err != nil {
		fmt.Println("Failed to load pack", path)
		return
	}
	if resourcePacks.Value() != "" {
		resourcePacks.SetValue(resourcePacks.Value() + "," + path)
	} else {
		resourcePacks.SetValue(path)
	}
	reloadResources()
}

func RemovePack(path string) {
	console.Text("Removing pack " + path)
	resource.RemovePack(path)
	var buf bytes.Buffer
	for _, pck := range strings.Split(resourcePacks.Value(), ",") {
		if pck != path {
			buf.WriteString(pck)
			buf.WriteRune(',')
		}
	}
	val := buf.String()
	if strings.HasPrefix(val, ",") {
		val = val[:len(val)-1]
	}
	resourcePacks.SetValue(val)
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
	cloudImage = nil
}
