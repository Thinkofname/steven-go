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
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/ui"
	"github.com/thinkofdeath/steven/ui/scene"
)

var (
	InvPlayer = playerInventory{}
)

type InventoryType interface {
	Draw(scene *scene.Type, inv *Inventory)
}

type Inventory struct {
	Type InventoryType

	Items []*ItemStack

	scene *scene.Type
}

func NewInventory(ty InventoryType, size int) *Inventory {
	return &Inventory{
		Type:  ty,
		Items: make([]*ItemStack, size),
		scene: scene.New(true),
	}
}

func (inv *Inventory) Update() {
	inv.scene.Hide()
	inv.scene = scene.New(true)
	inv.Type.Draw(inv.scene, inv)
}

func (inv *Inventory) Close() {
	inv.scene.Hide()
}

func (inv *Inventory) Hide() {
	inv.scene.Hide()
}

func (inv *Inventory) Show() {
	inv.scene.Show()
}

// Player

type playerInventory struct {
}

const invPlayerHotbarOffset = 36

func (playerInventory) Draw(s *scene.Type, inv *Inventory) {
	// Slots 36-44 are the hotbar
	Client.hotbarScene.Hide()
	Client.hotbarScene = scene.New(true)
	hs := Client.hotbarScene
	for i := invPlayerHotbarOffset; i < invPlayerHotbarOffset+9; i++ {
		if inv.Items[i] == nil {
			continue
		}
		img := ui.NewImage(render.GetTexture("missing_texture"), 6+40*float64(i-36), 6, 32, 32, 0, 0, 1, 1, 255, 255, 255)
		img.AttachTo(Client.hotbar)
		hs.AddDrawable(img.Attach(ui.Top, ui.Left))
	}
}
