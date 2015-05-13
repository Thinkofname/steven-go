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

	"github.com/go-gl/mathgl/mgl32"
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

var Client ClientState

func init() {
	initClient()
}

func initClient() {
	if Client.valid {
		Client.scene.Hide()
		// Cleanup
		for _, e := range Client.entities.entities {
			Client.entities.container.RemoveEntity(e)
		}
		if Client.entity != nil {
			Client.entities.container.RemoveEntity(Client.entity)
		}
		Client.playerList.free()
	}
	Client = ClientState{
		Bounds: vmath.AABB{
			Min: mgl32.Vec3{-0.3, 0, -0.3},
			Max: mgl32.Vec3{0.3, 1.8, 0.3},
		},
		valid: true,
		scene: scene.New(true),
	}
}

type ClientState struct {
	valid bool

	scene *scene.Type

	entity interface {
		PlayerComponent
		PositionComponent
		RotationComponent
		TargetPositionComponent
		TargetRotationComponent
	}
	entityAdded bool

	X, Y, Z    float64
	Yaw, Pitch float64

	Health float64
	Hunger float64

	Jumping                  bool
	VSpeed                   float64
	KeyState                 [5]bool
	OnGround, didTouchGround bool

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
	hotbarUI          *ui.Image
	currentHotbarSlot int
	lifeUI            []*ui.Image
	lifeFillUI        []*ui.Image
	foodUI            []*ui.Image
	foodFillUI        []*ui.Image

	chat       ChatUI
	playerList playerListUI
	entities   clientEntities

	delta float64
}

func (c *ClientState) init() {
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
	c.hotbarUI = ui.NewImage(widgets, -22*2+4, -2, 24*2, 24*2, 0, 22.0/256.0, 24.0/256.0, 24.0/256.0, 255, 255, 255).
		Attach(ui.Bottom, ui.Center)
	c.scene.AddDrawable(c.hotbarUI)

	// Hearts / Food
	for i := 0; i < 10; i++ {
		l := ui.NewImage(icons, 16*float64(i), -16-8-10, 18, 18, 16.0/256.0, 0, 9.0/256.0, 9.0/256.0, 255, 255, 255).
			Attach(ui.Top, ui.Left)
		l.Parent = hotbar
		c.scene.AddDrawable(l)
		c.lifeUI = append(c.lifeUI, l)
		l = ui.NewImage(icons, 16*float64(i), -16-8-10, 18, 18, (16+9*4)/256.0, 0, 9.0/256.0, 9.0/256.0, 255, 255, 255).
			Attach(ui.Top, ui.Left)
		l.Parent = hotbar
		c.scene.AddDrawable(l)
		c.lifeFillUI = append(c.lifeFillUI, l)

		f := ui.NewImage(icons, 16*float64(i), -16-8-10, 18, 18, 16.0/256.0, 27.0/256.0, 9.0/256.0, 9.0/256.0, 255, 255, 255).
			Attach(ui.Top, ui.Right)
		f.Parent = hotbar
		c.scene.AddDrawable(f)
		c.foodUI = append(c.foodUI, f)
		f = ui.NewImage(icons, 16*float64(i), -16-8-10, 18, 18, (16+9*4)/256.0, 27.0/256.0, 9.0/256.0, 9.0/256.0, 255, 255, 255).
			Attach(ui.Top, ui.Right)
		f.Parent = hotbar
		c.scene.AddDrawable(f)
		c.foodFillUI = append(c.foodFillUI, f)
	}

	// Exp bar
	c.scene.AddDrawable(
		ui.NewImage(icons, 0, 22*2+4, 182*2, 10, 0, 64.0/256.0, 182.0/256.0, 5.0/256.0, 255, 255, 255).
			Attach(ui.Bottom, ui.Center),
	)

	c.chat.init()
	c.initDebug()
	c.playerList.init()
	c.entities.init()

	c.initEntity()
}

func (c *ClientState) initEntity() {
	type clientEntity struct {
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		playerComponent
		playerModelComponent
	}
	ce := &clientEntity{}
	ub, _ := hex.DecodeString(profile.ID)
	copy(ce.uuid[:], ub)
	c.entity = ce
	ce.hasHead = false
	ce.bounds = vmath.NewAABB(-0.3, 0, -0.3, 0.6, 1.8, 0.6)
}

func (c *ClientState) renderTick(delta float64) {
	c.delta = delta
	c.hotbarUI.X = -184 + 24 + 40*float64(c.currentHotbarSlot)

	forward, yaw := c.calculateMovement()

	if c.GameMode.Fly() {
		c.X += forward * math.Cos(yaw) * -math.Cos(c.Pitch) * delta * 0.2
		c.Z -= forward * math.Sin(yaw) * -math.Cos(c.Pitch) * delta * 0.2
		c.Y -= forward * math.Sin(c.Pitch) * delta * 0.2
	} else if chunkMap[chunkPosition{int(c.X) >> 4, int(c.Z) >> 4}] != nil {
		speed := 4.317 / 60.0
		if c.KeyState[KeySprint] {
			speed = 5.612 / 60.0
		}
		c.X += forward * math.Cos(yaw) * delta * speed
		c.Z -= forward * math.Sin(yaw) * delta * speed
		if !c.OnGround {
			c.VSpeed -= 0.01 * delta
			if c.VSpeed < -0.3 {
				c.VSpeed = -0.3
			}
		} else if c.Jumping {
			c.VSpeed = 0.15
		} else {
			c.VSpeed = 0
		}
		c.Y += c.VSpeed * delta
	}

	if !c.GameMode.NoClip() {
		cx := c.X
		cy := c.Y
		cz := c.Z
		c.Y = render.Camera.Y - playerHeight
		c.Z = render.Camera.Z

		// We handle each axis separately to allow for a sliding
		// effect when pushing up against walls.

		bounds, xhit := c.checkCollisions(c.Bounds)
		c.X = float64(bounds.Min[0] + 0.3)
		c.copyToCamera()

		c.Z = cz
		bounds, zhit := c.checkCollisions(c.Bounds)
		c.Z = float64(bounds.Min[2] + 0.3)
		c.copyToCamera()

		// Half block jumps
		// Minecraft lets you 'jump' up 0.5 blocks
		// for slabs and stairs (or smaller blocks).
		// Currently we implement this as a teleport to the
		// top of the block if we could move there
		// but this isn't smooth.
		// TODO(Think) Improve this
		if (xhit || zhit) && c.OnGround {
			ox, oz := c.X, c.Z
			c.X, c.Z = cx, cz
			for i := 1.0 / 16.0; i <= 0.5; i += 1.0 / 16.0 {
				mini := c.Bounds
				mini.Shift(0, float32(i), 0)
				_, hit := c.checkCollisions(mini)
				if !hit {
					cy += i
					ox, oz = c.X, c.Z
					break
				}
			}
			c.X, c.Z = ox, oz
		}
		c.copyToCamera()

		c.Y = cy
		bounds, _ = c.checkCollisions(c.Bounds)
		c.Y = float64(bounds.Min.Y())

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

	c.copyToCamera()

	//  Highlights the target block
	c.highlightTarget()

	// Debug displays
	c.renderDebug()

	// Update our entity
	// TODO Should the entity be the main thing
	// instead of duplicating things in the client?
	ox := math.Cos(c.Yaw-math.Pi/2) * 0.3
	oz := -math.Sin(c.Yaw-math.Pi/2) * 0.3
	c.entity.SetTargetPosition(c.X-ox, c.Y, c.Z-oz)
	c.entity.SetYaw(-c.Yaw)
	c.entity.SetTargetYaw(-c.Yaw)

	c.playerList.render(delta)
	c.entities.tick()
}

func (c *ClientState) UpdateHealth(health float64) {
	const maxHealth = 20.0
	c.Health = health
	hp := (health / maxHealth) * float64(len(c.lifeFillUI))
	for i, img := range c.lifeFillUI {
		i := float64(i)
		if i+0.5 < hp {
			img.Visible = true
			img.TW = 9.0 / 256.0
			img.W = 18
		} else if i < hp {
			img.Visible = true
			img.TW = 4.5 / 256.0
			img.W = 9
		} else {
			img.Visible = false
		}
	}
}

func (c *ClientState) UpdateHunger(hunger float64) {
	const maxHunger = 20.0
	c.Hunger = hunger
	hp := (hunger / maxHunger) * float64(len(c.foodFillUI))
	for i, img := range c.foodFillUI {
		i := float64(i)
		if i+0.5 < hp {
			img.Visible = true
			img.TX = (16 + 9*4) / 256.0
			img.TW = 9.0 / 256.0
			img.W = 18
		} else if i < hp {
			img.Visible = true
			img.TX = (16+9*4)/256.0 + (4.5 / 256.0)
			img.TW = 4.5 / 256.0
			img.W = 9
		} else {
			img.Visible = false
		}
	}
}

func (c *ClientState) targetBlock() (x, y, z int, block Block) {
	const max = 4.0
	block = Blocks.Air.Base
	s := mgl32.Vec3{float32(c.X), float32(c.Y + playerHeight), float32(c.Z)}
	d := c.viewVector()

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
	prevN := float32(0.0)
	nextNA := next(aGen)
	nextNB := next(bGen)
	nextNC := next(cGen)
	for {
		nextN := float32(0.0)
		if nextNA < nextNB {
			if nextNA < nextNC {
				nextN = nextNA
				nextNA = next(aGen)
			} else {
				nextN = nextNC
				nextNC = next(cGen)
			}
		} else {
			if nextNB < nextNC {
				nextN = nextNB
				nextNB = next(bGen)
			} else {
				nextN = nextNC
				nextNC = next(cGen)
			}
		}
		if prevN == nextN {
			continue
		}
		final := false
		n := (prevN + nextN) / 2
		if nextN > max {
			final = true
			n = max
		}
		bv := s.Add(d.Mul(n))
		bx, by, bz := int(math.Floor(float64(bv.X()))), int(math.Floor(float64(bv.Y()))), int(math.Floor(float64(bv.Z())))
		b := chunkMap.Block(bx, by, bz)
		if _, ok := b.(*blockLiquid); !b.Is(Blocks.Air) && !ok {
			bb := b.CollisionBounds()
			for _, bound := range bb {
				bound.Shift(float32(bx), float32(by), float32(bz))
				if bound.IntersectsLine(s, d) {
					x, y, z = bx, by, bz
					block = b
					return
				}
			}
		}
		prevN = nextN
		if final {
			break
		}
	}
	return
}

// draws a box around the target block using the collision
// box for the shape
func (c *ClientState) highlightTarget() {
	if c.GameMode == gmSpecator {
		return
	}
	const lineSize = 1.0 / 128.0
	tx, ty, tz, b := c.targetBlock()
	if b.Is(Blocks.Air) {
		return
	}
	for _, b := range b.CollisionBounds() {
		b.Shift(float32(tx), float32(ty), float32(tz))

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
	bounds.Shift(float32(c.X), float32(c.Y), float32(c.Z))

	dir := mgl32.Vec3{
		-float32(render.Camera.X - c.X),
		-float32(render.Camera.Y - playerHeight - c.Y),
		-float32(render.Camera.Z - c.Z),
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
						bb.Shift(float32(x), float32(y), float32(z))
						if bb.Intersects(&bounds) {
							bounds.MoveOutOf(&bb, dir)
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
	render.Camera.X = c.X
	render.Camera.Y = c.Y + playerHeight
	render.Camera.Z = c.Z
	render.Camera.Yaw = c.Yaw
	render.Camera.Pitch = c.Pitch
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

	writeChan <- &protocol.PlayerPositionLook{
		X:        c.X,
		Y:        c.Y,
		Z:        c.Z,
		Yaw:      float32(-c.Yaw * (180 / math.Pi)),
		Pitch:    float32((-c.Pitch - math.Pi) * (180 / math.Pi)),
		OnGround: onGround,
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
