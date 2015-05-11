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

import "github.com/thinkofdeath/steven/entitysys"

var entityTypes = map[int]func() Entity{
	50: newCreeper,
	51: newSkeleton,
	52: newSpider,
	// 53: Giant Zombie, TODO: Do we need this?
	54: newZombie,
	55: newSlime,
	56: newGhast,
	57: newZombiePigman,
	58: newEnderman,
	59: newCaveSpider,
	60: newSilverfish,
	92: newCow,
}

type clientEntities struct {
	entities  map[int]Entity
	container *entitysys.Container
}

func (ce *clientEntities) init() {
	ce.container = entitysys.NewContainer()
	ce.entities = map[int]Entity{}
	ce.register()
}

func (ce *clientEntities) add(id int, e Entity) {
	ce.entities[id] = e
	ce.container.AddEntity(e)
}

func (ce *clientEntities) remove(id int) {
	e, ok := ce.entities[id]
	if !ok {
		return
	}
	delete(ce.entities, id)
	ce.container.RemoveEntity(e)
}

func (ce *clientEntities) tick() {
	ce.container.Tick()
}

type Entity interface{}
