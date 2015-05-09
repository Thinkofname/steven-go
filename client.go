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

package phteven

import (
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/thinkofdeath/phteven/protocol"
	"github.com/thinkofdeath/phteven/render"
	"github.com/thinkofdeath/phteven/type/direction"
	"github.com/thinkofdeath/phteven/type/vmath"
	"github.com/thinkofdeath/phteven/ui"
	"github.com/thinkofdeath/phteven/ui/scene"
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

	X, Y, Z    float64
	Yaw, Pitch float64

	Jumping  bool
	VSpeed   float64
	KeyState [5]bool
	OnGround bool

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
	c.scene.AddDrawable(
		ui.NewImage(widgets, 0, 0, 182*2, 22*2, 0, 0, 182.0/256.0, 22.0/256.0, 255, 255, 255).
			Attach(ui.Bottom, ui.Center),
	)
	c.hotbarUI = ui.NewImage(widgets, -22*2+4, -2, 24*2, 24*2, 0, 22.0/256.0, 24.0/256.0, 24.0/256.0, 255, 255, 255).
		Attach(ui.Bottom, ui.Center)
	c.scene.AddDrawable(c.hotbarUI)

	// Hearts / Food
	for i := 0; i < 10; i++ {
		l := ui.NewImage(icons, -182+9+16*float64(i), 22*2+16, 18, 18, 16.0/256.0, 0, 9.0/256.0, 9.0/256.0, 255, 255, 255).
			Attach(ui.Bottom, ui.Center)
		c.scene.AddDrawable(l)
		c.lifeUI = append(c.lifeUI, l)
		l = ui.NewImage(icons, -182+9+16*float64(i), 22*2+16, 18, 18, (16+9*4)/256.0, 0, 9.0/256.0, 9.0/256.0, 255, 255, 255).
			Attach(ui.Bottom, ui.Center)
		c.scene.AddDrawable(l)
		c.lifeFillUI = append(c.lifeFillUI, l)

		f := ui.NewImage(icons, 182+7-16*float64(10-i), 22*2+16, 18, 18, 16.0/256.0, 27.0/256.0, 9.0/256.0, 9.0/256.0, 255, 255, 255).
			Attach(ui.Bottom, ui.Center)
		c.scene.AddDrawable(f)
		c.foodUI = append(c.foodUI, l)
		f = ui.NewImage(icons, 182+7-16*float64(10-i), 22*2+16, 18, 18, (16+9*4)/256.0, 27.0/256.0, 9.0/256.0, 9.0/256.0, 255, 255, 255).
			Attach(ui.Bottom, ui.Center)
		c.scene.AddDrawable(f)
		c.foodFillUI = append(c.foodFillUI, l)
	}

	// Exp bar
	c.scene.AddDrawable(
		ui.NewImage(icons, 0, 22*2+4, 182*2, 10, 0, 64.0/256.0, 182.0/256.0, 5.0/256.0, 255, 255, 255).
			Attach(ui.Bottom, ui.Center),
	)

	c.chat.init()
	c.initDebug()
	c.playerList.init()
}

func (c *ClientState) renderTick(delta float64) {
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

	c.playerList.render(delta)
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
	_, c.OnGround = c.checkCollisions(ground)
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
	writeChan <- &protocol.PlayerPositionLook{
		X:        c.X,
		Y:        c.Y,
		Z:        c.Z,
		Yaw:      float32(-c.Yaw * (180 / math.Pi)),
		Pitch:    float32((-c.Pitch - math.Pi) * (180 / math.Pi)),
		OnGround: c.OnGround,
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
