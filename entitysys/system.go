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

// Package entitysys provides a simple entity component system
// for handling entities.
package entitysys

import "reflect"

// Stage is used to specify when a system is called.
type Stage int

const (
	// Add marks the system for being called when a entity is
	// added to a container.
	Add Stage = iota
	// Tick marks the system for being called when a entity is
	// ticked.
	Tick
	// Remove marks the system for being called when a entity is
	// removed from a container.
	Remove
)

// Container stores multiple systems and their entities.
type Container struct {
	systems     []*system
	preSystems  []*system
	postSystems []*system
}

// NewContainer creates a new Container.
func NewContainer() *Container {
	return &Container{}
}

// AddEntity adds the entity to all systems that are compatible
// with the entity.
func (c *Container) AddEntity(entity interface{}) {
	re := reflect.ValueOf(entity)
	for _, sys := range c.systems {
		if !sys.Matches(entity) {
			continue
		}
		se := &systemEntity{
			v:      re,
			params: make([]reflect.Value, len(sys.params)),
		}
		se.params = genParams(re, sys.params)

		sys.entities = append(sys.entities, se)
	}
	for _, sys := range c.preSystems {
		if !sys.Matches(entity) {
			continue
		}
		params := genParams(re, sys.params)
		sys.f.Call(params)
	}
}

// RemoveEntity removes the entity from all systems it is
// attached too.
func (c *Container) RemoveEntity(e interface{}) {
	re := reflect.ValueOf(e)
	for _, sys := range c.systems {
		if !sys.Matches(e) {
			continue
		}
	seLoop:
		for i, se := range sys.entities {
			if se.v == re {
				sys.entities = append(sys.entities[:i], sys.entities[i+1:]...)
				break seLoop
			}
		}
	}
	for _, sys := range c.postSystems {
		if !sys.Matches(e) {
			continue
		}
		params := genParams(re, sys.params)
		sys.f.Call(params)
	}
}

func genParams(re reflect.Value, sparams []reflect.Type) []reflect.Value {
	params := make([]reflect.Value, len(sparams))
	for i, p := range sparams {
		if p.Kind() == reflect.Interface {
			params[i] = re
		} else {
			params[i] = re.Elem().FieldByName(p.Elem().Name()).Addr()
		}
	}
	return params
}

// AddSystem adds the system to the container, the passed desc
// values will be used to match when an entity is added. f will
// called for all matching entities depending on the stage. All
// parameters to f are automatically added to matchers.
func (c *Container) AddSystem(stage Stage, f interface{}, matchers ...Matcher) {
	s := &system{
		f:        reflect.ValueOf(f),
		matchers: matchers,
	}
	t := s.f.Type()
	for i := 0; i < t.NumIn(); i++ {
		s.params = append(s.params, t.In(i))
		s.matchers = append(s.matchers, typeMatcher{Type: t.In(i)})
	}
	switch stage {
	case Add:
		c.preSystems = append(c.preSystems, s)
	case Tick:
		c.systems = append(c.systems, s)
	case Remove:
		c.postSystems = append(c.postSystems, s)
	}
}

// Tick ticks all systems and their entities.
func (c *Container) Tick() {
	for _, sys := range c.systems {
		for _, e := range sys.entities {
			sys.f.Call(e.params)
		}
	}
}

type system struct {
	f        reflect.Value
	params   []reflect.Type
	matchers []Matcher

	entities []*systemEntity
}

func (s *system) Matches(e interface{}) bool {
	for _, matcher := range s.matchers {
		if !matcher.Match(e) {
			return false
		}
	}
	return true
}

type systemEntity struct {
	v reflect.Value

	params []reflect.Value
}
