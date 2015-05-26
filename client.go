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
	"encoding/hex"
	"math"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/thinkofdeath/steven/chat"
	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/type/direction"
	"github.com/thinkofdeath/steven/type/vmath"
	"github.com/thinkofdeath/steven/ui"
	"github.com/thinkofdeath/steven/ui/scene"
)

const (
	playerHeight = 1.62
)

var Client *ClientState

func initClient() {
	if Client != nil {
		Client.scene.Hide()
		// Cleanup
		for _, e := range Client.entities.entities {
			Client.entities.container.RemoveEntity(e)
		}
		for _, e := range Client.blockBreakers {
			Client.entities.container.RemoveEntity(e)
		}
		if Client.entity != nil && Client.entityAdded {
			Client.entities.container.RemoveEntity(Client.entity)
		}
		Client.playerList.free()

		Client.playerInventory.Close()
		Client.hotbarScene.Hide()
	}
	newClient()
}

type cameraMode int

const (
	cameraNormal cameraMode = iota
	cameraBehind
	cameraFront
)

type ClientState struct {
	scene *scene.Type

	cameraMode cameraMode

	entity      *clientEntity
	entityAdded bool

	LX, LY, LZ float64
	X, Y, Z    float64
	Yaw, Pitch float64

	Health float64
	Hunger float64

	VSpeed                   float64
	KeyState                 [6]bool
	OnGround, didTouchGround bool
	isLeftDown               bool

	GameMode gameMode
	HardCore bool

	Bounds vmath.AABB

	debug struct {
		enabled  bool
		position *ui.Text
		facing   *ui.Text
		fps      *ui.Text
		memory   *ui.Text

		target     *ui.Text
		targetName *ui.Text
		targetInfo [][2]*ui.Text

		fpsValue  int
		frames    int
		lastCount time.Time
	}
	hotbar     *ui.Image
	hotbarUI   *ui.Image
	lifeUI     []*ui.Image
	lifeFillUI []*ui.Image
	foodUI     []*ui.Image
	foodFillUI []*ui.Image

	currentHotbarSlot, lastHotbarSlot int
	lastHotbarItem                    *ItemStack
	itemNameUI                        *ui.Formatted
	itemNameTimer                     float64

	network    networkManager
	chat       ChatUI
	playerList playerListUI
	entities   clientEntities

	playerInventory *Inventory
	hotbarScene     *scene.Type

	currentBreakingBlock    Block
	currentBreakingPos      Position
	maxBreakTime, breakTime float64
	swingTimer              float64
	breakEntity             BlockEntity
	blockBreakers           map[int]BlockEntity

	delta float64
}

func newClient() {
	c := &ClientState{
		Bounds: vmath.AABB{
			Min: mgl32.Vec3{-0.3, 0, -0.3},
			Max: mgl32.Vec3{0.3, 1.8, 0.3},
		},
		scene: scene.New(true),
	}
	Client = c
	c.playerInventory = NewInventory(InvPlayer, 45)
	c.hotbarScene = scene.New(true)
	c.network.init()
	c.currentBreakingBlock = Blocks.Air.Base
	c.blockBreakers = map[int]BlockEntity{}
	widgets := render.GetTexture("gui/widgets")
	icons := render.GetTexture("gui/icons")
	// Crosshair
	c.scene.AddDrawable(
		ui.NewImage(icons, 0, 0, 32, 32, 0, 0, 16.0/256.0, 16.0/256.0, 255, 255, 255).
			Attach(ui.Middle, ui.Center),
	)
	// Hotbar
	hotbar := ui.NewImage(widgets, 0, 0, 182*2, 22*2, 0, 0, 182.0/256.0, 22.0/256.0, 255, 255, 255).
		Attach(ui.Bottom, ui.Center)
	c.scene.AddDrawable(hotbar)
	c.hotbar = hotbar
	c.hotbarUI = ui.NewImage(widgets, -22*2+4, -2, 24*2, 24*2, 0, 22.0/256.0, 24.0/256.0, 24.0/256.0, 255, 255, 255).
		Attach(ui.Bottom, ui.Center)
	c.scene.AddDrawable(c.hotbarUI)

	// Hearts / Food
	for i := 0; i < 10; i++ {
		l := ui.NewImage(icons, 16*float64(i), -16-8-10, 18, 18, 16.0/256.0, 0, 9.0/256.0, 9.0/256.0, 255, 255, 255).
			Attach(ui.Top, ui.Left)
		l.AttachTo(hotbar)
		c.scene.AddDrawable(l)
		c.lifeUI = append(c.lifeUI, l)
		l = ui.NewImage(icons, 16*float64(i), -16-8-10, 18, 18, (16+9*4)/256.0, 0, 9.0/256.0, 9.0/256.0, 255, 255, 255).
			Attach(ui.Top, ui.Left)
		l.AttachTo(hotbar)
		c.scene.AddDrawable(l)
		c.lifeFillUI = append(c.lifeFillUI, l)

		f := ui.NewImage(icons, 16*float64(i), -16-8-10, 18, 18, 16.0/256.0, 27.0/256.0, 9.0/256.0, 9.0/256.0, 255, 255, 255).
			Attach(ui.Top, ui.Right)
		f.AttachTo(hotbar)
		c.scene.AddDrawable(f)
		c.foodUI = append(c.foodUI, f)
		f = ui.NewImage(icons, 16*float64(i), -16-8-10, 18, 18, (16+9*4)/256.0, 27.0/256.0, 9.0/256.0, 9.0/256.0, 255, 255, 255).
			Attach(ui.Top, ui.Right)
		f.AttachTo(hotbar)
		c.scene.AddDrawable(f)
		c.foodFillUI = append(c.foodFillUI, f)
	}

	// Exp bar
	c.scene.AddDrawable(
		ui.NewImage(icons, 0, 22*2+4, 182*2, 10, 0, 64.0/256.0, 182.0/256.0, 5.0/256.0, 255, 255, 255).
			Attach(ui.Bottom, ui.Center),
	)

	c.itemNameUI = ui.NewFormatted(chat.AnyComponent{Value: &chat.TextComponent{}}, 0, -16-8-10-16)
	c.itemNameUI.AttachTo(c.hotbar)
	c.scene.AddDrawable(c.itemNameUI.Attach(ui.Top, ui.Middle))

	c.chat.init()
	c.initDebug()
	c.playerList.init()
	c.entities.init()

	c.initEntity(false)
}

type clientEntity struct {
	positionComponent
	rotationComponent
	targetRotationComponent
	targetPositionComponent

	playerComponent
	playerModelComponent
}

func (c *ClientState) initEntity(head bool) {
	ce := &clientEntity{}
	ub, _ := hex.DecodeString(profile.ID)
	copy(ce.uuid[:], ub)
	c.entity = ce
	ce.hasHead = head
	ce.isFirstPerson = !head
	ce.manualMove = true
	ce.SetCurrentItem(c.lastHotbarItem)
}

func (c *ClientState) cycleCamera() {
	oldMode := c.cameraMode
	c.cameraMode = (c.cameraMode + 1) % 3
	if oldMode == cameraNormal || c.cameraMode == cameraNormal {
		// Reset the entity
		oldEntity := c.entity
		c.entities.container.RemoveEntity(oldEntity)
		c.initEntity(c.cameraMode != cameraNormal)
		c.entities.container.AddEntity(c.entity)
		c.entity.SetPosition(oldEntity.Position())
		c.entity.SetTargetPosition(oldEntity.TargetPosition())
	}
}

func (c *ClientState) renderTick(delta float64) {
	c.delta = delta
	c.hotbarUI.SetX(-184 + 24 + 40*float64(c.currentHotbarSlot))
	c.tickItemName()

	forward, yaw := c.calculateMovement()

	c.LX, c.LY, c.LZ = c.X, c.Y, c.Z
	lx, ly, lz := c.X, c.Y, c.Z

	if c.GameMode.Fly() {
		c.X += forward * math.Cos(yaw) * -math.Cos(c.Pitch) * delta * 0.2
		c.Z -= forward * math.Sin(yaw) * -math.Cos(c.Pitch) * delta * 0.2
		c.Y -= forward * math.Sin(c.Pitch) * delta * 0.2
	} else if chunkMap[chunkPosition{int(math.Floor(c.X)) >> 4, int(math.Floor(c.Z)) >> 4}] != nil {
		speed := 4.317 / 60.0
		if c.KeyState[KeySprint] {
			speed = 5.612 / 60.0
		}
		if _, ok := chunkMap.Block(int(math.Floor(c.X)), int(math.Floor(c.Y)), int(math.Floor(c.Z))).(*blockLiquid); ok {
			speed = 2.20 / 60.0
			if c.KeyState[KeyJump] {
				c.VSpeed = 0.05
			} else {
				c.VSpeed -= 0.005 * delta
				if c.VSpeed < -0.05 {
					c.VSpeed = -0.05
				}
			}
		} else if !c.OnGround {
			c.VSpeed -= 0.01 * delta
			if c.VSpeed < -0.3 {
				c.VSpeed = -0.3
			}
		} else if c.KeyState[KeyJump] {
			c.VSpeed = 0.15
		} else {
			c.VSpeed = 0
		}
		c.X += forward * math.Cos(yaw) * delta * speed
		c.Z -= forward * math.Sin(yaw) * delta * speed
		c.Y += c.VSpeed * delta
	}

	if !c.GameMode.NoClip() {
		cx := c.X
		cy := c.Y
		cz := c.Z
		c.Y = c.LY
		c.Z = c.LZ

		// We handle each axis separately to allow for a sliding
		// effect when pushing up against walls.

		bounds, xhit := c.checkCollisions(c.Bounds)
		c.X = float64(bounds.Min[0] + 0.3)
		c.LX = c.X

		c.Z = cz
		bounds, zhit := c.checkCollisions(c.Bounds)
		c.Z = float64(bounds.Min[2] + 0.3)
		c.LZ = c.Z

		// Half block jumps
		// Minecraft lets you 'jump' up 0.5 blocks
		// for slabs and stairs (or smaller blocks).
		// Currently we implement this as a teleport to the
		// top of the block if we could move there
		// but this isn't smooth.
		if (xhit || zhit) && c.OnGround {
			ox, oz := c.X, c.Z
			c.X, c.Z = cx, cz
			for i := 1.0 / 16.0; i <= 0.5; i += 1.0 / 16.0 {
				mini := c.Bounds.Shift(0, float32(i), 0)
				_, hit := c.checkCollisions(mini)
				if !hit {
					cy += i
					ox, oz = c.X, c.Z
					break
				}
			}
			c.X, c.Z = ox, oz
		}

		c.Y = cy
		bounds, yhit := c.checkCollisions(c.Bounds)
		c.Y = float64(bounds.Min.Y())
		if yhit {
			c.VSpeed = 0
		}

		c.checkGround()
	}

	c.Pitch = math.Mod(c.Pitch, math.Pi*2)
	c.Yaw = math.Mod(c.Yaw, math.Pi*2)
	if c.Pitch < 0 {
		c.Pitch += math.Pi * 2
	}
	if c.Yaw < 0 {
		c.Yaw += math.Pi * 2
	}

	ox := math.Cos(c.Yaw-math.Pi/2) * 0.25
	oz := -math.Sin(c.Yaw-math.Pi/2) * 0.25
	c.entity.SetTargetPosition(c.X-ox, c.Y, c.Z-oz)
	c.entity.SetTargetYaw(-c.Yaw)
	c.entity.SetTargetPitch(-c.Pitch - math.Pi)
	c.entity.walking = c.X != lx || c.Y != ly || c.Z != lz

	//  Highlights the target block
	c.highlightTarget()

	// Debug displays
	c.renderDebug()

	c.armTick()
	c.chat.Draw(delta)

	c.playerList.render(delta)
	c.entities.tick()
	c.copyToCamera()
}

func (c *ClientState) tickItemName() {
	item := c.playerInventory.Items[invPlayerHotbarOffset+c.currentHotbarSlot]
	if c.lastHotbarSlot != c.currentHotbarSlot || item != c.lastHotbarItem {
		c.lastHotbarSlot = c.currentHotbarSlot
		c.lastHotbarItem = item
		c.entity.SetCurrentItem(item)
		if item != nil {
			var name chat.AnyComponent
			if di, ok := item.Type.(DisplayTag); ok && di.DisplayName() != "" {
				name = chat.AnyComponent{Value: &chat.TextComponent{Text: di.DisplayName()}}
				chat.ConvertLegacy(name)
			} else {
				name = chat.AnyComponent{Value: &chat.TranslateComponent{Translate: item.Type.NameLocaleKey()}}
			}
			c.itemNameUI.Update(name)
			c.itemNameTimer = 120
		} else {
			c.itemNameUI.Update(chat.AnyComponent{Value: &chat.TextComponent{}})
			c.itemNameTimer = 0
		}
	}

	c.itemNameTimer -= Client.delta
	if c.itemNameTimer < 0 {
		c.itemNameTimer = 0
	}
	a := c.itemNameTimer / 30
	if a > 1 {
		a = 1
	}
	for _, txt := range c.itemNameUI.Text {
		txt.SetA(int(a * 255))
	}
}

func (c *ClientState) armTick() {
	if c.isLeftDown {
		c.swingTimer -= c.delta
		if c.swingTimer < 0 {
			c.swingTimer = 15
			c.entity.SwingArm()
			c.network.Write(&protocol.ArmSwing{})
			e := c.targetEntity()
			if ne, ok := e.(NetworkComponent); ok {
				c.network.Write(&protocol.UseEntity{
					TargetID: protocol.VarInt(ne.EntityID()),
					Type:     1, // Attack
				})
				return
			}
		}
	}
	pos, b, face, _ := c.targetBlock()
	if c.isLeftDown {
		if b != c.currentBreakingBlock || pos != c.currentBreakingPos {
			c.currentBreakingBlock = Blocks.Air.Base
			c.network.Write(&protocol.PlayerDigging{
				Status:   1, // Cancel
				Location: protocol.NewPosition(pos.X, pos.Y, pos.Z),
				Face:     directionToProtocol(face),
			})
			c.killBreakEntity()
		}
		if c.currentBreakingBlock.Is(Blocks.Air) {
			if math.IsInf(b.Hardness(), 1) {
				return
			}
			c.network.Write(&protocol.PlayerDigging{
				Status:   0, // Start
				Location: protocol.NewPosition(pos.X, pos.Y, pos.Z),
				Face:     directionToProtocol(face),
			})
			c.breakTime = b.Hardness() * 1.5 * 60.0
			c.maxBreakTime = c.breakTime
			c.currentBreakingBlock = b
			c.currentBreakingPos = pos
			c.breakEntity = newBlockBreakEntity()
			c.breakEntity.SetPosition(pos)
			c.breakEntity.(BlockBreakComponent).Update()
		} else {
			c.breakTime -= c.delta
			if c.breakTime < 0 {
				c.breakTime = 0
				c.currentBreakingBlock = Blocks.Air.Base
				c.network.Write(&protocol.PlayerDigging{
					Status:   2, // Finish
					Location: protocol.NewPosition(pos.X, pos.Y, pos.Z),
					Face:     directionToProtocol(face),
				})
				chunkMap.SetBlock(Blocks.Air.Base, pos.X, pos.Y, pos.Z)
				chunkMap.UpdateBlock(pos.X, pos.Y, pos.Z)
				c.killBreakEntity()
			} else {
				stage := int(9 - math.Min(9, 10*(c.breakTime/c.maxBreakTime)))
				if stage != c.breakEntity.(BlockBreakComponent).Stage() {
					c.breakEntity.(BlockBreakComponent).SetStage(stage)
					c.breakEntity.(BlockBreakComponent).Update()
				}
			}
		}
	} else if !c.currentBreakingBlock.Is(Blocks.Air) {
		c.currentBreakingBlock = Blocks.Air.Base
		c.network.Write(&protocol.PlayerDigging{
			Status:   1, // Cancel
			Location: protocol.NewPosition(pos.X, pos.Y, pos.Z),
			Face:     directionToProtocol(face),
		})
		c.killBreakEntity()
	}
}

func (c *ClientState) killBreakEntity() {
	if c.breakEntity != nil {
		c.entities.container.RemoveEntity(c.breakEntity)
		c.breakEntity = nil
	}
}

func (c *ClientState) UpdateHealth(health float64) {
	const maxHealth = 20.0
	c.Health = health
	hp := (health / maxHealth) * float64(len(c.lifeFillUI))
	for i, img := range c.lifeFillUI {
		i := float64(i)
		if i+0.5 < hp {
			img.SetDraw(true)
			img.SetTextureWidth(9.0 / 256.0)
			img.SetWidth(18)
		} else if i < hp {
			img.SetDraw(true)
			img.SetTextureWidth(4.5 / 256.0)
			img.SetWidth(9)
		} else {
			img.SetDraw(false)
		}
	}
	if health == 0.0 {
		setScreen(newRespawnScreen())
	}
}

func (c *ClientState) UpdateHunger(hunger float64) {
	const maxHunger = 20.0
	c.Hunger = hunger
	hp := (hunger / maxHunger) * float64(len(c.foodFillUI))
	for i, img := range c.foodFillUI {
		i := float64(i)
		if i+0.5 < hp {
			img.SetDraw(true)
			img.SetTextureX((16 + 9*4) / 256.0)
			img.SetTextureWidth(9.0 / 256.0)
			img.SetWidth(18)
		} else if i < hp {
			img.SetDraw(true)
			img.SetTextureX((16+9*4)/256.0 + (4.5 / 256.0))
			img.SetTextureWidth(4.5 / 256.0)
			img.SetWidth(9)
		} else {
			img.SetDraw(false)
		}
	}
}

func (c *ClientState) MouseAction(button glfw.MouseButton, down bool) {
	if button == glfw.MouseButtonLeft {
		c.isLeftDown = down
	} else if button == glfw.MouseButtonRight && down {
		e := c.targetEntity()
		if ne, ok := e.(NetworkComponent); ok {
			c.network.Write(&protocol.UseEntity{
				TargetID: protocol.VarInt(ne.EntityID()),
				Type:     0, // Interact
			})
			return
		}
		if c.playerInventory.Items[c.currentHotbarSlot+invPlayerHotbarOffset] != nil {
			c.network.Write(&protocol.PlayerBlockPlacement{
				Face: 0xFF,
			})
		}

		pos, b, face, cur := c.targetBlock()
		if b.Is(Blocks.Air) {
			return
		}
		c.entity.SwingArm()
		c.network.Write(&protocol.ArmSwing{})
		c.network.Write(&protocol.PlayerBlockPlacement{
			Location: protocol.NewPosition(pos.X, pos.Y, pos.Z),
			Face:     directionToProtocol(face),
			CursorX:  byte(cur.X() * 16),
			CursorY:  byte(cur.Y() * 16),
			CursorZ:  byte(cur.Z() * 16),
		})
	}
}

func directionToProtocol(d direction.Type) byte {
	switch d {
	case direction.Up:
		return 1
	case direction.Down:
		return 0
	default:
		return byte(d)
	}
}

func (c *ClientState) targetEntity() (e Entity) {
	s := mgl32.Vec3{float32(render.Camera.X), float32(render.Camera.Y), float32(render.Camera.Z)}
	d := c.viewVector()

	bounds := vmath.NewAABB(0, 0, 0, 1, 1, 1)
	traceRay(
		4,
		s, d,
		func(bx, by, bz int) bool {
			ents := chunkMap.EntitiesIn(bounds.Shift(float32(bx), float32(by), float32(bz)))
			for _, ee := range ents {
				ex, ey, ez := ee.(PositionComponent).Position()
				bo := ee.(SizeComponent).Bounds().Shift(float32(ex), float32(ey), float32(ez))
				if _, ok := bo.IntersectsLine(s, d); ok {
					e = ee
					return false
				}
			}

			b := chunkMap.Block(bx, by, bz)
			if _, ok := b.(*blockLiquid); !b.Is(Blocks.Air) && !ok {
				bb := b.CollisionBounds()
				for _, bound := range bb {
					bound = bound.Shift(float32(bx), float32(by), float32(bz))
					if _, ok := bound.IntersectsLine(s, d); ok {
						return false
					}
				}
			}
			return true
		},
	)
	return
}

func (c *ClientState) targetBlock() (pos Position, block Block, face direction.Type, cursor mgl32.Vec3) {
	s := mgl32.Vec3{float32(render.Camera.X), float32(render.Camera.Y), float32(render.Camera.Z)}
	d := c.viewVector()
	face = direction.Invalid

	block = Blocks.Air.Base
	bounds := vmath.NewAABB(0, 0, 0, 1, 1, 1)
	traceRay(
		4,
		s, d,
		func(bx, by, bz int) bool {
			ents := chunkMap.EntitiesIn(bounds.Shift(float32(bx), float32(by), float32(bz)))
			for _, ee := range ents {
				ex, ey, ez := ee.(PositionComponent).Position()
				bo := ee.(SizeComponent).Bounds().Shift(float32(ex), float32(ey), float32(ez))
				if _, ok := bo.IntersectsLine(s, d); ok {
					return false
				}
			}

			b := chunkMap.Block(bx, by, bz)
			if _, ok := b.(*blockLiquid); !b.Is(Blocks.Air) && !ok {
				bb := b.CollisionBounds()
				for _, bound := range bb {
					bound = bound.Shift(float32(bx), float32(by), float32(bz))
					if at, ok := bound.IntersectsLine(s, d); ok {
						pos = Position{bx, by, bz}
						block = b
						face = findFace(bound, at)
						cursor = at.Sub(mgl32.Vec3{float32(bx), float32(by), float32(bz)})
						return false
					}
				}
			}
			return true
		},
	)
	return
}

func findFace(bound vmath.AABB, at mgl32.Vec3) direction.Type {
	switch {
	case bound.Min.X() == at.X():
		return direction.West
	case bound.Max.X() == at.X():
		return direction.East
	case bound.Min.Y() == at.Y():
		return direction.Down
	case bound.Max.Y() == at.Y():
		return direction.Up
	case bound.Min.Z() == at.Z():
		return direction.North
	case bound.Max.Z() == at.Z():
		return direction.South
	}
	return direction.Up
}

func traceRay(max float32, s, d mgl32.Vec3, cb func(x, y, z int) bool) {
	type gen struct {
		count   int
		base, d float32
	}
	newGen := func(start, d float32) *gen {
		g := &gen{}
		if d > 0 {
			g.base = (float32(math.Ceil(float64(start))) - start) / d
		} else if d < 0 {
			d = float32(math.Abs(float64(d)))
			g.base = (start - float32(math.Floor(float64(start)))) / d
		}
		g.d = d
		return g
	}
	next := func(g *gen) float32 {
		g.count++
		if g.d == 0 {
			return float32(math.Inf(1))
		}
		return g.base + float32(g.count-1)/g.d
	}

	aGen := newGen(s.X(), d.X())
	bGen := newGen(s.Y(), d.Y())
	cGen := newGen(s.Z(), d.Z())
	nextNA := next(aGen)
	nextNB := next(bGen)
	nextNC := next(cGen)

	x, y, z := int(math.Floor(float64(s.X()))), int(math.Floor(float64(s.Y()))), int(math.Floor(float64(s.Z())))

	for {
		if !cb(x, y, z) {
			return
		}
		nextN := float32(0.0)
		if nextNA <= nextNB {
			if nextNA <= nextNC {
				nextN = nextNA
				nextNA = next(aGen)
				x += int(math.Copysign(1, float64(d.X())))
			} else {
				nextN = nextNC
				nextNC = next(cGen)
				z += int(math.Copysign(1, float64(d.Z())))
			}
		} else {
			if nextNB <= nextNC {
				nextN = nextNB
				nextNB = next(bGen)
				y += int(math.Copysign(1, float64(d.Y())))
			} else {
				nextN = nextNC
				nextNC = next(cGen)
				z += int(math.Copysign(1, float64(d.Z())))
			}
		}
		if nextN > max {
			break
		}
	}
}

// draws a box around the target block using the collision
// box for the shape
func (c *ClientState) highlightTarget() {
	if c.GameMode == gmSpecator {
		return
	}
	const lineSize = 1.0 / 128.0
	t, b, _, _ := c.targetBlock()
	if b.Is(Blocks.Air) {
		return
	}
	for _, b := range b.CollisionBounds() {
		b = b.Shift(float32(t.X), float32(t.Y), float32(t.Z))

		points := [][2]float64{
			{float64(b.Min.X()), float64(b.Min.Z())},
			{float64(b.Min.X()), float64(b.Max.Z())},
			{float64(b.Max.X()), float64(b.Min.Z())},
			{float64(b.Max.X()), float64(b.Max.Z())},
		}

		for _, p := range points {
			render.DrawBox(
				p[0]-lineSize, float64(b.Min.Y())-lineSize, p[1]-lineSize,
				p[0]+lineSize, float64(b.Max.Y())+lineSize, p[1]+lineSize,
				0, 0, 0, 255,
			)
		}

		topPoints := [][4]float64{
			{float64(b.Min.X()), float64(b.Min.Z()), float64(b.Max.X()), float64(b.Min.Z())},
			{float64(b.Min.X()), float64(b.Max.Z()), float64(b.Max.X()), float64(b.Max.Z())},
			{float64(b.Min.X()), float64(b.Min.Z()), float64(b.Min.X()), float64(b.Max.Z())},
			{float64(b.Max.X()), float64(b.Min.Z()), float64(b.Max.X()), float64(b.Max.Z())},
		}
		for _, p := range topPoints {
			p2 := p[2:]
			render.DrawBox(
				p[0]-lineSize, float64(b.Min.Y())-lineSize, p[1]-lineSize,
				p2[0]+lineSize, float64(b.Min.Y())+lineSize, p2[1]+lineSize,
				0, 0, 0, 255,
			)
			render.DrawBox(
				p[0]-lineSize, float64(b.Max.Y())-lineSize, p[1]-lineSize,
				p2[0]+lineSize, float64(b.Max.Y())+lineSize, p2[1]+lineSize,
				0, 0, 0, 255,
			)
		}
	}
}

func (c *ClientState) checkGround() {
	ground := vmath.AABB{
		Min: mgl32.Vec3{-0.3, -0.05, -0.3},
		Max: mgl32.Vec3{0.3, 0.0, 0.3},
	}
	prev := c.OnGround
	_, c.OnGround = c.checkCollisions(ground)
	if !prev && c.OnGround {
		c.didTouchGround = true
	}
}

func (c *ClientState) calculateMovement() (float64, float64) {
	forward := 0.0
	yaw := c.Yaw - math.Pi/2
	if c.KeyState[KeyForward] || c.KeyState[KeyBackwards] {
		forward = 1
		if c.KeyState[KeyBackwards] {
			yaw += math.Pi
		}
	}
	change := 0.0
	if c.KeyState[KeyLeft] {
		change = (math.Pi / 2) / (math.Abs(forward) + 1)
	}
	if c.KeyState[KeyRight] {
		change = -(math.Pi / 2) / (math.Abs(forward) + 1)
	}
	if c.KeyState[KeyRight] || c.KeyState[KeyLeft] {
		forward = 1
	}
	if c.KeyState[KeyBackwards] {
		yaw -= change
	} else {
		yaw += change
	}
	return forward, yaw
}

func (c *ClientState) facingDirection() direction.Type {
	viewVector := c.viewVector()
	for _, d := range direction.Values {
		if d.AsVec().Dot(viewVector) > 0.5 {
			return d
		}
	}
	return direction.Invalid
}

func (c *ClientState) viewVector() mgl32.Vec3 {
	return mgl32.Vec3{
		float32(math.Cos(c.Yaw-math.Pi/2) * -math.Cos(c.Pitch)),
		float32(-math.Sin(c.Pitch)),
		float32(-math.Sin(c.Yaw-math.Pi/2) * -math.Cos(c.Pitch)),
	}
}

func (c *ClientState) checkCollisions(bounds vmath.AABB) (vmath.AABB, bool) {
	bounds = bounds.Shift(float32(c.X), float32(c.Y), float32(c.Z))

	dir := mgl32.Vec3{
		-float32(c.LX - c.X),
		-float32(c.LY - c.Y),
		-float32(c.LZ - c.Z),
	}

	minX, minY, minZ := int(bounds.Min.X()-1), int(bounds.Min.Y()-1), int(bounds.Min.Z()-1)
	maxX, maxY, maxZ := int(bounds.Max.X()+1), int(bounds.Max.Y()+1), int(bounds.Max.Z()+1)

	hit := false
	for y := minY; y < maxY; y++ {
		for z := minZ; z < maxZ; z++ {
			for x := minX; x < maxX; x++ {
				b := chunkMap.Block(x, y, z)

				if b.Collidable() {
					for _, bb := range b.CollisionBounds() {
						bb = bb.Shift(float32(x), float32(y), float32(z))
						if bb.Intersects(bounds) {
							bounds = bounds.MoveOutOf(bb, dir)
							hit = true
						}
					}
				}
			}
		}
	}
	return bounds, hit
}

func (c *ClientState) copyToCamera() {
	x, y, z := c.entity.Position()

	ox := math.Cos(-c.entity.Yaw()-math.Pi/2) * 0.25
	oz := -math.Sin(-c.entity.Yaw()-math.Pi/2) * 0.25
	x += ox
	z += oz
	render.Camera.X = x
	render.Camera.Y = y + playerHeight
	render.Camera.Z = z
	render.Camera.Yaw = -c.entity.Yaw()
	render.Camera.Pitch = -c.entity.Pitch() + math.Pi
	switch c.cameraMode {
	case cameraBehind:
		render.Camera.X -= 4 * math.Cos(-c.entity.Yaw()-math.Pi/2) * -math.Cos(-c.entity.Pitch()+math.Pi)
		render.Camera.Z += 4 * math.Sin(-c.entity.Yaw()-math.Pi/2) * -math.Cos(-c.entity.Pitch()+math.Pi)
		render.Camera.Y += 4 * math.Sin(-c.entity.Pitch()+math.Pi)
	case cameraFront:
		render.Camera.X += 4 * math.Cos(-c.entity.Yaw()-math.Pi/2) * -math.Cos(-c.entity.Pitch()+math.Pi)
		render.Camera.Z -= 4 * math.Sin(-c.entity.Yaw()-math.Pi/2) * -math.Cos(-c.entity.Pitch()+math.Pi)
		render.Camera.Y -= 4 * math.Sin(-c.entity.Pitch()+math.Pi)
		render.Camera.Yaw += math.Pi
		render.Camera.Pitch = -render.Camera.Pitch
	}
}

func (c *ClientState) tick() {
	// Now you may be wondering why we have to spam movement
	// packets (any of the Player* move/look packets) 20 times
	// a second instead of only sending when something changes.
	// This is because the server only ticks certain parts of
	// the player when a movement packet is recieved meaning
	// if we sent them any slower health regen would be slowed
	// down as well and various other things too (potions, speed
	// hack check). This also has issues if we send them too
	// fast as well since we will regen health at much faster
	// rates than normal players and some modded servers will
	// (correctly) detect this as cheating. Its Minecraft
	// what did you expect?
	// TODO(Think) Use the smaller packets when possible

	// Force the server to know when touched the ground
	// otherwise if it happens between ticks the server
	// will think we are flying.
	onGround := c.OnGround
	if c.didTouchGround {
		c.didTouchGround = false
		onGround = true
	}

	if c.Health > 0 {
		c.network.Write(&protocol.PlayerPositionLook{
			X:        c.X,
			Y:        c.Y,
			Z:        c.Z,
			Yaw:      float32(-c.Yaw * (180 / math.Pi)),
			Pitch:    float32((-c.Pitch - math.Pi) * (180 / math.Pi)),
			OnGround: onGround,
		})
	}
}

type gameMode int

const (
	gmSurvival gameMode = iota
	gmCreative
	gmAdventure
	gmSpecator
)

func (g gameMode) Fly() bool {
	switch g {
	case gmCreative, gmSpecator:
		return true
	}
	return false
}

func (g gameMode) NoClip() bool {
	switch g {
	case gmSpecator:
		return true
	}
	return false
}

type teleportFlag byte

const (
	teleportRelX teleportFlag = 1 << iota
	teleportRelY
	teleportRelZ
	teleportRelYaw
	teleportRelPitch
)

func calculateTeleport(flag teleportFlag, flags byte, base, val float64) float64 {
	if flags&byte(flag) != 0 {
		return base + val
	}
	return val
}
