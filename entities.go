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

import "github.com/thinkofdeath/steven/type/vmath"

func newPlayer() Entity {
	type player struct {
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		playerComponent
		playerModelComponent
	}
	p := &player{}
	p.bounds = vmath.NewAABB(-0.3, 0, -0.3, 0.6, 1.8, 0.6)
	return p
}

func newCreeper() Entity {
	type creeper struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	c := &creeper{
		debugComponent: debugComponent{16, 117, 55},
	}
	c.networkID = 50
	c.bounds = vmath.NewAABB(-0.2, 0, -0.2, 0.4, 1.5, 0.4)
	return c
}

func newSkeleton() Entity {
	type skeleton struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	s := &skeleton{
		debugComponent: debugComponent{255, 255, 255},
	}
	s.networkID = 51
	s.bounds = vmath.NewAABB(-0.3, 0, -0.3, 0.6, 1.8, 0.6)
	return s
}

func newSpider() Entity {
	type spider struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	s := &spider{
		debugComponent: debugComponent{59, 7, 7},
	}
	s.networkID = 52
	s.bounds = vmath.NewAABB(-0.7, 0, -0.7, 1.4, 0.9, 1.4)
	return s
}

func newZombie() Entity {
	type zombie struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	z := &zombie{
		debugComponent: debugComponent{17, 114, 156},
	}
	z.networkID = 54
	z.bounds = vmath.NewAABB(-0.3, 0, -0.3, 0.6, 1.8, 0.6)
	return z
}

func newSlime() Entity {
	type slime struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	s := &slime{
		debugComponent: debugComponent{17, 114, 156},
	}
	s.networkID = 55
	s.bounds = vmath.NewAABB(-0.5, 0, -0.5, 1, 1, 1)
	return s
}

func newGhast() Entity {
	type ghast struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	g := &ghast{
		debugComponent: debugComponent{191, 191, 191},
	}
	g.networkID = 56
	g.bounds = vmath.NewAABB(-2, 0, -2, 4, 4, 4)
	return g
}

func newZombiePigman() Entity {
	type zombiePigman struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	z := &zombiePigman{
		debugComponent: debugComponent{204, 110, 198},
	}
	z.networkID = 57
	z.bounds = vmath.NewAABB(-0.3, 0, -0.3, 0.6, 1.8, 0.6)
	return z
}

func newEnderman() Entity {
	type enderman struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	e := &enderman{
		debugComponent: debugComponent{74, 0, 69},
	}
	e.networkID = 58
	e.bounds = vmath.NewAABB(-0.3, 0, -0.3, 0.6, 2.9, 0.6)
	return e
}

func newCaveSpider() Entity {
	type caveSpider struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	c := &caveSpider{
		debugComponent: debugComponent{0, 116, 232},
	}
	c.networkID = 59
	c.bounds = vmath.NewAABB(-0.35, 0, -0.35, 0.7, 0.5, 0.7)
	return c
}

func newSilverfish() Entity {
	type silverfish struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	s := &silverfish{
		debugComponent: debugComponent{128, 128, 128},
	}
	s.networkID = 60
	s.bounds = vmath.NewAABB(-0.2, 0, -0.2, 0.4, 0.3, 0.4)
	return s
}

func newBlaze() Entity {
	type blaze struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	b := &blaze{
		debugComponent: debugComponent{184, 61, 0},
	}
	b.networkID = 61
	b.bounds = vmath.NewAABB(-0.3, 0, -0.3, 0.6, 1.8, 0.6)
	return b
}

func newMagmaCube() Entity {
	type magmaCube struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	m := &magmaCube{
		debugComponent: debugComponent{186, 28, 28},
	}
	m.networkID = 62
	m.bounds = vmath.NewAABB(-0.5, 0, -0.5, 1, 1, 1)
	return m
}

func newEnderDragon() Entity {
	type enderDragon struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	e := &enderDragon{
		debugComponent: debugComponent{122, 59, 117},
	}
	e.networkID = 63
	e.bounds = vmath.NewAABB(-8, 0, -8, 16, 8, 16)
	return e
}

func newWither() Entity {
	type wither struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	w := &wither{
		debugComponent: debugComponent{64, 64, 64},
	}
	w.networkID = 64
	w.bounds = vmath.NewAABB(-0.45, 0, -0.45, 0.9, 3.5, 0.9)
	return w
}

func newBat() Entity {
	type bat struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	b := &bat{
		debugComponent: debugComponent{8, 8, 8},
	}
	b.networkID = 65
	b.bounds = vmath.NewAABB(-0.25, 0, -0.25, 0.5, 0.9, 0.5)
	return b
}

func newWitch() Entity {
	type witch struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	w := &witch{
		debugComponent: debugComponent{87, 64, 0},
	}
	w.networkID = 66
	w.bounds = vmath.NewAABB(-0.3, 0, -0.3, 0.6, 1.8, 0.6)
	return w
}

func newEndermite() Entity {
	type endermite struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	e := &endermite{
		debugComponent: debugComponent{69, 47, 71},
	}
	e.networkID = 67
	e.bounds = vmath.NewAABB(-0.2, 0, -0.2, 0.4, 0.3, 0.4)
	return e
}

func newGuardian() Entity {
	type guardian struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	g := &guardian{
		debugComponent: debugComponent{69, 47, 71},
	}
	g.networkID = 68
	g.bounds = vmath.NewAABB(-0.425, 0, -0.425, 0.85, 0.85, 0.85)
	return g
}

func newPig() Entity {
	type pig struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	p := &pig{
		debugComponent: debugComponent{252, 0, 194},
	}
	p.networkID = 90
	p.bounds = vmath.NewAABB(-0.45, 0, -0.45, 0.9, 0.9, 0.9)
	return p
}

func newSheep() Entity {
	type sheep struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	s := &sheep{
		debugComponent: debugComponent{232, 232, 232},
	}
	s.networkID = 91
	s.bounds = vmath.NewAABB(-0.45, 0, -0.45, 0.9, 1.3, 0.9)
	return s
}

func newCow() Entity {
	type cow struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	c := &cow{
		debugComponent: debugComponent{125, 52, 0},
	}
	c.networkID = 92
	c.bounds = vmath.NewAABB(-0.45, 0, -0.45, 0.9, 1.3, 0.9)
	return c
}

func newChicken() Entity {
	type chicken struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	c := &chicken{
		debugComponent: debugComponent{217, 217, 217},
	}
	c.networkID = 93
	c.bounds = vmath.NewAABB(-0.2, 0, -0.2, 0.4, 0.7, 0.4)
	return c
}

func newSquid() Entity {
	type squid struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	s := &squid{
		debugComponent: debugComponent{84, 39, 245},
	}
	s.networkID = 94
	s.bounds = vmath.NewAABB(-0.475, 0, -0.475, 0.95, 0.95, 0.95)
	return s
}

func newWolf() Entity {
	type wolf struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	w := &wolf{
		debugComponent: debugComponent{148, 148, 148},
	}
	w.networkID = 95
	w.bounds = vmath.NewAABB(-0.3, 0, -0.3, 0.6, 0.8, 0.6)
	return w
}

func newMooshroom() Entity {
	type mooshroom struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	m := &mooshroom{
		debugComponent: debugComponent{145, 41, 0},
	}
	m.networkID = 96
	m.bounds = vmath.NewAABB(-0.45, 0, -0.45, 0.9, 1.3, 0.9)
	return m
}

func newSnowman() Entity {
	type snowman struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	s := &snowman{
		debugComponent: debugComponent{225, 225, 255},
	}
	s.networkID = 97
	s.bounds = vmath.NewAABB(-0.35, 0, -0.35, 0.7, 1.9, 0.7)
	return s
}

func newOcelot() Entity {
	type ocelot struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	o := &ocelot{
		debugComponent: debugComponent{242, 222, 0},
	}
	o.networkID = 98
	o.bounds = vmath.NewAABB(-0.3, 0, -0.3, 0.6, 0.8, 0.6)
	return o
}

func newIronGolem() Entity {
	type ironGolem struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	i := &ironGolem{
		debugComponent: debugComponent{125, 125, 125},
	}
	i.networkID = 99
	i.bounds = vmath.NewAABB(-0.7, 0, -0.7, 1.4, 2.9, 1.4)
	return i
}

func newHorse() Entity {
	type horse struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	h := &horse{
		debugComponent: debugComponent{191, 156, 0},
	}
	h.networkID = 100
	h.bounds = vmath.NewAABB(-0.7, 0, -0.7, 1.4, 1.6, 1.4)
	return h
}

func newRabbit() Entity {
	type rabbit struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	r := &rabbit{
		debugComponent: debugComponent{181, 123, 42},
	}
	r.networkID = 101
	r.bounds = vmath.NewAABB(-0.3, 0, -0.3, 0.6, 0.7, 0.6)
	return r
}

func newVillager() Entity {
	type villager struct {
		networkComponent
		positionComponent
		rotationComponent
		targetRotationComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	v := &villager{
		debugComponent: debugComponent{212, 183, 142},
	}
	v.networkID = 120
	v.bounds = vmath.NewAABB(-0.3, 0, -0.3, 0.6, 1.8, 0.6)
	return v
}
