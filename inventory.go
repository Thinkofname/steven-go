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

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/ui"
	"github.com/thinkofdeath/steven/ui/scene"
)

var (
	InvPlayer = playerInventory{}

	invScreen = &inventoryScreen{}
)

type InventoryType interface {
	Draw(scene *scene.Type, inv *Inventory)
}

type Inventory struct {
	Type InventoryType
	ID   int

	Items []*ItemStack

	scene *scene.Type
}

func NewInventory(ty InventoryType, id, size int) *Inventory {
	return &Inventory{
		Type:  ty,
		ID:    id,
		Items: make([]*ItemStack, size),
		scene: scene.New(true),
	}
}

func (inv *Inventory) Update() {
	was := inv.scene.IsVisible()
	inv.scene.Hide()
	inv.scene = scene.New(was)
	inv.Type.Draw(inv.scene, inv)
}

func (inv *Inventory) Close() {
	inv.scene.Hide()
	Client.network.Write(&protocol.CloseWindow{ID: byte(inv.ID)})
}

func (inv *Inventory) Hide() {
	inv.scene.Hide()
}

func (inv *Inventory) Show() {
	inv.scene.Show()
}

func openInventory(inv *Inventory) {
	Client.activeInventory = inv
	Client.activeInventory.Show()
	Client.activeInventory.Update()
	setScreen(invScreen)
}

func closeInventory() {
	if inv := Client.activeInventory; inv != nil {
		inv.Close()
		Client.activeInventory = nil
		setScreen(nil)
	}
	Client.playerInventory.Update()
}

type inventoryScreen struct {
	prev glfw.KeyCallback

	activeSlot int
	inWindow   bool

	cursorItem     *ItemStack
	cursorIcon     *ui.Container
	lastMX, lastMY float64
	scene          *scene.Type
}

func (i *inventoryScreen) init() {
	i.prev = window.SetKeyCallback(i.onKey)
	i.activeSlot = -1
	i.cursorItem = nil
	if i.scene != nil {
		i.scene.Hide()
	}
	i.scene = scene.New(true)
}
func (i *inventoryScreen) tick(delta float64) {}

func (i *inventoryScreen) hover(x, y float64, w, h int) {
	i.lastMX, i.lastMY = x, y
	if i.cursorIcon != nil {
		i.cursorIcon.SetX(x - 16)
		i.cursorIcon.SetY(y - 16)
	}
	ui.Hover(x, y, w, h)
}
func (i *inventoryScreen) click(down bool, x, y float64, w, h int) {
	if down {
		if i.activeSlot != -1 {
			item := Client.activeInventory.Items[i.activeSlot]
			Client.activeInventory.Items[i.activeSlot] = i.cursorItem
			Client.network.Write(&protocol.ClickWindow{
				ID:           byte(Client.activeInventory.ID),
				Slot:         int16(i.activeSlot),
				Button:       0,
				Mode:         0,
				ActionNumber: 42,
				ClickedItem:  ItemStackToProtocol(nil),
			})
			i.setCursor(item)
		} else if !i.inWindow {
			Client.network.Write(&protocol.ClickWindow{
				ID:           byte(Client.activeInventory.ID),
				Slot:         int16(-999),
				Button:       0,
				Mode:         0,
				ActionNumber: 42,
				ClickedItem:  ItemStackToProtocol(nil),
			})
			i.setCursor(nil)
		}
		return
	}
	ui.Click(x, y, w, h)
}

func (i *inventoryScreen) setCursor(item *ItemStack) {
	i.scene.Hide()
	i.scene = scene.New(true)
	i.cursorItem = item
	i.cursorIcon = nil
	if item != nil {
		i.cursorIcon = createItemIcon(item, i.scene, i.lastMX-16, i.lastMY-16)
		i.cursorIcon.SetLayer(100)
	}
	Client.activeInventory.Update()
}

func (i *inventoryScreen) onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}
	if key == glfw.KeyE || key == glfw.KeyEscape {
		closeInventory()
	}
}

func (i *inventoryScreen) remove() {
	window.SetKeyCallback(i.prev)
	i.scene.Hide()
}

// Player

type playerInventory struct {
}

const invPlayerHotbarOffset = 36

func (playerInventory) Draw(s *scene.Type, inv *Inventory) {
	full := Client.activeInventory == Client.playerInventory

	if !full {
		// Slots 36-44 are the hotbar
		Client.hotbarScene.Hide()
		Client.hotbarScene = scene.New(true)
		hs := Client.hotbarScene
		for i := invPlayerHotbarOffset; i < invPlayerHotbarOffset+9; i++ {
			if inv.Items[i] == nil {
				continue
			}
			item := inv.Items[i]
			container := createItemIcon(item, hs, 6+40*float64(i-36), 6).
				Attach(ui.Top, ui.Left)
			container.AttachTo(Client.hotbar)
		}
		return
	}

	background := ui.NewImage(
		render.GetTexture("gui/container/inventory"),
		0, 0, 176*2, 166*2,
		0, 0, 176/256.0, 166/256.0,
		255, 255, 255,
	)
	s.AddDrawable(background.Attach(ui.Middle, ui.Center))

	check := ui.NewContainer(0, 0, 176*2, 166*2)
	s.AddDrawable(check.Attach(ui.Middle, ui.Center))
	check.HoverFunc = func(over bool) {
		invScreen.inWindow = over
	}

	var slotPositions = [45][2]float64{
		0: {144, 36}, // Craft-out
		// Craft-In
		1: {88, 26},
		2: {106, 26},
		3: {88, 44},
		4: {106, 44},
		// Armor
		5: {8, 8},
		6: {8, 26},
		7: {8, 44},
		8: {8, 62},
	}
	for i := 9; i <= 35; i++ {
		x := i % 9
		y := (i / 9) - 1
		slotPositions[i] = [2]float64{
			8 + 18*float64(x), 84 + 18*float64(y),
		}
	}
	for i := 0; i < 9; i++ {
		slotPositions[i+36] = [2]float64{
			8 + 18*float64(i), 142,
		}
	}

	solid := render.GetTexture("solid")

	for i, pos := range slotPositions {
		i := i
		ctn := ui.NewContainer(pos[0]*2, pos[1]*2, 32, 32)
		ctn.AttachTo(background)
		s.AddDrawable(ctn)

		item := inv.Items[i]
		if item != nil {
			container := createItemIcon(item, s, pos[0]*2, pos[1]*2)
			container.AttachTo(background)
		} else if i >= 5 && i <= 8 {
			tex := render.GetTexture([]string{
				"items/empty_armor_slot_helmet",
				"items/empty_armor_slot_chestplate",
				"items/empty_armor_slot_leggings",
				"items/empty_armor_slot_boots",
			}[i-5])
			img := ui.NewImage(tex, pos[0]*2, pos[1]*2, 32, 32, 0, 0, 1, 1, 255, 255, 255)
			img.AttachTo(background)
			s.AddDrawable(img)
		}

		highlight := ui.NewImage(solid, pos[0]*2, pos[1]*2, 32, 32, 0, 0, 1, 1, 255, 255, 255)
		highlight.SetA(0)
		highlight.AttachTo(background)
		highlight.SetLayer(25)
		s.AddDrawable(highlight)

		ctn.HoverFunc = func(over bool) {
			if over {
				highlight.SetA(100)
				invScreen.activeSlot = i
			} else {
				highlight.SetA(0)
				if i == invScreen.activeSlot {
					invScreen.activeSlot = -1
				}
			}
		}
	}
}

func createItemIcon(item *ItemStack, scene *scene.Type, x, y float64) *ui.Container {
	mdl := getModel(item.Type.Name())

	container := ui.NewContainer(x, y, 32, 32)
	if mdl == nil || mdl.builtIn == builtInGenerated {
		var tex render.TextureInfo
		if mdl == nil {
			tex = render.GetTexture("missing_texture")

			img := ui.NewImage(tex, 0, 0, 32, 32, 0, 0, 1, 1, 255, 255, 255)
			img.AttachTo(container)
			scene.AddDrawable(img.Attach(ui.Top, ui.Left))
		} else {
			for i := 0; i < 9; i++ {
				v := fmt.Sprintf("layer%d", i)
				if _, ok := mdl.textureVars[v]; !ok {
					break
				}
				tex = mdl.lookupTexture("#" + v)

				img := ui.NewImage(tex, 0, 0, 32, 32, 0, 0, 1, 1, 255, 255, 255)
				img.AttachTo(container)
				scene.AddDrawable(img.Attach(ui.Top, ui.Left))
			}
		}
	} else if mdl.builtIn == builtInFalse {
		var blk Block
		if bt, ok := item.Type.(*blockItem); ok {
			blk = bt.block
		}
		u := modelToUI(mdl, blk)
		u.AttachTo(container)
		scene.AddDrawable(u.Attach(ui.Top, ui.Left))
	}
	if dam, ok := item.Type.(ItemDamagable); ok && dam.Damage() > 0 {
		val := 1.0 - (float64(dam.Damage()) / float64(dam.MaxDamage()))
		bar := ui.NewImage(render.GetTexture("solid"), 0, 0, 32*val, 2, 0, 0, 1, 1,
			int(255*(1.0-val)), int(255*val), 0,
		)
		bar.SetLayer(2)
		bar.AttachTo(container)
		scene.AddDrawable(bar.Attach(ui.Bottom, ui.Left))
	}
	if item.Count > 1 {
		txt := ui.NewText(fmt.Sprint(item.Count), -2, -2, 255, 255, 255).
			Attach(ui.Bottom, ui.Right)
		txt.AttachTo(container)
		txt.SetLayer(2)
		scene.AddDrawable(txt)
	}
	return container
}
