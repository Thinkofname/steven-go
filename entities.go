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
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	p := &player{
		debugComponent: debugComponent{255, 0, 255},
	}
	p.bounds = vmath.NewAABB(-0.3, 0, -0.3, 0.6, 1.8, 0.6)
	return p
}

func newCreeper() Entity {
	type creeper struct {
		networkComponent
		positionComponent
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

func newCow() Entity {
	type cow struct {
		networkComponent
		positionComponent
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
